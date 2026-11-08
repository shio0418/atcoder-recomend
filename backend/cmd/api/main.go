package main

import (
	"atcoder-recomend/internal/cf"
	"atcoder-recomend/internal/client"
	"atcoder-recomend/internal/embedding"
	"atcoder-recomend/internal/model"
	"atcoder-recomend/internal/recommend"
	"atcoder-recomend/internal/repository"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pgvector/pgvector-go"
)

type recommendResponse struct {
	Review []model.SolvedProblem `json:"review"`
	New    []model.Problem       `json:"new"`
}

func recommendHandler(conn repository.DBTX) echo.HandlerFunc {
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

		baseInterval, _ := strconv.ParseFloat(c.QueryParam("base_interval"), 64)
		if baseInterval == 0 {
			baseInterval = 5
		}
		growthRate, _ := strconv.ParseFloat(c.QueryParam("growth_rate"), 64)
		if growthRate == 0 {
			growthRate = 2
		}
		difficultyWeight, _ := strconv.ParseFloat(c.QueryParam("difficulty_weight"), 64)
		if difficultyWeight == 0 {
			difficultyWeight = 1000
		}
		overdueCap, _ := strconv.ParseFloat(c.QueryParam("overdue_cap"), 64)
		if overdueCap == 0 {
			overdueCap = 90
		}
		reviewCandidates := recommend.ReviewCandidates(solvedProblems, time.Now(), baseInterval, growthRate, difficultyWeight, overdueCap)
		newCandidates := recommend.NewCandidates(tagStats, rating, candidates)

		return c.JSON(http.StatusOK, recommendResponse{
			Review: reviewCandidates,
			New:    newCandidates,
		})
	}
}

func similarHandler(conn repository.DBTX) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		problemID := c.QueryParam("problem_id")
		if problemID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "problem_id is required"})
		}

		similar, err := repository.FindSimilarProblems(ctx, conn, problemID, 5)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "類似問題取得失敗"})
		}
		return c.JSON(http.StatusOK, similar)
	}
}

func similarByTagHandler(conn repository.DBTX) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		problemID := c.QueryParam("problem_id")
		if problemID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "problem_id is required"})
		}

		similar, err := repository.FindSimilarProblemsByTag(ctx, conn, problemID, 5)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タグによる類似問題取得失敗"})
		}
		return c.JSON(http.StatusOK, similar)
	}
}

func searchProblemHandler(conn repository.DBTX) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		query := c.QueryParam("q")
		if query == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "q is required"})
		}
		var difficultyMin, difficultyMax *int

		if v := c.QueryParam("difficulty_min"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "difficulty_minは整数で指定してください"})
			}
			difficultyMin = &n
		}
		if v := c.QueryParam("difficulty_max"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "difficulty_maxは整数で指定してください"})
			}
			difficultyMax = &n
		}

		embedResp, err := embedding.EmbedBatch(ctx, map[string]string{"query": "query " + query})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "クエリの埋め込みに失敗"})
		}

		queryVector := pgvector.NewVector(embedResp["query"].Values)

		results, err := repository.SearchProblemsByEmbedding(ctx, conn, queryVector, difficultyMin, difficultyMax, 5)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "テキストによる検索に失敗"})
		}
		return c.JSON(http.StatusOK, results)
	}
}

func cfRecommendHandler(conn repository.DBTX) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		problemIDs, err := cf.Recommend(ctx, "uwaaaaaaa", 10)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "CF推薦取得失敗"})
		}
		return c.JSON(http.StatusOK, problemIDs)
	}
}

func main() {
	_ = godotenv.Load()
	ctx := context.Background()
	pool, err := repository.ConnectPool(ctx)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer pool.Close()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
	}))

	e.GET("/api/recommend", recommendHandler(pool))
	e.GET("/api/similar", similarHandler(pool))
	e.GET("/api/similarByTag", similarByTagHandler(pool))
	e.GET("/api/search", searchProblemHandler(pool))
	e.GET("/api/cfRecommend", cfRecommendHandler(pool))
	e.Logger.Fatal(e.Start(":8080"))
}
