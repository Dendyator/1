package memorystorage

import (
	"errors"
	"sync"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; exists {
		return errors.New("event already exists")
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(id string, newEvent storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return errors.New("event not found")
	}
	s.events[id] = newEvent
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return errors.New("event not found")
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	if !exists {
		return event, errors.New("event not found")
	}
	return event, nil
}

func (s *Storage) ListEvents() ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}
	return events, nil
}
