package task

import (
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
	ID          int
	Title       string
	Description string
	Status      Status
	Priority    Priority
	DueDate     time.Time
	Category    string
	ParentID    int
	Tags        []string
	AssigneeID  int
	CreatorID   int
	CreatedAt   time.Time
}
