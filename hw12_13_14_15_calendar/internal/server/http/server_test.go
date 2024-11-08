package internalhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"                       //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage"                      //nolint
	memorystorage "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage/memory" //nolint
	"github.com/stretchr/testify/assert"
)

func TestCreateEventHandler(t *testing.T) {
	store := memorystorage.New()
	logg := logger.New("info")
	handler := createEventHandler(store, logg)

	t.Run("successful creation", func(t *testing.T) {
		event := storage.Event{
			ID:    "1",
			Title: "Test Event",
		}
		buf, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(buf))
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("duplicate event", func(t *testing.T) {
		event := storage.Event{
			ID:    "1",
			Title: "Duplicate Event",
		}
		buf, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(buf))
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetEventHandler(t *testing.T) {
	store := memorystorage.New()
	logg := logger.New("info")
	handler := getEventHandler(store, logg) // updated to getEventHandler

	event := storage.Event{
		ID:    "1",
		Title: "Test Event",
	}
	_ = store.CreateEvent(event)

	t.Run("successful get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var respEvent storage.Event
		_ = json.NewDecoder(rec.Body).Decode(&respEvent)
		assert.Equal(t, event, respEvent)
	})

	t.Run("event not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/2", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestUpdateEventHandler(t *testing.T) {
	store := memorystorage.New()
	logg := logger.New("info")
	handler := updateEventHandler(store, logg) // updated to updateEventHandler

	event := storage.Event{
		ID:    "1",
		Title: "Test Event",
	}
	_ = store.CreateEvent(event)

	t.Run("successful update", func(t *testing.T) {
		updatedEvent := storage.Event{
			ID:    "1",
			Title: "Updated Event",
		}
		buf, _ := json.Marshal(updatedEvent)
		req := httptest.NewRequest(http.MethodPut, "/events/1", bytes.NewBuffer(buf))
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("update not found", func(t *testing.T) {
		updatedEvent := storage.Event{
			ID:    "2",
			Title: "New Event",
		}
		buf, _ := json.Marshal(updatedEvent)
		req := httptest.NewRequest(http.MethodPut, "/events/2", bytes.NewBuffer(buf))
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestDeleteEventHandler(t *testing.T) {
	store := memorystorage.New()
	logg := logger.New("info")
	handler := deleteEventHandler(store, logg) // updated to deleteEventHandler

	event := storage.Event{
		ID:    "1",
		Title: "Test Event",
	}
	_ = store.CreateEvent(event)

	t.Run("successful delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("delete not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/2", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
