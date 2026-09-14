-- +goose Up
ALTER TABLE submissions ADD COLUMN submission_id BIGINT NOT NULL UNIQUE;

-- +goose Down
ALTER TABLE submissions DROP COLUMN submission_id;
