package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents a database connection pool.
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB creates a new database connection pool.
func NewDB(ctx context.Context, connString string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool.
func (db *DB) Close() {
	db.Pool.Close()
}

// RunMigrations runs schema.sql against the database.
func (db *DB) RunMigrations(ctx context.Context) error {
	schema, err := os.ReadFile("internal/storage/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	_, err = db.Pool.Exec(ctx, string(schema))
	return err
}

// LogEntry struct
type LogEntry struct {
	Time      time.Time
	TenantID  string
	ServiceID string
	Level     string
	Message   string
	Metadata  []byte
}

// InsertLog inserts a log entry into the timeseries DB.
func (db *DB) InsertLog(ctx context.Context, log LogEntry) error {
	_, err := db.Pool.Exec(ctx,
		"INSERT INTO log_entries (time, tenant_id, service_id, level, message, metadata) VALUES ($1, $2, $3, $4, $5, $6)",
		log.Time, log.TenantID, log.ServiceID, log.Level, log.Message, log.Metadata)
	return err
}

// MetricPoint struct
type MetricPoint struct {
	Time       time.Time
	TenantID   string
	ServiceID  string
	MetricName string
	Value      float64
	Labels     []byte
}

// InsertMetric inserts a metric data point into the timeseries DB.
func (db *DB) InsertMetric(ctx context.Context, metric MetricPoint) error {
	_, err := db.Pool.Exec(ctx,
		"INSERT INTO metric_points (time, tenant_id, service_id, metric_name, value, labels) VALUES ($1, $2, $3, $4, $5, $6)",
		metric.Time, metric.TenantID, metric.ServiceID, metric.MetricName, metric.Value, metric.Labels)
	return err
}

// Anomaly struct
type Anomaly struct {
	ID               string
	ServiceID        string
	DetectedAt       time.Time
	AnomalyType      string
	Description      string
	RootCauseSummary string
	LikelyCause      string
	SuggestedFix     string
}

// InsertAnomaly stores a detected anomaly and AI summary.
func (db *DB) InsertAnomaly(ctx context.Context, anomaly Anomaly) error {
	_, err := db.Pool.Exec(ctx,
		`INSERT INTO anomalies (service_id, detected_at, anomaly_type, description, root_cause_summary, likely_cause, suggested_fix) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		anomaly.ServiceID, anomaly.DetectedAt, anomaly.AnomalyType, anomaly.Description, anomaly.RootCauseSummary, anomaly.LikelyCause, anomaly.SuggestedFix)
	return err
}

// GetRecentLogs retrieves recent logs for AI analysis.
func (db *DB) GetRecentLogs(ctx context.Context, serviceID string, limit int) ([]LogEntry, error) {
	rows, err := db.Pool.Query(ctx, 
		"SELECT time, service_id, level, message, metadata FROM log_entries WHERE service_id = $1 ORDER BY time DESC LIMIT $2", 
		serviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.Time, &l.ServiceID, &l.Level, &l.Message, &l.Metadata); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (db *DB) GetRecentLogsForTenant(ctx context.Context, tenantID string, limit int) ([]LogEntry, error) {
	rows, err := db.Pool.Query(ctx, 
		`SELECT time, service_id, level, message, metadata 
		FROM log_entries 
		WHERE tenant_id = $1 
		ORDER BY time DESC LIMIT $2`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.Time, &l.ServiceID, &l.Level, &l.Message, &l.Metadata); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (db *DB) GetRecentMetricsForTenant(ctx context.Context, tenantID string, limit int) ([]MetricPoint, error) {
	rows, err := db.Pool.Query(ctx, 
		`SELECT time, service_id, metric_name, value, labels 
		FROM metric_points 
		WHERE tenant_id = $1 
		ORDER BY time DESC LIMIT $2`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []MetricPoint
	for rows.Next() {
		var m MetricPoint
		if err := rows.Scan(&m.Time, &m.ServiceID, &m.MetricName, &m.Value, &m.Labels); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (db *DB) GetRecentAnomaliesForTenant(ctx context.Context, tenantID string, limit int) ([]Anomaly, error) {
	rows, err := db.Pool.Query(ctx, 
		`SELECT id, service_id, detected_at, anomaly_type, description, root_cause_summary, likely_cause, suggested_fix 
		FROM anomalies 
		WHERE tenant_id = $1 
		ORDER BY detected_at DESC LIMIT $2`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anomalies []Anomaly
	for rows.Next() {
		var a Anomaly
		if err := rows.Scan(&a.ID, &a.ServiceID, &a.DetectedAt, &a.AnomalyType, &a.Description, &a.RootCauseSummary, &a.LikelyCause, &a.SuggestedFix); err != nil {
			return nil, err
		}
		anomalies = append(anomalies, a)
	}
	return anomalies, nil
}
