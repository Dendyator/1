package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/config"  //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"  //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint

	internalhttp "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/server/http"     //nolint
	memorystorage "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage/memory" //nolint
	sqlstorage "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage/sql"       //nolint
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg := config.LoadConfig(configFile)
	logg := logger.New(cfg.Logger.Level)

	var store storage.Interface
	if cfg.Database.Driver == "in-memory" {
		store = memorystorage.New()
		logg.Info("Using in-memory storage")
	} else {
		var err error
		store, err = sqlstorage.New(cfg.Database.DSN)
		if err != nil {
			logg.Error("Failed to initialize SQL storage: " + err.Error())
			return // Использую return вместо os.Exit
		}
		logg.Info("Using SQL storage")
	}

	serverCfg := internalhttp.ServerConfig{
		Host: cfg.Server.Host,
		Port: cfg.Server.Port,
	}

	server := internalhttp.NewServer(serverCfg, logg, store)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(timeoutCtx); err != nil {
			logg.Error("Failed to stop HTTP server: " + err.Error())
		}
	}()

	logg.Info("Calendar server is running...")
	if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logg.Error("Failed to start HTTP server: " + err.Error())
		return
	}
}
