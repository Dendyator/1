package sqlstorage

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" //nolint
	"github.com/jmoiron/sqlx"          //nolint

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint
)

type Storage struct {
	DB *sqlx.DB
}

func New(dsn string) (*Storage, error) {
	log.Println("Using DSN:", dsn)
	var db *sqlx.DB
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Open("postgres", dsn)

		if err == nil {
			if err = db.Ping(); err == nil {
				log.Println("Successfully connected to the database!")
				return &Storage{DB: db}, nil
			}
		}
		log.Printf("Failed to connect to database: %v. Retrying...\n", err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

func (s *Storage) CreateEvent(event storage.Event) error {
	query := `INSERT INTO events (id, title, description, start_time, end_time) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.DB.Exec(query, event.ID, event.Title, event.Description, event.StartTime, event.EndTime)
	return err
}

func (s *Storage) UpdateEvent(id string, newEvent storage.Event) error {
	query := `UPDATE events SET title = $1, description = $2, start_time = $3, end_time = $4 WHERE id = $5`
	res, err := s.DB.Exec(query, newEvent.Title, newEvent.Description, newEvent.StartTime, newEvent.EndTime, id)
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
	res, err := s.DB.Exec(query, id)
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
	err := s.DB.Get(&event, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return event, errors.New("event not found")
	}
	return event, err
}

func (s *Storage) ListEvents() ([]storage.Event, error) {
	var events []storage.Event
	query := `SELECT id, title, description, start_time, end_time FROM events`
	err := s.DB.Select(&events, query)
	return events, err
}

func (s *Storage) DeleteOldEvents(before time.Time) error {
	query := `DELETE FROM events WHERE end_time < $1`
	_, err := s.DB.Exec(query, before)
	return err
}
