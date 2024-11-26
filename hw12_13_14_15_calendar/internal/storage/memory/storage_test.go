package memorystorage

import (
	"testing"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/stretchr/testify/assert"                             //nolint:depguard
)

func TestStorage(t *testing.T) {
	s := New()
	event1 := storage.Event{ID: "1", Title: "Event 1"}
	event2 := storage.Event{ID: "2", Title: "Event 2"}
	eventUpdated := storage.Event{ID: "1", Title: "Event 1 Updated"}

	t.Run("Create Event", func(t *testing.T) {
		err := s.CreateEvent(event1)
		assert.NoError(t, err)
		err = s.CreateEvent(event1)
		assert.Error(t, err)
	})

	t.Run("Get Event", func(t *testing.T) {
		event, err := s.GetEvent("1")
		assert.NoError(t, err)
		assert.Equal(t, event1, event)
		_, err = s.GetEvent("invalid")
		assert.Error(t, err)
	})

	t.Run("Update Event", func(t *testing.T) {
		err := s.UpdateEvent("1", eventUpdated)
		assert.NoError(t, err)
		event, err := s.GetEvent("1")
		assert.NoError(t, err)
		assert.Equal(t, eventUpdated, event)
		err = s.UpdateEvent("invalid", event2)
		assert.Error(t, err)
	})

	t.Run("Delete Event", func(t *testing.T) {
		err := s.DeleteEvent("1")
		assert.NoError(t, err)
		_, err = s.GetEvent("1")
		assert.Error(t, err)
		err = s.DeleteEvent("invalid")
		assert.Error(t, err)
	})

	t.Run("List Events", func(t *testing.T) {
		err := s.CreateEvent(event1)
		assert.NoError(t, err)
		err = s.CreateEvent(event2)
		assert.NoError(t, err)
		events, err := s.ListEvents()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(events))
	})
}
