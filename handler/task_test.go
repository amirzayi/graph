package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amirzayi/graph/handler"
	"github.com/amirzayi/graph/task"
	"github.com/amirzayi/graph/testhelper"
	"github.com/gin-gonic/gin"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	h := handler.CreateTask(&mockTaskService{
		NewFunc: func(ctx context.Context, nt task.NewTask, i int) (task.Task, error) {
			if i < 1 {
				return task.Task{}, errors.New("bad user_id")
			}
			if nt.Title == "" {
				return task.Task{}, errors.New("bad title")
			}
			if nt.Priority == 0 {
				return task.Task{}, errors.New("bad priority")
			}
			return task.Task{}, nil
		},
	})

	t.Run("empty request body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("valid request body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		data := map[string]any{
			"title":       "Deploy security patch",
			"description": "Apply security patch to production",
			"priority":    "medium",
			"due_date":    "2026-09-08T09:00:00Z",
			"assignee_id": 456,
			"category":    "security",
			"tags":        []string{"security", "deployment"},
		}
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(data)
		testhelper.MustNilError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", buf)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("valid request body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Set("user_id", 13)
		data := map[string]any{
			"title":       "Deploy security patch",
			"description": "Apply security patch to production",
			"priority":    "medium",
			"due_date":    "2026-09-08T09:00:00Z",
			"assignee_id": 456,
			"category":    "security",
			"tags":        []string{"security", "deployment"},
		}
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(data)
		testhelper.MustNilError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", buf)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusCreated)
	})
}

func TestGetTask(t *testing.T) {
	t.Parallel()
	h := handler.GetTask(&mockTaskService{
		GetFunc: func(ctx context.Context, i int64) (task.Task, error) {
			if i == 1000 {
				return task.Task{}, task.ErrNotFound
			}
			if i <= 0 {
				return task.Task{}, errors.New("bad id")
			}
			return task.Task{}, nil
		},
	})

	t.Run("invalid id param", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("zero id", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "0")
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "1000")
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusNotFound)
	})

	t.Run("found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "1")
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusOK)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Parallel()
	h := handler.DeleteTask(&mockTaskService{
		DeleteFunc: func(ctx context.Context, i1 int64, i2 int) error {
			if i1 == 755 {
				return task.ErrNotFound
			}
			if i2 == 755 {
				return task.ErrTaskNotBelongs
			}
			if i1%2 == 0 {
				return errors.New("some internal error")
			}
			return nil
		},
	})

	t.Run("invalid id param", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("zero id", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "0")
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "755")
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusNotFound)
	})

	t.Run("not belong", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "504")
		c.Set("user_id", 755)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusForbidden)
	})

	t.Run("internal error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "500")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "1")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusOK)
	})
}

func TestChangeStatus(t *testing.T) {
	t.Parallel()
	h := handler.ChangeStatus(&mockTaskService{
		ChangeStatusFunc: func(ctx context.Context, i1 int64, i2 task.Status, i3 int) error {
			if i1 == 755 {
				return task.ErrNotFound
			}
			if i3 == 755 {
				return task.ErrTaskNotBelongs
			}
			if i2 == 0 {
				return task.ErrChangeStatusNotValid
			}
			if i1%2 == 0 {
				return errors.New("some internal error")
			}
			return nil
		},
	})

	t.Run("invalid id param", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "abc")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("invalid json body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(`{invalid}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "123")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("invalid status", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(`{"status":"invalid_status"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "123")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(`{"status":"done"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "755")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusNotFound)
	})

	t.Run("not belong", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(`{"status":"done"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "123")
		c.Set("user_id", 755)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusForbidden)
	})

	t.Run("internal error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(`{"status":"done"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "500")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(`{"status":"done"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "1")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusOK)
	})
}

func TestChangeAssignee(t *testing.T) {
	t.Parallel()
	h := handler.ChangeAssignee(&mockTaskService{
		ChangeAssigneeFunc: func(ctx context.Context, i1 int64, i2 int, i3 int) error {
			if i1 == 755 {
				return task.ErrNotFound
			}
			if i1%2 == 0 {
				return errors.New("some internal error")
			}
			return nil
		},
	})

	t.Run("invalid id param", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "abc")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("invalid json body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(`{invalid}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "123")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		body := `{"assignee_id": 500}`
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "755")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusNotFound)
	})

	t.Run("internal error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		body := `{"assignee_id": 500}`
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "500")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		body := `{"assignee_id": 500}`
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/tasks", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "1")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusOK)
	})
}

func TestChangePriority(t *testing.T) {
	t.Parallel()
	h := handler.ChangePriority(&mockTaskService{
		ChangePriorityFunc: func(ctx context.Context, i1 int64, i2 task.Priority, i3 int) error {
			if i1 == 755 {
				return task.ErrNotFound
			}
			if i2 == 999 {
				return task.ErrInvalidPriority
			}
			if i1%2 == 0 {
				return errors.New("some internal error")
			}
			return nil
		},
	})

	t.Run("invalid id param", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/tasks", http.NoBody)
		c.AddParam("id", "abc")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("invalid json body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/tasks", bytes.NewBufferString(`{invalid}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "123")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		body := `{"priority":"high"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/tasks", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "755")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusNotFound)
	})

	t.Run("internal error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		body := `{"priority":"high"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/tasks", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "500")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		body := `{"priority":"high"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/tasks", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "1")
		c.Set("user_id", 1000)
		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusOK)
	})
}

func TestPaginatedListTask(t *testing.T) {
	t.Parallel()

	h := handler.PaginatedListTask(&mockTaskService{
		ListFunc: func(ctx context.Context, req task.ListRequest) ([]task.Task, int64, error) {
			mockTasks := []task.Task{
				{
					ID:       1,
					Title:    "Task 1",
					Status:   task.StatusInProgress,
					Priority: task.PriorityHigh,
				},
				{
					ID:       2,
					Title:    "Task 2",
					Status:   task.StatusDone,
					Priority: task.PriorityLow,
				},
				{
					ID:       3,
					Title:    "Task 3",
					Status:   task.StatusCancelled,
					Priority: task.PriorityMedium,
				},
				{
					ID:       4,
					Title:    "Task 4",
					Status:   task.StatusTodo,
					Priority: task.PriorityLow,
				},
			}
			if req.PageSize > 100 {
				return nil, 0, errors.New("some internal error")
			}
			if req.Page == 0 && req.PageSize == 0 {
				return nil, 0, nil
			}
			if req.AssigneeID == 755 {
				return mockTasks[:2], 2, nil
			}

			return mockTasks, 4, nil
		},
	})

	t.Run("empty data", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)

		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusOK)

		var response map[string]any
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testhelper.MustEqual(t, err, nil)
		testhelper.MapMustContainKeys(t, response, "total", "data")
		testhelper.MustEqual(t, response["total"].(float64), float64(0))
	})

	t.Run("internal error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		req.URL.RawQuery = "page_size=101"
		c.Request = req

		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("contains data", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", http.NoBody)
		req.URL.RawQuery = "page=1&page_size=50"
		c.Request = req

		h(c)
		testhelper.MustEqual(t, rec.Code, http.StatusOK)

		var response map[string]any
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testhelper.MustEqual(t, err, nil)
		testhelper.MapMustContainKeys(t, response, "total", "data")
		testhelper.MustEqual(t, response["total"].(float64), float64(4))

		data, ok := response["data"].([]any)
		if !ok {
			t.Fatalf("expected array but got, %T", response["data"])
		}
		if len(data) != 4 {
			t.Fatalf("expected 4 items, but got %d", len(data))
		}
	})
}
