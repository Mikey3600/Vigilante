package anomaly

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/user/vigilante/internal/ai"
	"github.com/user/vigilante/internal/storage"
)

// Engine runs the periodic anomaly checks.
type Engine struct {
	DB *storage.DB
	AI *ai.Client
}

// Start begins a ticker to scan for anomalies every 30 seconds.
func (e *Engine) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.checkMetrics(ctx)
		}
	}
}

func (e *Engine) checkMetrics(ctx context.Context) {
	// Query for services that had an average latency above 200ms in the last 1 minute
	rows, err := e.DB.Pool.Query(ctx, `
		SELECT service_id, AVG(value) as avg_latency
		FROM metric_points
		WHERE metric_name = 'latency' AND time > NOW() - INTERVAL '1 minute'
		GROUP BY service_id
		HAVING AVG(value) > 200
	`)
	if err != nil {
		log.Printf("Failed to query anomalies: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var serviceID string
		var avgLatency float64
		if err := rows.Scan(&serviceID, &avgLatency); err != nil {
			continue
		}

		// Prevent duplicate anomaly generation within a short timeframe
		var recentCount int
		e.DB.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM anomalies 
			WHERE service_id = $1 AND detected_at > NOW() - INTERVAL '5 minutes'
		`, serviceID).Scan(&recentCount)
		if recentCount > 0 {
			continue // Already reported recently
		}

		log.Printf("Anomaly detected for service %s! Fetching logs for AI analysis...", serviceID)

		// Fetch logs for AI
		logs, err := e.DB.GetRecentLogs(ctx, serviceID, 50)
		if err != nil {
			log.Printf("Failed to fetch logs: %v", err)
			continue
		}

		if e.AI == nil {
			log.Println("AI Engine not configured, skipping AI analysis")
			continue
		}

		// Proceed to AI analysis
		anomalyDesc := fmt.Sprintf("High average latency detected: %.2fms", avgLatency)
		report, err := e.AI.AnalyzeLogs(ctx, logs, anomalyDesc)
		if err != nil {
			log.Printf("AI analysis failed: %v", err)
			continue
		}

		err = e.DB.InsertAnomaly(ctx, storage.Anomaly{
			ServiceID:        serviceID,
			AnomalyType:      "Latency Spike",
			Description:      anomalyDesc,
			RootCauseSummary: report.Summary,
			LikelyCause:      report.LikelyCause,
			SuggestedFix:     report.SuggestedFix,
		})
		
		if err != nil {
			log.Printf("Failed to insert AI anomaly: %v", err)
		} else {
			log.Printf("AI Anomaly reported for service %s", serviceID)
		}
	}
}
