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
	ServiceID string
	Level     string
	Message   string
	Metadata  []byte
}

// InsertLog inserts a log entry into the timeseries DB.
func (db *DB) InsertLog(ctx context.Context, log LogEntry) error {
	_, err := db.Pool.Exec(ctx,
		"INSERT INTO log_entries (time, service_id, level, message, metadata) VALUES ($1, $2, $3, $4, $5)",
		log.Time, log.ServiceID, log.Level, log.Message, log.Metadata)
	return err
}

// MetricPoint struct
type MetricPoint struct {
	Time       time.Time
	ServiceID  string
	MetricName string
	Value      float64
	Labels     []byte
}

// InsertMetric inserts a metric data point into the timeseries DB.
func (db *DB) InsertMetric(ctx context.Context, metric MetricPoint) error {
	_, err := db.Pool.Exec(ctx,
		"INSERT INTO metric_points (time, service_id, metric_name, value, labels) VALUES ($1, $2, $3, $4, $5)",
		metric.Time, metric.ServiceID, metric.MetricName, metric.Value, metric.Labels)
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
