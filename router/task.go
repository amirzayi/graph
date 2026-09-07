package router

import (
	"github.com/amirzayi/graph/handler"
	"github.com/amirzayi/graph/task"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	TaskService task.Service
}

func Bootstrap(r gin.IRouter, deps Dependencies) {
	gp := r.Group("/api/v1/tasks")

	gp.POST("", handler.CreateTask(deps.TaskService))
	gp.GET("", handler.PaginatedListTask(deps.TaskService))
	gp.GET("/:id", handler.GetTask(deps.TaskService))
	gp.DELETE("/:id", handler.DeleteTask(deps.TaskService))
	gp.PATCH("/:id/status", handler.ChangeStatus(deps.TaskService))
	gp.PATCH("/:id/assignee", handler.ChangeAssignee(deps.TaskService))
	gp.PATCH("/:id/priority", handler.ChangePriority(deps.TaskService))

}
