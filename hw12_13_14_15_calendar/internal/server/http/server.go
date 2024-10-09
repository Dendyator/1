package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Server struct {
	httpServer *http.Server
}

type ServerConfig struct {
	Host string
	Port string
}

type Application interface {
	CreateEvent(event storage.Event) error
}

func NewServer(cfg ServerConfig, logg *logger.Logger, _ Application) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", rootHandler())

	srv := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           loggingMiddleware(logg)(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		err := s.Stop(context.Background())
		if err != nil {
			return
		}
	}()

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			return
		}
	}
}
