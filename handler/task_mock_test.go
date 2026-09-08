package handler_test

import (
	"context"

	"github.com/amirzayi/graph/task"
)

type mockTaskService struct {
	NewFunc            func(context.Context, task.NewTask, int) (task.Task, error)
	ChangeStatusFunc   func(context.Context, int64, task.Status, int) error
	ChangeAssigneeFunc func(context.Context, int64, int, int) error
	ChangePriorityFunc func(context.Context, int64, task.Priority, int) error
	GetFunc            func(context.Context, int64) (task.Task, error)
	DeleteFunc         func(context.Context, int64, int) error
	ListFunc           func(context.Context, task.ListRequest) ([]task.Task, int64, error)
}

func (m *mockTaskService) New(ctx context.Context, arg task.NewTask, creatorID int) (task.Task, error) {
	return m.NewFunc(ctx, arg, creatorID)
}
func (m *mockTaskService) ChangeStatus(ctx context.Context, id int64, newStatus task.Status, currentUserID int) error {
	return m.ChangeStatusFunc(ctx, id, newStatus, currentUserID)
}
func (m *mockTaskService) ChangeAssignee(ctx context.Context, id int64, newAssigneeID, currentUserID int) error {
	return m.ChangeAssigneeFunc(ctx, id, newAssigneeID, currentUserID)
}
func (m *mockTaskService) ChangePriority(ctx context.Context, id int64, newPriority task.Priority, currentUserID int) error {
	return m.ChangePriorityFunc(ctx, id, newPriority, currentUserID)
}
func (m *mockTaskService) Get(ctx context.Context, id int64) (task.Task, error) {
	return m.GetFunc(ctx, id)
}
func (m *mockTaskService) Delete(ctx context.Context, id int64, currentUserID int) error {
	return m.DeleteFunc(ctx, id, currentUserID)
}
func (m *mockTaskService) List(ctx context.Context, req task.ListRequest) ([]task.Task, int64, error) {
	return m.ListFunc(ctx, req)
}

var _ task.Service = (*mockTaskService)(nil)
