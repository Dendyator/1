package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" //nolint

	api "github.com/Dendyator/1/hw12_13_14_15_calendar/api/pb"                            //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/config"                       //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"                       //nolint
	internalgrpc "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/server/grpc"     //nolint
	internalhttp "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/server/http"     //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage"                      //nolint
	memorystorage "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage/memory" //nolint
	sqlstorage "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage/sql"       //nolint
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yaml", "Path to configuration file")
}

func main() {
	versionFlag := flag.Bool("version", false, "Display version information")
	flag.Parse()

	if *versionFlag {
		printVersion()
		return
	}

	cfg := config.LoadConfig(configFile)
	logg := logger.New(cfg.Logger.Level)

	logg.Info("Using DSN: " + cfg.Database.DSN)

	time.Sleep(5 * time.Second)

	var store storage.Interface
	if cfg.Database.Driver == "in-memory" {
		store = memorystorage.New()
		logg.Info("Using in-memory storage")
	} else {
		var err error
		store, err = sqlstorage.New(cfg.Database.DSN)
		if err != nil {
			logg.Error("Failed to initialize SQL storage: " + err.Error())
			return
		}
		logg.Info("Using SQL storage")
	}

	serverCfg := internalhttp.ServerConfig{
		Host: cfg.Server.Host,
		Port: cfg.Server.Port,
	}

	httpServer := internalhttp.NewServer(serverCfg, logg, store)

	grpcServer := grpc.NewServer()
	apiServer := internalgrpc.NewGRPCServer(store, logg)
	api.RegisterEventServiceServer(grpcServer, apiServer)
	reflection.Register(grpcServer)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(timeoutCtx); err != nil {
			logg.Error("Failed to stop HTTP server: " + err.Error())
		}

		grpcServer.GracefulStop()
	}()

	go func() {
		listener, err := net.Listen("tcp", ":50051") //nolint
		if err != nil {
			logg.Error("Failed to listen on port 50051: " + err.Error())
			return
		}
		logg.Info("GRPC server is running on port 50051...")
		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logg.Error("Failed to start GRPC server: " + err.Error())
			return
		}
	}()

	logg.Info("Calendar server is running...")
	if err := httpServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logg.Error("Failed to start HTTP server: " + err.Error())
		return
	}
}
