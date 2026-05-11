package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/user/vigilante/internal/storage"
	"google.golang.org/api/option"
)

type Client struct{ cli *genai.Client }
type RootCauseReport struct { Summary string `json:"summary"`; LikelyCause string `json:"likely_cause"`; SuggestedFix string `json:"suggested_fix"` }

func NewClient(ctx context.Context) (*Client,error){ k:=os.Getenv("GEMINI_API_KEY"); if k==""{ return nil, fmt.Errorf("GEMINI_API_KEY is not set")}; c,err:=genai.NewClient(ctx, option.WithAPIKey(k)); if err!=nil{return nil,err}; return &Client{cli:c},nil }

func (c *Client) AnalyzeLogs(ctx context.Context, logs []storage.LogEntry, anomalyMeta string) (*RootCauseReport, error) {
	var last error
	for i:=0;i<3;i++ { reqCtx,cancel:=context.WithTimeout(ctx,30*time.Second); model:=c.cli.GenerativeModel("gemini-1.5-flash"); model.ResponseMIMEType="application/json"; resp,err:=model.GenerateContent(reqCtx, genai.Text(fmt.Sprintf("Anomaly: %s Logs: %+v",anomalyMeta,logs))); cancel(); if err!=nil{ last=err; time.Sleep(time.Duration(1<<i)*time.Second); continue }; if len(resp.Candidates)==0 { last=fmt.Errorf("no candidates"); continue }; txt,ok:=resp.Candidates[0].Content.Parts[0].(genai.Text); if !ok { last=fmt.Errorf("unexpected response"); continue }; var r RootCauseReport; if err:=json.Unmarshal([]byte(txt),&r); err!=nil{ last=err; continue }; return &r,nil }
	return &RootCauseReport{Summary:"AI analysis unavailable",LikelyCause:"unknown",SuggestedFix:"investigate logs manually"}, last
}
