package ingestion

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/user/vigilante/internal/storage"
)

<<<<<<< HEAD
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
=======
type LogPayload struct {
	Time     string          `json:"time"`
	Level    string          `json:"level" binding:"required"`
	Message  string          `json:"message" binding:"required"`
	Metadata json.RawMessage `json:"metadata"`
}

func (l *LogPayload) UnmarshalJSON(data []byte) error {
	type alias LogPayload
	tmp := struct {
		*alias
		Timestamp *string `json:"timestamp"`
	}{
		alias: (*alias)(l),
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if (l.Time == "") && tmp.Timestamp != nil {
		l.Time = *tmp.Timestamp
	}

	return nil
}

type MetricPayload struct {
	Time       string          `json:"time" binding:"required"`
	MetricName string          `json:"metric_name" binding:"required"`
	Value      float64         `json:"value"`
	Labels     json.RawMessage `json:"labels"`
}

type BatchLogPayload struct {
	ServiceID string       `json:"service_id" binding:"required"`
	Logs      []LogPayload `json:"logs" binding:"required,min=1"`
}

type BatchMetricPayload struct {
	ServiceID string          `json:"service_id" binding:"required"`
	Metrics   []MetricPayload `json:"metrics" binding:"required,min=1"`
>>>>>>> 1d7e6cff7174607f532472fa8ffd51a5c3841fb3
}

func ProcessLogs(ctx context.Context, db *storage.DB, tenantID string, payload BatchLogPayload) error {
	if payload.ServiceID == "" {
		return errors.New("service_id required")
	}
<<<<<<< HEAD
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
=======

	for _, l := range payload.Logs {
		t := time.Now()
		if l.Time != "" {
			parsed, err := time.Parse(time.RFC3339, l.Time)
			if err != nil {
				return err
			}
			t = parsed
		}

>>>>>>> 1d7e6cff7174607f532472fa8ffd51a5c3841fb3
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
<<<<<<< HEAD
=======

>>>>>>> 1d7e6cff7174607f532472fa8ffd51a5c3841fb3
	return nil
}

func ProcessMetrics(ctx context.Context, db *storage.DB, tenantID string, payload BatchMetricPayload) error {
	if payload.ServiceID == "" {
		return errors.New("service_id required")
	}
<<<<<<< HEAD
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
=======

	for _, m := range payload.Metrics {
		t, err := time.Parse(time.RFC3339, m.Time)
		if err != nil {
			return err
		}
		if err := db.InsertMetric(ctx, storage.MetricPoint{Time: t, TenantID: tenantID, ServiceID: payload.ServiceID, MetricName: m.MetricName, Value: m.Value, Labels: m.Labels}); err != nil {
>>>>>>> 1d7e6cff7174607f532472fa8ffd51a5c3841fb3
			return err
		}
	}
	return nil
<<<<<<< HEAD
}
=======
}
>>>>>>> 1d7e6cff7174607f532472fa8ffd51a5c3841fb3
