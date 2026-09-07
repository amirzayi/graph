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
)

var (
	ErrTaskNotBelongs       = errors.New("task not belog to you")
	ErrChangeStatusNotValid = errors.New("requested status couldn't apply to present task status")
	ErrAlreadyAssigned      = errors.New("task already assigned")
)

type Service interface {
	New(context.Context, NewTask) (Task, error)
	ChangeStatus(context.Context, int, Status) error
	ChangeAssignee(context.Context, int, int) error
	ChangePriority(context.Context, int, Priority) error
}

type service struct {
	repo  Repository
	audit Audit
}

func NewService(repo Repository, audit Audit) service {
	return service{repo: repo, audit: audit}
}

func (s *service) New(ctx context.Context, creatorID int, arg NewTask) (Task, error) {
	t := Task{
		Title:       arg.Title,
		Description: arg.Description,
		Priority:    arg.Priority,
		Status:      StatusTodo,
		AssigneeID:  arg.AssigneeID,
		CreatorID:   creatorID,
		CreatedAt:   time.Now(),
	}
	id, err := s.repo.Create(ctx, t)
	if err != nil {
		return Task{}, err
	}
	t.ID = id
	if err = s.audit.Log(EventTaskCreated, arg.CreatorID); err != nil {
		slog.Error("create task: audit log failed", "error", err)
	}
	return t, nil
}

func (s *service) ChangeStatus(ctx context.Context, id, currentUserID int, newStatus Status) error {
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
	if err = s.audit.Log(EventTaskStatusChanged, currentUserID); err != nil {
		slog.Error("change status: audit log failed", "error", err)
	}
	return nil
}

func (s *service) ChangeAssignee(ctx context.Context, id, newAssigneeID int) error {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if t.AssigneeID == newAssigneeID {
		return ErrAlreadyAssigned
	}
	if err = s.repo.ChangeAssignee(ctx, id, newAssigneeID); err != nil {
		return err
	}
	if err = s.audit.Log(EventTaskAssigneeChanged, id); err != nil {
		slog.Error("change assignee: audit log failed", "error", err)
	}
	return nil
}

func (s *service) ChangePriority(ctx context.Context, id int, newPriority Priority) error {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	// due to keep idempotent
	if t.Priority == newPriority {
		return nil
	}
	if err = s.repo.ChangePriority(ctx, id, newPriority); err != nil {
		return err
	}
	if err = s.audit.Log(EventTaskPriorityChanged, id); err != nil {
		slog.Error("change priority: audit log failed", "error", err)
	}
	return nil
}
