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
type Anomaly struct {
	ID, TenantID, ServiceID, AnomalyType, Description, RootCauseSummary, LikelyCause, SuggestedFix, Severity string
	DetectedAt                                                                                               time.Time
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
	rows, err := db.Pool.Query(ctx, "SELECT time, tenant_id, service_id, level, message, metadata FROM log_entries WHERE tenant_id=$1 AND service_id=$2 ORDER BY time DESC LIMIT $3", tenantID, serviceID, limit)
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
