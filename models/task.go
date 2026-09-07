package models

import (
	"time"

	"github.com/lib/pq"
)

type Task struct {
	ID          int64 `gorm:"primarykey"`
	Title       string
	Description string
	Status      int
	Priority    int
	DueDate     time.Time
	Category    string
	ParentID    int64
	Tags        pq.StringArray `gorm:"type:text[]"`
	AssigneeID  int
	CreatorID   int
	CreatedAt   time.Time
}

func (Task) TableName() string {
	return "tasks"
}
