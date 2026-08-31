package classifier

import (
	"atcoder-recomend/internal/model"
	"context"
	"encoding/json"
	"os"

	"google.golang.org/genai"
)

func ClassifyBatch(ctx context.Context, problems []model.Problem) ([]ClassificationResult, error) {
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

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return nil, err
	}

	var result []ClassificationResult
	if err := json.Unmarshal([]byte(resp.Text()), &result); err != nil {
		return nil, err
	}
	return result, nil
}
