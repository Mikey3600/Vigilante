package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/user/vigilante/internal/storage"
	"google.golang.org/api/option"
)

// Client wraps the Gemini API interface.
type Client struct {
	cli *genai.Client
}

// NewClient initializes the Gemini client mapping it to the provided API key.
func NewClient(ctx context.Context) (*Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to init genai client: %w", err)
	}

	return &Client{cli: client}, nil
}

// RootCauseReport maps the desired JSON format from Gemini.
type RootCauseReport struct {
	Summary      string `json:"summary"`
	LikelyCause  string `json:"likely_cause"`
	SuggestedFix string `json:"suggested_fix"`
}

// AnalyzeLogs asks Gemini 1.5 Flash to determine the root cause of an anomaly based on recent logs.
func (c *Client) AnalyzeLogs(ctx context.Context, logs []storage.LogEntry, anomalyMeta string) (*RootCauseReport, error) {
	model := c.cli.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"

	prompt := fmt.Sprintf(`You are a site reliability engineer diagnosing an issue.
Anomaly: %s
Logs Context:
%v
Generate a JSON object with 'summary', 'likely_cause', and 'suggested_fix'.
`, anomalyMeta, logs)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response candidates generated")
	}

	var report RootCauseReport
	part := resp.Candidates[0].Content.Parts[0]
	
	switch p := part.(type) {
	case genai.Text:
		if err := json.Unmarshal([]byte(p), &report); err != nil {
			return nil, err
		}
	}

	return &report, nil
}
