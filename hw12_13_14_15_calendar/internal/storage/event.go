package storage

import "time"

type Event struct {
	ID          string
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type Interface interface {
	CreateEvent(event Event) error
	UpdateEvent(id string, newEvent Event) error
	DeleteEvent(id string) error
	GetEvent(id string) (Event, error)
	ListEvents() ([]Event, error)
}
