-- +goose Up
ALTER TABLE problems ADD COLUMN tag_probabilities JSONB;

-- +goose Down
ALTER TABLE problems DROP COLUMN tag_probabilities;
