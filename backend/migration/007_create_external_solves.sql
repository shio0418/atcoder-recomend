-- +goose Up
CREATE TABLE IF NOT EXISTS external_solves (
    id SERIAL PRIMARY KEY,
    user_handle TEXT NOT NULL,
    problem_id TEXT NOT NULL,
    UNIQUE(user_handle, problem_id)
);

-- +goose Down
DROP TABLE IF EXISTS external_solves;