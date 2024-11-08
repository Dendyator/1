-- +goose Up
CREATE TABLE IF NOT EXISTS events (
                                      id SERIAL PRIMARY KEY,
                                      title VARCHAR(255) NOT NULL,
                                      description TEXT,
                                      start_time TIMESTAMP NOT NULL,
                                      end_time TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS events;
