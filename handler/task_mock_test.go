package handler_test

import (
	"context"

	"github.com/amirzayi/graph/task"
)

type mockTaskService struct {
	NewFunc            func(context.Context, task.NewTask, int, string) (task.Task, error)
	ChangeStatusFunc   func(context.Context, int64, task.Status, int, string) error
	ChangeAssigneeFunc func(context.Context, int64, int, int, string) error
	ChangePriorityFunc func(context.Context, int64, task.Priority, int, string) error
	GetFunc            func(context.Context, int64, string) (task.Task, error)
	DeleteFunc         func(context.Context, int64, int, string) error
	ListFunc           func(context.Context, task.ListRequest, string) ([]task.Task, int64, error)
}

func (m *mockTaskService) New(ctx context.Context, arg task.NewTask, creatorID int, traceID string) (task.Task, error) {
	return m.NewFunc(ctx, arg, creatorID, traceID)
}
func (m *mockTaskService) ChangeStatus(ctx context.Context, id int64, newStatus task.Status, currentUserID int, traceID string) error {
	return m.ChangeStatusFunc(ctx, id, newStatus, currentUserID, traceID)
}
func (m *mockTaskService) ChangeAssignee(ctx context.Context, id int64, newAssigneeID, currentUserID int, traceID string) error {
	return m.ChangeAssigneeFunc(ctx, id, newAssigneeID, currentUserID, traceID)
}
func (m *mockTaskService) ChangePriority(ctx context.Context, id int64, newPriority task.Priority, currentUserID int, traceID string) error {
	return m.ChangePriorityFunc(ctx, id, newPriority, currentUserID, traceID)
}
func (m *mockTaskService) Get(ctx context.Context, id int64, traceID string) (task.Task, error) {
	return m.GetFunc(ctx, id, traceID)
}
func (m *mockTaskService) Delete(ctx context.Context, id int64, currentUserID int, traceID string) error {
	return m.DeleteFunc(ctx, id, currentUserID, traceID)
}
func (m *mockTaskService) List(ctx context.Context, req task.ListRequest, traceID string) ([]task.Task, int64, error) {
	return m.ListFunc(ctx, req, traceID)
}

var _ task.Service = (*mockTaskService)(nil)
