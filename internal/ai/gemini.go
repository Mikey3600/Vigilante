package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/user/vigilante/internal/storage"
)

type Client struct{}

type RootCauseReport struct {
	Summary      string `json:"summary"`
	LikelyCause  string `json:"likely_cause"`
	SuggestedFix string `json:"suggested_fix"`
}

func NewClient(ctx context.Context) (*Client, error) {
	if os.Getenv("GROQ_API_KEY") == "" {
		return nil, fmt.Errorf("GROQ_API_KEY is not set")
	}
	return &Client{}, nil
}

func (c *Client) AnalyzeLogs(ctx context.Context, logs []storage.LogEntry, anomalyMeta string) (*RootCauseReport, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	url := "https://api.groq.com/openai/v1/chat/completions"

	prompt := fmt.Sprintf(`You are a site reliability engineer diagnosing an incident.
Anomaly: %s
Recent logs: %v
Respond ONLY with a valid JSON object with exactly these three fields: summary, likely_cause, suggested_fix. No markdown, no backticks, just raw JSON.`, anomalyMeta, logs)

	body := map[string]any{
		"model": "llama-3.1-8b-instant",
		"messages": []map[string]any{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
	}

	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	log.Printf("GROQ RAW: %s", string(data))

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices: %s", string(data))
	}

	text := result.Choices[0].Message.Content
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var report RootCauseReport
	if err := json.Unmarshal([]byte(text), &report); err != nil {
		report = RootCauseReport{
			Summary:      text,
			LikelyCause:  "See summary",
			SuggestedFix: "See summary",
		}
	}
	return &report, nil
}