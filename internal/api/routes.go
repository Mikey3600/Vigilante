package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/user/vigilante/internal/auth"
	"github.com/user/vigilante/internal/ingestion"
	"github.com/user/vigilante/internal/storage"
)

func SetupRouter(db *storage.DB) *gin.Engine {
	start := time.Now(); r:=gin.New(); r.Use(gin.Recovery(), gin.Logger(), requestID())
	r.GET("/health", func(c *gin.Context){ dbs:="ok"; if err:=db.Ping(c.Request.Context()); err!=nil{dbs="down"}; c.JSON(http.StatusOK, gin.H{"status":"ok","db":dbs,"version":"1.0.0","uptime":time.Since(start).String()}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	v1:=r.Group("/api/v1"); v1.Use(auth.JWTMiddleware())
	v1.POST("/logs", func(c *gin.Context){ var p ingestion.BatchLogPayload; if err:=c.ShouldBindJSON(&p); err!=nil{jsonErr(c,400,"INVALID_REQUEST",err.Error()); return}; if err:=ingestion.ProcessLogs(c,db,auth.GetTenantID(c),p); err!=nil{jsonErr(c,400,"INVALID_LOGS",err.Error()); return}; c.JSON(200,gin.H{"status":"accepted"}) })
	v1.POST("/metrics", func(c *gin.Context){ var p ingestion.BatchMetricPayload; if err:=c.ShouldBindJSON(&p); err!=nil{jsonErr(c,400,"INVALID_REQUEST",err.Error()); return}; if err:=ingestion.ProcessMetrics(c,db,auth.GetTenantID(c),p); err!=nil{jsonErr(c,400,"INVALID_METRICS",err.Error()); return}; c.JSON(200,gin.H{"status":"accepted"}) })
	return r
}
func requestID() gin.HandlerFunc { return func(c *gin.Context){id:=uuid.NewString(); c.Set("request_id",id); c.Writer.Header().Set("X-Request-ID",id); c.Next()} }
func jsonErr(c *gin.Context, status int, code,msg string){ rid,_:=c.Get("request_id"); c.JSON(status, gin.H{"error":msg,"code":code,"request_id":rid}) }
