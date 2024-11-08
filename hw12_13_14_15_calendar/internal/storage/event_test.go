package storage_test

import (
	"testing"
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/stretchr/testify/assert"
)

func TestEventStructure(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)

	event := storage.Event{
		ID:          "12345",
		Title:       "Test Event",
		Description: "This is a test description",
		StartTime:   startTime,
		EndTime:     endTime,
	}

	assert.Equal(t, "12345", event.ID)
	assert.Equal(t, "Test Event", event.Title)
	assert.Equal(t, "This is a test description", event.Description)
	assert.Equal(t, startTime, event.StartTime)
	assert.Equal(t, endTime, event.EndTime)
}
