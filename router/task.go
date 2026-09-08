package router

import (
	"github.com/amirzayi/graph/handler"
	"github.com/amirzayi/graph/task"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Bootstrap(r gin.IRouter, taskService task.Service) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	gp := r.Group("/api/v1/tasks")

	gp.POST("", handler.CreateTask(taskService))
	gp.GET("", handler.PaginatedListTask(taskService))
	gp.GET("/:id", handler.GetTask(taskService))
	gp.DELETE("/:id", handler.DeleteTask(taskService))
	gp.PATCH("/:id/status", handler.ChangeStatus(taskService))
	gp.PATCH("/:id/assignee", handler.ChangeAssignee(taskService))
	gp.PATCH("/:id/priority", handler.ChangePriority(taskService))

}
