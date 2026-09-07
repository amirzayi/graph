package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/amirzayi/graph/task"
	"github.com/gin-gonic/gin"
)

type taskResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date,omitzero"`
	Category    string    `json:"category,omitempty"`
	ParentID    int64     `json:"parent_id,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	AssigneeID  int       `json:"assignee_id,omitempty"`
	CreatorID   int       `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type createTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	AssigneeID  int       `json:"assignee_id,omitempty"`
	ParentID    int64     `json:"parent_id,omitempty"`
	Category    string    `json:"category,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}

func CreateTask(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var in createTaskRequest
		if err := ctx.ShouldBindJSON(&in); err != nil {
			ctx.JSON(400, gin.H{"error": err})
			return
		}

		t, err := taskService.New(ctx.Request.Context(), task.NewTask{
			Title:       in.Title,
			Description: in.Description,
			Priority:    convertTaskPriorityTextToEnum(in.Priority),
			DueDate:     in.DueDate,
			AssigneeID:  in.AssigneeID,
			CreatorID:   ctx.GetInt("user_id"),
			ParentID:    in.ParentID,
			Category:    in.Category,
			Tags:        in.Tags,
		}, ctx.GetInt("user_id"))
		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
			return
		}
		ctx.JSON(201, convertTaskDomainToResponse(t))
	}
}

func GetTask(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		t, err := taskService.Get(ctx.Request.Context(), id)
		if err != nil {
			if errors.Is(err, task.ErrNotFound) {
				ctx.JSON(404, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(500, http.StatusText(500))
			return
		}
		ctx.JSON(200, convertTaskDomainToResponse(t))
	}
}

func DeleteTask(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		if err = taskService.Delete(ctx.Request.Context(), id, ctx.GetInt("user_id")); err != nil {
			if errors.Is(err, task.ErrNotFound) {
				ctx.JSON(404, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(500, http.StatusText(500))
			return
		}
		ctx.Status(204)
	}
}

func PaginatedListTask(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		page, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
		pageSize, _ := strconv.ParseInt(ctx.Query("page_size"), 10, 64)
		assigneeID, _ := strconv.ParseInt(ctx.Query("assignee_id"), 10, 64)
		tasks, total, err := taskService.List(ctx.Request.Context(), task.ListRequest{
			Page:       int(page),
			PageSize:   int(pageSize),
			Status:     convertTaskStatusTextToEnum(ctx.Query("status")),
			AssigneeID: int(assigneeID),
		})
		if err != nil {
			ctx.JSON(500, http.StatusText(500))
			return
		}
		taskReponse := make([]taskResponse, 0, len(tasks))
		for _, t := range tasks {
			taskReponse = append(taskReponse, convertTaskDomainToResponse(t))
		}
		ctx.JSON(200, gin.H{
			"data":  taskReponse,
			"total": total,
		})
	}
}

func ChangeStatus(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		var in struct {
			Status string `json:"status"`
		}
		if err := ctx.ShouldBindJSON(&in); err != nil {
			ctx.JSON(400, gin.H{"error": err})
			return
		}
		if err = taskService.ChangeStatus(ctx.Request.Context(), id, convertTaskStatusTextToEnum(in.Status), ctx.GetInt("user_id")); err != nil {
			if errors.Is(err, task.ErrNotFound) {
				ctx.JSON(404, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, task.ErrChangeStatusNotValid) {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, task.ErrChangeStatusNotValid) {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, task.ErrTaskNotBelongs) {
				ctx.JSON(403, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(500, http.StatusText(500))
			return
		}
		ctx.Status(204)
	}
}

func ChangeAssignee(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		var in struct {
			AssigneeID int `json:"assignee_id"`
		}
		if err := ctx.ShouldBindJSON(&in); err != nil {
			ctx.JSON(400, gin.H{"error": err})
			return
		}
		if err = taskService.ChangeAssignee(ctx.Request.Context(), id, in.AssigneeID, ctx.GetInt("user_id")); err != nil {
			if errors.Is(err, task.ErrNotFound) {
				ctx.JSON(404, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(500, http.StatusText(500))
			return
		}
		ctx.Status(204)
	}
}

func ChangePriority(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		var in struct {
			Priority string `json:"priority"`
		}
		if err := ctx.ShouldBindJSON(&in); err != nil {
			ctx.JSON(400, gin.H{"error": err})
			return
		}
		if err = taskService.ChangePriority(ctx.Request.Context(), id, convertTaskPriorityTextToEnum(in.Priority), ctx.GetInt("user_id")); err != nil {
			if errors.Is(err, task.ErrNotFound) {
				ctx.JSON(404, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, task.ErrInvalidPriority) {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(500, http.StatusText(500))
			return
		}
		ctx.Status(204)
	}
}
