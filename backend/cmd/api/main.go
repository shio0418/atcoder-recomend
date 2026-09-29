package main

import (
	"atcoder-recomend/internal/client"
	"atcoder-recomend/internal/model"
	"atcoder-recomend/internal/recommend"
	"atcoder-recomend/internal/repository"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type recommendResponse struct {
	Review []model.SolvedProblem `json:"review"`
	New    []model.Problem       `json:"new"`
}

func recommendHandler(conn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		rating, err := client.FetchRating("uwaaaaaaa")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "rating取得失敗: " + err.Error()})
		}

		tagStats, err := repository.GetTagStats(ctx, conn)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "tagStats取得失敗: " + err.Error()})
		}

		solvedProblems, err := repository.GetSolvedProblems(ctx, conn)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "solvedProblems取得失敗: " + err.Error()})
		}

		candidates, err := repository.ListCandidateProblems(ctx, conn)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "candidates取得失敗: " + err.Error()})
		}

		reviewCandidates := recommend.ReviewCandidates(solvedProblems, time.Now())
		newCandidates := recommend.NewCandidates(tagStats, rating, candidates)

		return c.JSON(http.StatusOK, recommendResponse{
			Review: reviewCandidates,
			New:    newCandidates,
		})
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env読み込み失敗: %v", err)
	}
	ctx := context.Background()
	conn, err := repository.Connect(ctx)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer conn.Close(ctx)

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/api/recommend", recommendHandler(conn))
	e.Logger.Fatal(e.Start(":8080"))
}
