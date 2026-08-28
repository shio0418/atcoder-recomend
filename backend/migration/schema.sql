CREATE TABLE problems(
    id SERIAL PRIMARY KEY,
    problem_id TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    difficulty INTEGER,
    tags TEXT[]
);

CREATE TABLE submissions(
    id SERIAL PRIMARY KEY,
    problem_id TEXT NOT NULL references problems(problem_id),
    result TEXT NOT NULL,
    code TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now()
);