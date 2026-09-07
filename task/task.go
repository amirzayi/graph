package task

import (
	"context"
	"slices"
	"time"
)

type Status int

const (
	statusDefault Status = iota
	StatusTodo
	StatusInProgress
	StatusDone
	StatusCancelled
)

type Priority int

const (
	priorityDefault Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
)

var statusStateMachine = map[Status][]Status{
	StatusTodo:       {StatusInProgress, StatusCancelled},
	StatusInProgress: {StatusDone, StatusCancelled},
	StatusDone:       nil,
	StatusCancelled:  nil,
}

func (s Status) isValidNext(newStatus Status) bool {
	return slices.Contains(statusStateMachine[s], newStatus)
}

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      Status
	Priority    Priority
	DueDate     time.Time
	Category    string
	ParentID    int64
	Tags        []string
	AssigneeID  int
	CreatorID   int
	CreatedAt   time.Time
}

type NewTask struct {
	Title       string
	Description string
	Priority    Priority
	DueDate     time.Time
	AssigneeID  int
	CreatorID   int
	ParentID    int64
	Category    string
	Tags        []string
}

type Audit interface {
	Log(event string, creator int) error
}

type Repository interface {
	Get(context.Context, int64) (Task, error)
	Create(context.Context, Task) (int64, error)
	ChangeStatus(context.Context, int64, Status) error
	ChangeAssignee(context.Context, int64, int) error
	ChangePriority(context.Context, int64, Priority) error
	List(context.Context, ListRequest) ([]Task, int64, error)
	Delete(ctx context.Context, id int64) error
}
