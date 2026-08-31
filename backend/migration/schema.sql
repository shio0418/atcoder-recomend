CREATE TABLE problems(
    id SERIAL PRIMARY KEY,
    problem_id TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    difficulty INTEGER,
    tags TEXT[]
);

CREATE TABLE submissions(
    id SERIAL PRIMARY KEY,
    submission_id BIGINT NOT NULL UNIQUE,
    problem_id TEXT NOT NULL REFERENCES problems(problem_id),
    result TEXT NOT NULL,
    code TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX submissions_problem_id_idx ON submissions(problem_id)
