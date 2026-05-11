package ingestion

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/user/vigilante/internal/storage"
)

type BatchLogPayload struct {
	ServiceID string `json:"service_id" binding:"required"`
	Logs      []struct {
		Time     string          `json:"time"`
		Level    string          `json:"level" binding:"required"`
		Message  string          `json:"message" binding:"required"`
		Metadata json.RawMessage `json:"metadata"`
	} `json:"logs" binding:"required,min=1"`
}

type BatchMetricPayload struct {
	ServiceID string `json:"service_id" binding:"required"`
	Metrics   []struct {
		Time       string          `json:"time"`
		MetricName string          `json:"metric_name" binding:"required"`
		Value      float64         `json:"value"`
		Labels     json.RawMessage `json:"labels"`
	} `json:"metrics" binding:"required,min=1"`
}

func ProcessLogs(ctx context.Context, db *storage.DB, tenantID string, payload BatchLogPayload) error {
	if payload.ServiceID == "" {
		return errors.New("service_id required")
	}
	for _, l := range payload.Logs {
		var t time.Time
		if l.Time == "" {
			t = time.Now()
		} else {
			var err error
			t, err = time.Parse(time.RFC3339, l.Time)
			if err != nil {
				return err
			}
		}
		if err := db.InsertLog(ctx, storage.LogEntry{
			Time:      t,
			TenantID:  tenantID,
			ServiceID: payload.ServiceID,
			Level:     l.Level,
			Message:   l.Message,
			Metadata:  l.Metadata,
		}); err != nil {
			return err
		}
	}
	return nil
}

func ProcessMetrics(ctx context.Context, db *storage.DB, tenantID string, payload BatchMetricPayload) error {
	if payload.ServiceID == "" {
		return errors.New("service_id required")
	}
	for _, m := range payload.Metrics {
		var t time.Time
		if m.Time == "" {
			t = time.Now()
		} else {
			var err error
			t, err = time.Parse(time.RFC3339, m.Time)
			if err != nil {
				return err
			}
		}
		if err := db.InsertMetric(ctx, storage.MetricPoint{
			Time:       t,
			TenantID:   tenantID,
			ServiceID:  payload.ServiceID,
			MetricName: m.MetricName,
			Value:      m.Value,
			Labels:     m.Labels,
		}); err != nil {
			return err
		}
	}
	return nil
}