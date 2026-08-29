package client

import (
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
