package router

import (
	"strconv"
	"time"

	"github.com/amirzayi/graph/handler"
	"github.com/amirzayi/graph/metric"
	"github.com/amirzayi/graph/task"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Bootstrap(r gin.IRouter, taskService task.Service) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Use(PrometheusMiddleware())

	gp := r.Group("/api/v1/tasks")
	gp.POST("", handler.CreateTask(taskService))
	gp.GET("", handler.PaginatedListTask(taskService))
	gp.GET("/:id", handler.GetTask(taskService))
	gp.DELETE("/:id", handler.DeleteTask(taskService))
	gp.PATCH("/:id/status", handler.ChangeStatus(taskService))
	gp.PATCH("/:id/assignee", handler.ChangeAssignee(taskService))
	gp.PATCH("/:id/priority", handler.ChangePriority(taskService))
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()

		metric.RequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		metric.RequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)
	}
}
