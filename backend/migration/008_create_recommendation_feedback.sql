-- +goose Up
CREATE TABLE IF NOT EXISTS recommendation_feedback (
    id SERIAL PRIMARY KEY,
    problem_id TEXT NOT NULL REFERENCES problems(problem_id),
    source TEXT NOT NULL,   -- 'review' | 'new' | 'cf' | 'search'
    feedback TEXT NOT NULL, -- 'solved' | 'needed_help'
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS recommendation_feedback;