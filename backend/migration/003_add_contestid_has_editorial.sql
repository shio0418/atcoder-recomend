-- +goose Up
ALTER TABLE problems ADD COLUMN contest_id TEXT;
ALTER TABLE problems ADD COLUMN has_editorial BOOLEAN;

-- +goose Down
ALTER TABLE problems DROP COLUMN contest_id;
ALTER TABLE problems DROP COLUMN has_editorial;
