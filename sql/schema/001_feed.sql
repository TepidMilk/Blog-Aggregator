-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID ON DELETE CASCADE,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
);

-- +goose Down
DROP TABLE feeds;