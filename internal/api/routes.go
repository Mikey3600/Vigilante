package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/vigilante/internal/ai"
	"github.com/user/vigilante/internal/auth"
	"github.com/user/vigilante/internal/ingestion"
	"github.com/user/vigilante/internal/storage"
)

// SetupRouter returns the full Gin router.
func SetupRouter(db *storage.DB, aiClient *ai.Client) *gin.Engine {
	r := gin.Default()

	r.StaticFile("/", "./index.html")
	r.StaticFile("/dashboard", "./index.html")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.Use(auth.JWTMiddleware())
	{
		api.POST("/logs", func(c *gin.Context) {
			var payload ingestion.BatchLogPayload
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := ingestion.ProcessLogs(c.Request.Context(), db, payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process logs"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "accepted"})
		})

		api.POST("/metrics", func(c *gin.Context) {
			var payload ingestion.BatchMetricPayload
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := ingestion.ProcessMetrics(c.Request.Context(), db, payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process metrics"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "accepted"})
		})

		api.GET("/dashboard", func(c *gin.Context) {
			tenantID := auth.GetTenantID(c)
			
			logs, _ := db.GetRecentLogsForTenant(c.Request.Context(), tenantID, 50)
			if logs == nil {
				logs = []storage.LogEntry{}
			}
			metrics, _ := db.GetRecentMetricsForTenant(c.Request.Context(), tenantID, 50)
			if metrics == nil {
				metrics = []storage.MetricPoint{}
			}
			anomalies, _ := db.GetRecentAnomaliesForTenant(c.Request.Context(), tenantID, 10)
			if anomalies == nil {
				anomalies = []storage.Anomaly{}
			}

			c.JSON(http.StatusOK, gin.H{
				"logs": logs,
				"metrics": metrics,
				"anomalies": anomalies,
			})
		})

		api.POST("/analyze", func(c *gin.Context) {
			tenantID := auth.GetTenantID(c)

			type AnalyzeRequest struct {
				ServiceID string `json:"service_id"`
				AnomalyType string `json:"anomaly_type"`
				Description string `json:"description"`
			}
			
			var req AnalyzeRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			logs, err := db.GetRecentLogs(c.Request.Context(), req.ServiceID, 50)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch context logs"})
				return
			}
			
			if aiClient == nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "AI processor unavailable"})
				return
			}

			report, err := aiClient.AnalyzeLogs(c.Request.Context(), logs, req.AnomalyType + ": " + req.Description)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "AI analysis failed"})
				return
			}

			c.JSON(http.StatusOK, report)
		})
	}

	return r
}
