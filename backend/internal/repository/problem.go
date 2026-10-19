package repository

import (
	"atcoder-recomend/internal/model"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

func InsertProblem(ctx context.Context, conn DBTX, p model.Problem) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO problems (problem_id, title, difficulty, tags, contest_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (problem_id) DO UPDATE SET contest_id = EXCLUDED.contest_id`,
		p.ProblemID, p.Title, p.Difficulty, p.Tags, p.ContestID,
	)
	return err
}

func ListUntaggedProblems(ctx context.Context, conn DBTX, limit int) ([]model.Problem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problem_id, title, difficulty, contest_id, statement, editorial_text, has_editorial
		FROM problems WHERE tags IS NULL AND statement IS NOT NULL LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Problem])
}

func ListUnembeddedProblems(ctx context.Context, conn DBTX, limit int) ([]model.Problem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problem_id, title, difficulty, contest_id, statement, editorial_text, has_editorial
		FROM problems WHERE embedding IS NULL AND statement IS NOT NULL LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Problem])
}

func UpdateProblemsEmbedding(ctx context.Context, conn DBTX, problemID string, values []float32, embedding_model string) error {
	embedding := pgvector.NewVector(values)
	_, err := conn.Exec(ctx,
		`UPDATE problems SET embedding = $1, embedding_model = $2 WHERE problem_id = $3`,
		embedding, embedding_model, problemID,
	)
	return err
}

// ListUnscrapedProblems returns problems that still need their statement
// (and editorial, if any) scraped. Separate from ListUntaggedProblems so a
// restart after a crash mid-scrape doesn't re-scrape problems that already
// succeeded — only tags are missing for those.
func ListUnscrapedProblems(ctx context.Context, conn DBTX, limit int) ([]model.Problem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problem_id, contest_id FROM problems WHERE tags IS NULL AND statement IS NULL LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Problem])
}

func UpdateProblemTags(ctx context.Context, conn DBTX, problemID string, tags []string, probabilities map[string]float64) error {
	_, err := conn.Exec(ctx,
		`UPDATE problems SET tags = $1, tag_probabilities = $2 WHERE problem_id = $3`,
		tags, probabilities, problemID,
	)
	return err
}

func UpdateProblemScrapedContent(ctx context.Context, conn DBTX, problemID string, statement string, editorialText *string, hasEditorial bool) error {
	_, err := conn.Exec(ctx,
		`UPDATE problems SET statement = $1, editorial_text = $2, has_editorial = $3 WHERE problem_id = $4`,
		statement, editorialText, hasEditorial, problemID,
	)
	return err
}

func ListCandidateProblems(ctx context.Context, conn DBTX) ([]model.Problem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problems.problem_id, problems.title, problems.difficulty, problems.tags, problems.contest_id
		FROM problems
		LEFT JOIN submissions ON problems.problem_id = submissions.problem_id AND submissions.result = 'AC'
		WHERE problems.tags IS NOT NULL
		AND problems.difficulty IS NOT NULL
		AND submissions.problem_id IS NULL`,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Problem])
}
