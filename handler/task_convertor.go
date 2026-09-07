package handler

import (
	"strings"

	"github.com/amirzayi/graph/task"
)

func convertTaskPriorityTextToEnum(p string) task.Priority {
	switch strings.TrimSpace(strings.ToLower(p)) {
	case "high":
		return task.PriorityHigh
	case "medium":
		return task.PriorityMedium
	case "low":
		return task.PriorityLow
	default:
		return task.Priority(0)
	}
}

func convertTaskPriorityToText(p task.Priority) string {
	switch p {
	case task.PriorityHigh:
		return "High"
	case task.PriorityMedium:
		return "Medium"
	case task.PriorityLow:
		return "Low"
	default:
		return ""
	}
}

func convertTaskStatusTextToEnum(s string) task.Status {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "cancelled":
		return task.StatusCancelled
	case "done":
		return task.StatusDone
	case "inprogress":
		return task.StatusInProgress
	case "todo":
		return task.StatusTodo
	default:
		return task.Status(0)
	}
}

func convertTaskStatusToText(s task.Status) string {
	switch s {
	case task.StatusCancelled:
		return "Cancelled"
	case task.StatusDone:
		return "Done"
	case task.StatusInProgress:
		return "InProgress"
	case task.StatusTodo:
		return "Todo"
	default:
		return ""
	}
}

func convertTaskDomainToResponse(t task.Task) taskResponse {
	return taskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      convertTaskStatusToText(t.Status),
		Priority:    convertTaskPriorityToText(t.Priority),
		DueDate:     t.DueDate,
		Category:    t.Category,
		ParentID:    t.ParentID,
		Tags:        t.Tags,
		AssigneeID:  t.AssigneeID,
		CreatorID:   t.CreatorID,
		CreatedAt:   t.CreatedAt,
	}
}
