package classifier

import (
	"atcoder-recomend/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type jevRequest struct {
	Model     string                       `json:"model"`
	State     map[string]string            `json:"state"`
	Questions map[string]jevChoiceQuestion `json:"questions"`
}

type jevChoiceQuestion struct {
	Type         string            `json:"type"`
	Instructions string            `json:"instructions"`
	Criteria     map[string]string `json:"criteria"`
}

type jevResponse struct {
	Answers map[string]struct {
		Choice        string             `json:"choice"`
		Probabilities map[string]float64 `json:"probabilities"`
	} `json:"answers"`
}

const endpoint = "https://api.typesafe.ai/v1/systemone"
const probThreshold = 0.05

func buildCriteria() map[string]string {
	criteria := make(map[string]string, len(TagCandidates))
	for _, tag := range TagCandidates {
		criteria[tag] = tag
	}
	return criteria
}

var jevCriteria = buildCriteria()

func buildJevRequest(p model.Problem) jevRequest {
	state := map[string]string{
		"title": p.Title,
		"statement": *p.Statement,
	}
	if p.HasEditorial {
		state["editorial_text"] = *p.EditorialText
	}
	questions := map[string]jevChoiceQuestion{
		"tags": jevChoiceQuestion{
			Type:         "choice",
			Instructions: "この競技プログラミングの問題を解くために主に必要となるアルゴリズム・データ構造の分野を、選択肢の中から選んでください。複数の分野が関係する場合でも、最も本質的な1つを選んでください。当てはまるものが無い場合は、無理に選ばず確信度を低くしてください。",
			Criteria:     jevCriteria,
		},
	}
	return jevRequest{
		Model:     "jev-latest",
		State:     state,
		Questions: questions,
	}
}

func callJev(ctx context.Context, req jevRequest) (jevResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return jevResponse{}, err
	}
	HttpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return jevResponse{}, err
	}

	APIKey := os.Getenv("JEV_API_KEY")
	HttpReq.Header.Set("Authorization", "Bearer "+APIKey)
	HttpReq.Header.Set("Content-Type", "application/json")

	const maxAttempts = 5
	backoff := time.Second
	var resp *http.Response
	var result jevResponse
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}
		resp, err = http.DefaultClient.Do(HttpReq)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			err = fmt.Errorf("jev status %d: %s", resp.StatusCode, respBody)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
				return jevResponse{}, err
			}
			continue
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		if err == nil {
			break
		}
	}
	if err != nil {
		return jevResponse{}, err
	}
	return result, nil
}

func classifyOne(ctx context.Context, p model.Problem) (ClassificationResult, error) {
	req := buildJevRequest(p)

	resp, err := callJev(ctx, req)
	if err != nil {
		return ClassificationResult{}, err
	}

	probabilities := resp.Answers["tags"].Probabilities
	tags := selectTags(probabilities)

	return ClassificationResult{
		ProblemID:     p.ProblemID,
		Tags:          tags,
		Probabilities: probabilities,
	}, nil
}


func selectTags(probabilities map[string]float64) []string {
	var result []string
	for tag, prob := range probabilities {
		if prob >= probThreshold {
			result = append(result, tag)
		}
	}
	return result
}

func ClassifyBatch(ctx context.Context, problems []model.Problem) ([]ClassificationResult, error) {
	var results []ClassificationResult
	for _, p := range problems {
		result, err := classifyOne(ctx, p)
		if err != nil {
			log.Printf("タグ分類失敗")
			continue
		}
		results = append(results, result)
	}
	return results, nil
}
