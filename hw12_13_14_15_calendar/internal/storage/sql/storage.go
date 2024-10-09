package sqlstorage

import (
	"database/sql"
	"errors"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/jmoiron/sqlx"                                        //nolint:depguard
)

type Storage struct {
	db *sqlx.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) CreateEvent(event storage.Event) error {
	query := `INSERT INTO events (id, title, description, start_time, end_time) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(query, event.ID, event.Title, event.Description, event.StartTime, event.EndTime)
	return err
}

func (s *Storage) UpdateEvent(id string, newEvent storage.Event) error {
	query := `UPDATE events SET title = $1, description = $2, start_time = $3, end_time = $4 WHERE id = $5`
	res, err := s.db.Exec(query, newEvent.Title, newEvent.Description, newEvent.StartTime, newEvent.EndTime, id)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return errors.New("event not found")
	}

	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	query := `DELETE FROM events WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return errors.New("event not found")
	}

	return nil
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	var event storage.Event
	query := `SELECT id, title, description, start_time, end_time FROM events WHERE id = $1`
	err := s.db.Get(&event, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return event, errors.New("event not found")
	}
	return event, err
}

func (s *Storage) ListEvents() ([]storage.Event, error) {
	var events []storage.Event
	query := `SELECT id, title, description, start_time, end_time FROM events`
	err := s.db.Select(&events, query)
	return events, err
}
