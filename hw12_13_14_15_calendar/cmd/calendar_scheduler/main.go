package main

import (
	"encoding/json"
	"flag"
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/config"                 //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"                 //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/rabbitmq"               //nolint
	sqlstorage "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage/sql" //nolint
	"github.com/google/uuid"                                                        //nolint
	_ "github.com/lib/pq"                                                           //nolint
)

type Notification struct {
	EventID   uuid.UUID `json:"eventId"`
	Title     string    `json:"title"`
	StartTime int64     `json:"startTime"`
}

func main() {
	configPath := flag.String("config", "configs/scheduler_config.yaml",
		"Path to configuration file")
	flag.Parse()

	cfg := config.LoadConfig(*configPath)
	logg := logger.New(cfg.Logger.Level)

	rabbit, err := rabbitmq.New(cfg.RabbitMQ.DSN, logg)
	if err != nil {
		logg.Error("Failed to connect to RabbitMQ: " + err.Error())
		return
	}
	defer rabbit.Close()

	err = rabbit.DeclareQueue("notifications")
	if err != nil {
		logg.Error("Failed to declare RabbitMQ queue: " + err.Error())
		return
	}

	store, err := sqlstorage.New(cfg.Database.DSN)
	if err != nil {
		logg.Error("Failed to connect to database: " + err.Error())
		return
	}

	for {
		events, err := store.ListEvents()
		if err != nil {
			logg.Error("Failed to list events: " + err.Error())
			continue
		}

		for _, event := range events {
			if time.Until(event.StartTime) < 24*time.Hour {
				notification := Notification{
					EventID:   event.ID,
					Title:     event.Title,
					StartTime: event.StartTime.Unix(),
				}

				body, err := json.Marshal(notification)
				if err != nil {
					logg.Error("Failed to marshal notification: " + err.Error())
					continue
				}

				err = rabbit.Publish("notifications", body)
				if err != nil {
					logg.Error("Failed to publish notification: " + err.Error())
				} else {
					logg.Info("Successfully published notification for event: " + notification.Title)
				}
			}
		}

		err = store.DeleteOldEvents(time.Now().AddDate(-1, 0, 0))
		if err != nil {
			logg.Error("Failed to delete old events: " + err.Error())
		}

		time.Sleep(cfg.Scheduler.Interval)
	}
}
