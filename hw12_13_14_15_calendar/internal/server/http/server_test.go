package internalhttp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/stretchr/testify/assert"
)

type MockApp struct{}

func (m *MockApp) CreateEvent(_ storage.Event) error {
	return nil
}

func TestServer(t *testing.T) {
	cfg := ServerConfig{
		Host: "localhost",
		Port: "8080",
	}
	logg := logger.New("info")
	app := &MockApp{}

	server := NewServer(cfg, logg, app)

	ts := httptest.NewServer(server.httpServer.Handler)
	defer ts.Close()

	t.Run("Root Handler", func(t *testing.T) {
		// Создаем HTTP-клиент с таймаутом
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		// Создаем запрос с контекстом
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL, nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "Hello, World!", string(body))
	})

	t.Run("Start and Stop Server", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := server.Start(ctx)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				t.Error(err)
			}
		}()

		time.Sleep(time.Second)

		cancel()
		time.Sleep(time.Second)
	})
}
