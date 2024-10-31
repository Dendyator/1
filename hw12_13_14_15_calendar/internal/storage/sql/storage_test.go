package sqlstorage_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"                                                //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage"                //nolint
	sqlstorage "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage/sql" //nolint
	"github.com/jmoiron/sqlx"                                                       //nolint
	"github.com/stretchr/testify/assert"                                            //nolint
)

func TestCreateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	storages := &sqlstorage.Storage{DB: sqlxDB}
	event := storage.Event{
		ID:          "1",
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
	}

	mock.ExpectExec("INSERT INTO events").
		WithArgs(event.ID, event.Title, event.Description, event.StartTime, event.EndTime).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = storages.CreateEvent(event)
	assert.NoError(t, err)
}

func TestUpdateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	storages := &sqlstorage.Storage{DB: sqlxDB}
	event := storage.Event{
		ID:          "1",
		Title:       "Updated Event",
		Description: "This is an updated event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
	}

	mock.ExpectExec("UPDATE events").
		WithArgs(event.Title, event.Description, event.StartTime, event.EndTime, event.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storages.UpdateEvent(event.ID, event)
	assert.NoError(t, err)
}

func TestDeleteEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	storage := &sqlstorage.Storage{DB: sqlxDB}
	mock.ExpectExec("DELETE FROM events").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = storage.DeleteEvent("1")
	assert.NoError(t, err)
}

func TestGetEventFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	storages := &sqlstorage.Storage{DB: sqlxDB}
	event := storage.Event{
		ID:          "1",
		Title:       "Found Event",
		Description: "This event is found",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
	}

	mock.ExpectQuery("SELECT id, title, description, start_time, end_time FROM events WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(mock.NewRows([]string{"id", "title", "description", "start_time", "end_time"}).
			AddRow(event.ID, event.Title, event.Description, event.StartTime, event.EndTime))

	retEvent, err := storages.GetEvent("1")
	assert.NoError(t, err)
	assert.Equal(t, event, retEvent)
}

func TestGetEventNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	storages := &sqlstorage.Storage{DB: sqlxDB}
	mock.ExpectQuery("SELECT id, title, description, start_time, end_time FROM events WHERE id = \\$1").
		WithArgs("2").
		WillReturnError(errors.New("event not found"))
	_, err = storages.GetEvent("2")
	assert.Error(t, err)
	assert.Equal(t, "event not found", err.Error())
}

func TestListEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	storages := &sqlstorage.Storage{DB: sqlxDB}
	events := []storage.Event{
		{
			ID:          "1",
			Title:       "Event 1",
			Description: "Description 1",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(1 * time.Hour),
		},
		{
			ID:          "2",
			Title:       "Event 2",
			Description: "Description 2",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(1 * time.Hour),
		},
	}

	mock.ExpectQuery("SELECT id, title, description, start_time, end_time FROM events").
		WillReturnRows(mock.NewRows([]string{"id", "title", "description", "start_time", "end_time"}).
			AddRow(events[0].ID, events[0].Title, events[0].Description, events[0].StartTime, events[0].EndTime).
			AddRow(events[1].ID, events[1].Title, events[1].Description, events[1].StartTime, events[1].EndTime))

	retEvents, err := storages.ListEvents()
	assert.NoError(t, err)
	assert.Equal(t, events, retEvents)
}
