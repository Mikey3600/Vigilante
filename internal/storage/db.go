package storage

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

type DB struct{ Pool *pgxpool.Pool }

type LogEntry struct {
	Time                                time.Time
	TenantID, ServiceID, Level, Message string
	Metadata                            []byte
}
type MetricPoint struct {
	Time                            time.Time
	TenantID, ServiceID, MetricName string
	Value                           float64
	Labels                          []byte
}

type DashboardSnapshot struct {
	Latency   float64       `json:"latency"`
	CPU       float64       `json:"cpu"`
	Memory    float64       `json:"memory"`
	Metrics   []MetricPoint `json:"metrics"`
	Logs      []LogEntry    `json:"logs"`
	Anomalies []Anomaly     `json:"anomalies"`
}
type Anomaly struct {
	ID, TenantID, ServiceID, AnomalyType, Description, RootCauseSummary, LikelyCause, SuggestedFix, Severity string
	DetectedAt                                                                                               time.Time
}
type Service struct {
	ID, TenantID, Name string
	CreatedAt          time.Time
}

func NewDB(ctx context.Context, connString string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 25
	cfg.MinConns = 5
	cfg.MaxConnLifetime = time.Hour
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}
func (db *DB) Close()                         { db.Pool.Close() }
func (db *DB) Ping(ctx context.Context) error { return db.Pool.Ping(ctx) }
func (db *DB) RunMigrations(ctx context.Context) error {
	_, err := db.Pool.Exec(ctx, schemaSQL)
	return err
}
func (db *DB) InsertLog(ctx context.Context, l LogEntry) error {
	_, err := db.Pool.Exec(ctx, "INSERT INTO log_entries (time, tenant_id, service_id, level, message, metadata) VALUES ($1,$2,$3,$4,$5,$6)", l.Time, l.TenantID, l.ServiceID, l.Level, l.Message, l.Metadata)
	return err
}
func (db *DB) InsertMetric(ctx context.Context, m MetricPoint) error {
	_, err := db.Pool.Exec(ctx, "INSERT INTO metric_points (time, tenant_id, service_id, metric_name, value, labels) VALUES ($1,$2,$3,$4,$5,$6)", m.Time, m.TenantID, m.ServiceID, m.MetricName, m.Value, m.Labels)
	return err
}
func (db *DB) GetRecentLogs(ctx context.Context, tenantID, serviceID string, limit int) ([]LogEntry, error) {
	q := "SELECT l.time, s.tenant_id::text, l.service_id::text, l.level, l.message, l.metadata FROM log_entries l JOIN services s ON s.id=l.service_id WHERE s.tenant_id=$1"
	args := []interface{}{tenantID}
	if serviceID != "" {
		q += " AND l.service_id=$2"
		args = append(args, serviceID)
	}
	q += " ORDER BY l.time DESC"
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}
	rows, err := db.Pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []LogEntry{}
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.Time, &l.TenantID, &l.ServiceID, &l.Level, &l.Message, &l.Metadata); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (db *DB) GetRecentAnomalies(ctx context.Context, tenantID string, limit int) ([]Anomaly, error) {
	rows, err := db.Pool.Query(ctx, "SELECT a.id::text, s.tenant_id::text, a.service_id::text, a.detected_at, a.anomaly_type, a.description, a.root_cause_summary, a.likely_cause, a.suggested_fix FROM anomalies a JOIN services s ON s.id=a.service_id WHERE s.tenant_id=$1 ORDER BY a.detected_at DESC LIMIT $2", tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Anomaly{}
	for rows.Next() {
		var a Anomaly
		if err := rows.Scan(&a.ID, &a.TenantID, &a.ServiceID, &a.DetectedAt, &a.AnomalyType, &a.Description, &a.RootCauseSummary, &a.LikelyCause, &a.SuggestedFix); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (db *DB) GetRecentMetrics(ctx context.Context, tenantID, serviceID string, limit int) ([]MetricPoint, error) {
	q := "SELECT m.time, s.tenant_id::text, m.service_id::text, m.metric_name, m.value, m.labels FROM metric_points m JOIN services s ON s.id=m.service_id WHERE s.tenant_id=$1"
	args := []interface{}{tenantID}
	if serviceID != "" {
		q += " AND m.service_id=$2"
		args = append(args, serviceID)
	}
	q += " ORDER BY m.time DESC"
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}
	rows, err := db.Pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []MetricPoint{}
	for rows.Next() {
		var m MetricPoint
		if err := rows.Scan(&m.Time, &m.TenantID, &m.ServiceID, &m.MetricName, &m.Value, &m.Labels); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func (db *DB) GetServices(ctx context.Context, tenantID string) ([]Service, error) {
	rows, err := db.Pool.Query(ctx, "SELECT id::text, tenant_id::text, name, created_at FROM services WHERE tenant_id=$1 ORDER BY created_at DESC", tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Service{}
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.TenantID, &s.Name, &s.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}
