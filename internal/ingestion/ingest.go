package ingestion

import (
	"context"
	"encoding/json"
	"time"

	"github.com/user/vigilante/internal/storage"
)

// BatchLogPayload maps the incoming bulk log struct.
type BatchLogPayload struct {
	ServiceID string `json:"service_id"`
	Logs      []struct {
		Time     string          `json:"time"`
		Level    string          `json:"level"`
		Message  string          `json:"message"`
		Metadata json.RawMessage `json:"metadata"`
	} `json:"logs"`
}

// ProcessLogs validates and inserts batch logs.
func ProcessLogs(ctx context.Context, db *storage.DB, payload BatchLogPayload) error {
	for _, l := range payload.Logs {
		t, err := time.Parse(time.RFC3339, l.Time)
		if err != nil {
			t = time.Now()
		}
		entry := storage.LogEntry{
			Time:      t,
			ServiceID: payload.ServiceID,
			Level:     l.Level,
			Message:   l.Message,
			Metadata:  l.Metadata,
		}
		if err := db.InsertLog(ctx, entry); err != nil {
			return err
		}
	}
	return nil
}

// BatchMetricPayload maps bulk metrics.
type BatchMetricPayload struct {
	ServiceID string `json:"service_id"`
	Metrics   []struct {
		Time       string          `json:"time"`
		MetricName string          `json:"metric_name"`
		Value      float64         `json:"value"`
		Labels     json.RawMessage `json:"labels"`
	} `json:"metrics"`
}

// ProcessMetrics handles processing block metrics points.
func ProcessMetrics(ctx context.Context, db *storage.DB, payload BatchMetricPayload) error {
	for _, m := range payload.Metrics {
		t, err := time.Parse(time.RFC3339, m.Time)
		if err != nil {
			t = time.Now()
		}
		pt := storage.MetricPoint{
			Time:       t,
			ServiceID:  payload.ServiceID,
			MetricName: m.MetricName,
			Value:      m.Value,
			Labels:     m.Labels,
		}
		if err := db.InsertMetric(ctx, pt); err != nil {
			return err
		}
	}
	return nil
}
