package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/config"   //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"   //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/rabbitmq" //nolint
	"github.com/google/uuid"                                          //nolint
)

type Notification struct {
	EventID   uuid.UUID `json:"eventId"`
	Title     string    `json:"title"`
	StartTime int64     `json:"startTime"`
}

func main() {
	configPath := flag.String("config", "configs/sender_config.yaml", "Path to configuration file")
	flag.Parse()

	cfg := config.LoadConfig(*configPath)
	logg := logger.New(cfg.Logger.Level)

	rabbit, err := rabbitmq.New(cfg.RabbitMQ.DSN, logg)
	if err != nil {
		logg.Error("Failed to connect to RabbitMQ: " + err.Error())
		return
	}
	defer func() {
		rabbit.Close()
		logg.Info("RabbitMQ connection closed")
	}()

	deliveries, err := rabbit.Consume("notifications")
	if err != nil {
		logg.Error("Failed to consume from RabbitMQ: " + err.Error())
		return
	}

	logg.Info("Started consuming from RabbitMQ")

	for msg := range deliveries {
		logg.Info("Received notification: " + string(msg.Body))
		err := processNotification(msg.Body)
		if err != nil {
			logg.Error("Failed to process notification: " + err.Error())
		} else {
			logg.Info("Successfully processed notification")
		}
	}
}

func processNotification(body []byte) error {
	var notification Notification
	err := json.Unmarshal(body, &notification)
	if err != nil {
		return fmt.Errorf("failed to unmarshal notification: %w", err)
	}

	fmt.Println("Processing notification:", notification)

	return nil
}
