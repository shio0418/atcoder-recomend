-- +goose Up
ALTER TABLE problems ADD COLUMN statement TEXT;
ALTER TABLE problems ADD COLUMN editorial_text TEXT;

-- +goose Down
ALTER TABLE problems DROP COLUMN statement;
ALTER TABLE problems DROP COLUMN editorial_text;
