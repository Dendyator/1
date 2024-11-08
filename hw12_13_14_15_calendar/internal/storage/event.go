package storage

import "time"

type Event struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
}

type Interface interface {
	CreateEvent(event Event) error
	UpdateEvent(id string, newEvent Event) error
	DeleteEvent(id string) error
	GetEvent(id string) (Event, error)
	ListEvents() ([]Event, error)
	DeleteOldEvents(before time.Time) error
}
