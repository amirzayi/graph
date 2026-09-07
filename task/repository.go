package task

import (
	"context"

	"github.com/amirzayi/graph/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type sqlRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) sqlRepository {
	return sqlRepository{db: db}
}

func (r sqlRepository) Get(ctx context.Context, id int64) (Task, error) {
	t, err := gorm.G[models.Task](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return Task{}, err
	}
	return Task{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      Status(t.Status),
		Priority:    Priority(t.Priority),
		DueDate:     t.DueDate,
		Category:    t.Category,
		ParentID:    t.ParentID,
		Tags:        t.Tags,
		AssigneeID:  t.AssigneeID,
		CreatorID:   t.CreatorID,
		CreatedAt:   t.CreatedAt,
	}, nil
}

func (r sqlRepository) Create(ctx context.Context, t Task) (int64, error) {
	ts := models.Task{
		Title:       t.Title,
		Description: t.Description,
		Status:      int(t.Status),
		Priority:    int(t.Priority),
		DueDate:     t.DueDate,
		Category:    t.Category,
		ParentID:    t.ParentID,
		Tags:        pq.StringArray(t.Tags),
		AssigneeID:  t.AssigneeID,
		CreatorID:   t.CreatorID,
	}
	err := gorm.G[models.Task](r.db).Create(ctx, &ts)
	if err != nil {
		return 0, err
	}
	return t.ID, nil
}

func (r sqlRepository) ChangeStatus(ctx context.Context, id int64, newStatus Status) error {
	_, err := gorm.G[models.Task](r.db).Where("id = ?", id).Updates(ctx, models.Task{Status: int(newStatus)})
	return err
}

func (r sqlRepository) ChangeAssignee(ctx context.Context, id int64, assigneeID int) error {
	_, err := gorm.G[models.Task](r.db).Where("id = ?", id).Updates(ctx, models.Task{AssigneeID: assigneeID})
	return err
}

func (r sqlRepository) ChangePriority(ctx context.Context, id int64, newPriority Priority) error {
	_, err := gorm.G[models.Task](r.db).Where("id = ?", id).Updates(ctx, models.Task{Priority: int(newPriority)})
	return err
}

func (r sqlRepository) List(ctx context.Context, req ListRequest) (ListResponse, error) {
	q := r.db.Model(&models.Task{}).WithContext(ctx)

	if req.AssigneeID > 0 {
		q = q.Where("assignee_id = ?", req.AssigneeID)
	}
	if req.Status != statusDefault {
		q = q.Where("status = ?", req.Status)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return ListResponse{}, err
	}
	if count == 0 {
		return ListResponse{}, nil
	}

	var tasks []models.Task
	if err := q.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Find(&tasks).Error; err != nil {
		return ListResponse{}, err
	}

	response := ListResponse{Total: count, Tasks: make([]Task, 0, len(tasks))}
	for _, t := range tasks {
		response.Tasks = append(response.Tasks, Task{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      Status(t.Status),
			Priority:    Priority(t.Priority),
			DueDate:     t.DueDate,
			Category:    t.Category,
			ParentID:    t.ParentID,
			Tags:        t.Tags,
			AssigneeID:  t.AssigneeID,
			CreatorID:   t.CreatorID,
			CreatedAt:   t.CreatedAt,
		})
	}
	return response, nil
}
