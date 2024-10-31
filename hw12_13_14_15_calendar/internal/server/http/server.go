package internalhttp

import (
	"context"
	"encoding/json"
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

func NewServer(cfg ServerConfig, logg *logger.Logger, store storage.Interface) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", rootHandler())
	mux.HandleFunc("/events", createEventHandler(store, logg))
	mux.HandleFunc("/events/", eventHandler(store, logg))

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

func createEventHandler(store storage.Interface, logg *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var event storage.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			logg.Errorf("Failed to decode event: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if err := store.CreateEvent(event); err != nil {
			logg.Errorf("Failed to create event: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logg.Infof("Event created: %s", event.ID)
		w.WriteHeader(http.StatusCreated)
	}
}

func eventHandler(store storage.Interface, logg *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/events/"):]

		switch r.Method {
		case http.MethodGet:
			event, err := store.GetEvent(id)
			if err != nil {
				logg.Errorf("Failed to get event: %v", err)
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(event)

		case http.MethodPut:
			var event storage.Event
			if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
				logg.Errorf("Failed to decode event: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			if err := store.UpdateEvent(id, event); err != nil {
				logg.Errorf("Failed to update event: %v", err)
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}

			logg.Infof("Event updated: %s", id)
			w.WriteHeader(http.StatusOK)

		case http.MethodDelete:
			if err := store.DeleteEvent(id); err != nil {
				logg.Errorf("Failed to delete event: %v", err)
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}

			logg.Infof("Event deleted: %s", id)
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
