-- +goose Up
CREATE EXTENSION IF NOT EXISTS vector;
ALTER TABLE problems ADD COLUMN embedding vector(384);
ALTER TABLE problems ADD COLUMN embedding_model TEXT;

CREATE TABLE IF NOT EXISTS inference_calls (
    id SERIAL PRIMARY KEY,
    call_type TEXT NOT NULL,
    latency_ms INTEGER NOT NULL,
    success BOOLEAN NOT NULL,
    item_count INTEGER,
    token_count INTEGER,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS inference_calls;
ALTER TABLE problems DROP COLUMN embedding vector(384);
ALTER TABLE problems DROP COLUMN embedding_model TEXT;