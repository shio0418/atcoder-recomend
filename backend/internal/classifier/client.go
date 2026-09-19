package classifier

import (
	"atcoder-recomend/internal/model"
	"context"
	"encoding/json"
	"os"
	"time"

	"google.golang.org/genai"
)

func ClassifyBatchWithGemini(ctx context.Context, problems []model.Problem) ([]ClassificationResult, error) {
	prompt := BuildPrompt(problems)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"problem_id": {Type: genai.TypeString},
					"tags": {
						Type:  genai.TypeArray,
						Items: &genai.Schema{Type: genai.TypeString, Enum: TagCandidates},
					},
				},
				Required: []string{"problem_id", "tags"},
			},
		},
	}

	const maxAttempts = 5
	backoff := time.Second
	var resp *genai.GenerateContentResponse
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}
		resp, err = client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	var result []ClassificationResult
	if err := json.Unmarshal([]byte(resp.Text()), &result); err != nil {
		return nil, err
	}
	return result, nil
}
