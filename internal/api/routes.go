package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/vigilante/internal/auth"
	"github.com/user/vigilante/internal/ingestion"
	"github.com/user/vigilante/internal/storage"
)

// SetupRouter returns the full Gin router.
func SetupRouter(db *storage.DB) *gin.Engine {
	r := gin.Default()

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
	}

	return r
}
