package handler

import (
	"errors"
	"log/slog"
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

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a new task with the provided details
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body createTaskRequest true "Task creation request"
// @Success 201
// @Failure 400
// @Failure 500
// @Router / [post]
func CreateTask(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var in createTaskRequest
		if err := ctx.ShouldBindJSON(&in); err != nil {
			ctx.JSON(400, gin.H{"error": err})
			return
		}
		traceID := ctx.GetString("trace_id")
		t, err := taskService.New(ctx.Request.Context(), task.NewTask{
			Title:       in.Title,
			Description: in.Description,
			Priority:    convertTaskPriorityTextToEnum(in.Priority),
			DueDate:     in.DueDate,
			AssigneeID:  in.AssigneeID,
			ParentID:    in.ParentID,
			Category:    in.Category,
			Tags:        in.Tags,
		}, ctx.GetInt("user_id"), traceID)
		if err != nil {
			slog.Error("api failed", "err", err, "trace", traceID)
			ctx.JSON(500, gin.H{"error": err})
			return
		}
		ctx.JSON(201, convertTaskDomainToResponse(t))
	}
}

// GetTask godoc
// @Summary Get task by ID
// @Description Returns a single task by its ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200
// @Failure 404
// @Failure 500
// @Router /{id} [get]
func GetTask(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		traceID := ctx.GetString("trace_id")
		t, err := taskService.Get(ctx.Request.Context(), id, traceID)
		if err != nil {
			slog.Error("api failed", "err", err, "trace", traceID)
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

// DeleteTask godoc
// @Summary Delete a task
// @Description Permanently deletes a task by ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /{id} [delete]
func DeleteTask(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		traceID := ctx.GetString("trace_id")
		if err = taskService.Delete(ctx.Request.Context(), id, ctx.GetInt("user_id"), traceID); err != nil {
			slog.Error("api failed", "err", err, "trace", traceID)
			if errors.Is(err, task.ErrNotFound) {
				ctx.JSON(404, gin.H{"error": err.Error()})
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

// PaginatedListTask godoc
// @Summary Get paginated list of tasks
// @Description Returns a paginated list of tasks with optional filters
// @Tags Tasks
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param page_size query int false "Number of items per page" default(10) minimum(1) maximum(100)
// @Param status query string false "Filter by status" Enums(Todo, InProgress, Done, Cancelled)
// @Param assignee_id query int false "Filter by assignee ID"
// @Success 200
// @Failure 500
// @Router / [get]
func PaginatedListTask(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		page, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
		pageSize, _ := strconv.ParseInt(ctx.Query("page_size"), 10, 64)
		assigneeID, _ := strconv.ParseInt(ctx.Query("assignee_id"), 10, 64)
		traceID := ctx.GetString("trace_id")
		tasks, total, err := taskService.List(ctx.Request.Context(), task.ListRequest{
			Page:       int(page),
			PageSize:   int(pageSize),
			Status:     convertTaskStatusTextToEnum(ctx.Query("status")),
			AssigneeID: int(assigneeID),
		}, traceID)
		if err != nil {
			slog.Error("api failed", "err", err, "trace", traceID)
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

type changeStatusRequest struct {
	Status string `json:"status"`
}

// ChangeStatus godoc
// @Summary Change task status
// @Description Updates only the status of a task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body changeStatusRequest true "Status update request"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{id}/status [patch]
func ChangeStatus(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		var in changeStatusRequest
		if err := ctx.ShouldBindJSON(&in); err != nil {
			ctx.JSON(400, gin.H{"error": err})
			return
		}
		traceID := ctx.GetString("trace_id")
		if err = taskService.ChangeStatus(ctx.Request.Context(), id, convertTaskStatusTextToEnum(in.Status), ctx.GetInt("user_id"), ctx.GetString("trace_id")); err != nil {
			slog.Error("api failed", "err", err, "trace", traceID)
			if errors.Is(err, task.ErrNotFound) {
				ctx.JSON(404, gin.H{"error": err.Error()})
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

type changeAssigneeRequest struct {
	AssigneeID int `json:"assignee_id"`
}

// ChangeAssignee godoc
// @Summary Change task assignee
// @Description Updates the assignee of a tChangeAssigneeask
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body changeAssigneeRequest true "Assignee update request"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{id}/assignee [patch]
func ChangeAssignee(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		var in changeAssigneeRequest
		if err := ctx.ShouldBindJSON(&in); err != nil {
			ctx.JSON(400, gin.H{"error": err})
			return
		}
		traceID := ctx.GetString("trace_id")
		if err = taskService.ChangeAssignee(ctx.Request.Context(), id, in.AssigneeID, ctx.GetInt("user_id"), ctx.GetString("trace_id")); err != nil {
			slog.Error("api failed", "err", err, "trace", traceID)
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

type changePriorityRequest struct {
	Priority string `json:"priority"`
}

// ChangePriority godoc
// @Summary Change task priority
// @Description Updates the priority of a task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body changePriorityRequest true "Priority update request"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{id}/priority [patch]
func ChangePriority(taskService task.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id param"})
			return
		}
		var in changePriorityRequest
		if err := ctx.ShouldBindJSON(&in); err != nil {
			ctx.JSON(400, gin.H{"error": err})
			return
		}
		traceID := ctx.GetString("trace_id")
		if err = taskService.ChangePriority(ctx.Request.Context(), id, convertTaskPriorityTextToEnum(in.Priority), ctx.GetInt("user_id"), traceID); err != nil {
			slog.Error("api failed", "err", err, "trace", traceID)
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
