-- +goose Up
ALTER TABLE feeds
ADD last_fethced_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fethced_at;