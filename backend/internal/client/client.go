package client

import (
	"atcoder-recomend/internal/model"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

func FetchSubmissions(user string, fromSecond int64) ([]SubmissionResponse, error) {
	base, err := url.Parse("https://kenkoooo.com/atcoder/atcoder-api/v3/user/submissions")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("user", user)
	params.Add("from_second", strconv.FormatInt(fromSecond, 10))
	base.RawQuery = params.Encode()
	finalURL := base.String()

	resp, err := http.Get(finalURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var submissions []SubmissionResponse

	if err := json.NewDecoder(resp.Body).Decode(&submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}

func FetchProblems() ([]model.Problem, error) {
	problemUrl := "https://kenkoooo.com/atcoder/resources/problems.json"
	problemResp, err := http.Get(problemUrl)
	if err != nil {
		return nil, err
	}
	defer problemResp.Body.Close()
	var problems []ProblemResponse
	if err := json.NewDecoder(problemResp.Body).Decode(&problems); err != nil {
		return nil, err
	}

	problemModelUrl := "https://kenkoooo.com/atcoder/resources/problem-models.json"
	problemModelResp, err := http.Get(problemModelUrl)
	if err != nil {
		return nil, err
	}
	defer problemModelResp.Body.Close()

	var models map[string]ProblemModel
	if err := json.NewDecoder(problemModelResp.Body).Decode(&models); err != nil {
		return nil, err
	}

	result := make([]model.Problem, 0, len(problems))
	for _, p := range problems {
		difficulty := models[p.ID].Difficulty
		result = append(result, model.Problem{
			ProblemID:  p.ID,
			Title:      p.Title,
			Difficulty: difficulty,
		})
	}
	return result, nil
}
