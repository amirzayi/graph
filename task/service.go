package task

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

const (
	EventTaskCreated         = "task.created"
	EventTaskStatusChanged   = "task.status_changed"
	EventTaskAssigneeChanged = "task.status_changed"
	EventTaskPriorityChanged = "task.priority_changed"
	EventTaskDeleted         = "task.deleted"
)

var (
	ErrTaskNotBelongs       = errors.New("task not belog to you")
	ErrChangeStatusNotValid = errors.New("requested status couldn't apply to present task status")
	ErrDeleteTask           = errors.New("cannot delete done or in progress task")
	ErrInvalidPriority      = errors.New("invalid priority")
)

type Service interface {
	New(context.Context, NewTask, int, string) (Task, error)
	ChangeStatus(context.Context, int64, Status, int, string) error
	ChangeAssignee(context.Context, int64, int, int, string) error
	ChangePriority(context.Context, int64, Priority, int, string) error
	Get(context.Context, int64, string) (Task, error)
	Delete(context.Context, int64, int, string) error
	List(context.Context, ListRequest, string) ([]Task, int64, error)
}

type service struct {
	repo  Repository
	audit Audit
}

func NewService(repo Repository, audit Audit) *service {
	return &service{repo: repo, audit: audit}
}

func (s *service) New(ctx context.Context, arg NewTask, creatorID int, traceID string) (Task, error) {
	t := Task{
		Title:       arg.Title,
		Description: arg.Description,
		Priority:    arg.Priority,
		Status:      StatusTodo,
		AssigneeID:  arg.AssigneeID,
		CreatorID:   creatorID,
		CreatedAt:   time.Now(),
		ParentID:    arg.ParentID,
		DueDate:     arg.DueDate,
		Category:    arg.Category,
		Tags:        arg.Tags,
	}
	id, err := s.repo.Create(ctx, t)
	if err != nil {
		return Task{}, err
	}
	t.ID = id
	if err = s.audit.Log(EventTaskCreated, creatorID, traceID); err != nil {
		slog.Error("create task: audit log failed", "error", err)
	}
	return t, nil
}

func (s *service) Get(ctx context.Context, id int64, traceID string) (Task, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) ChangeStatus(ctx context.Context, id int64, newStatus Status, currentUserID int, traceID string) error {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	// due to keep idempotent
	if t.Status == newStatus {
		return nil
	}
	if !t.Status.isValidNext(newStatus) {
		return ErrChangeStatusNotValid
	}
	if t.AssigneeID != currentUserID {
		return ErrTaskNotBelongs
	}
	if err = s.audit.Log(EventTaskStatusChanged, currentUserID, traceID); err != nil {
		slog.Error("change status: audit log failed", "error", err)
	}
	return nil
}

func (s *service) ChangeAssignee(ctx context.Context, id int64, newAssigneeID, currentUserID int, traceID string) error {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	// keep idempotent
	if t.AssigneeID == newAssigneeID {
		return nil
	}
	if err = s.repo.ChangeAssignee(ctx, id, newAssigneeID); err != nil {
		return err
	}
	if err = s.audit.Log(EventTaskAssigneeChanged, currentUserID, traceID); err != nil {
		slog.Error("change assignee: audit log failed", "error", err)
	}
	return nil
}

func (s *service) ChangePriority(ctx context.Context, id int64, newPriority Priority, currentUserID int, traceID string) error {
	if newPriority == priorityDefault {
		return ErrInvalidPriority
	}
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	// due to keep idempotent
	if t.Priority == newPriority {
		println("equal")
		return nil
	}
	println("changing")
	if err = s.repo.ChangePriority(ctx, id, newPriority); err != nil {
		return err
	}
	println("done")

	if err = s.audit.Log(EventTaskPriorityChanged, currentUserID, traceID); err != nil {
		slog.Error("change priority: audit log failed", "error", err)
	}
	return nil
}

type ListRequest struct {
	Page       int
	PageSize   int
	Status     Status
	AssigneeID int
}
type ListResponse struct {
	Tasks []Task
	Total int64
}

func (s *service) List(ctx context.Context, req ListRequest, traceID string) ([]Task, int64, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	return s.repo.List(ctx, req)
}

func (s *service) Delete(ctx context.Context, id int64, currentUserID int, traceID string) error {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		// keep idempotent: already deleted
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		return err
	}
	if t.Status == StatusDone || t.Status == StatusInProgress {
		return ErrDeleteTask
	}
	if t.CreatorID != currentUserID {
		return ErrTaskNotBelongs
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return err
	}
	if err = s.audit.Log(EventTaskDeleted, currentUserID, traceID); err != nil {
		slog.Error("delete task: audit log failed", "error", err)
	}
	return nil
}
