package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amirzayi/graph/audit"
	"github.com/amirzayi/graph/models"
	"github.com/amirzayi/graph/router"
	"github.com/amirzayi/graph/task"
	"github.com/amirzayi/graph/testhelper"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var httpSrvAddress string

func TestMain(m *testing.M) {
	ginrouter := gin.Default()
	srv := httptest.NewServer(ginrouter)
	defer srv.Close()
	httpSrvAddress = srv.URL

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}
	_ = db.AutoMigrate(&models.Task{})

	router.Bootstrap(ginrouter, task.NewService(task.NewSQLRepository(db), audit.NewDiscarded()))
	m.Run()
}

func TestListTask(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		resp, err := http.Get(httpSrvAddress + "/api/v1/tasks")
		testhelper.MustNilError(t, err)
		testhelper.MustEqual(t, resp.StatusCode, http.StatusOK)
		defer func() {
			_ = resp.Body.Close()
		}()

		var got map[string]any
		err = json.NewDecoder(resp.Body).Decode(&got)
		testhelper.MustNilError(t, err)
		testhelper.MapMustContainKeys(t, got, "total", "data")

		if v, ok := got["total"].(float64); true {
			if !ok {
				t.Fatalf("expected numeric total, but is %T", got["total"])
			}
			if v != 0 {
				t.Fatalf("expected zero value total, but is %f", v)
			}
		}
	})

	t.Run("contains data", func(t *testing.T) {
		for range 5 {
			createTask(t)
		}
		resp, err := http.Get(httpSrvAddress + "/api/v1/tasks")
		testhelper.MustNilError(t, err)
		testhelper.MustEqual(t, resp.StatusCode, http.StatusOK)
		defer func() {
			_ = resp.Body.Close()
		}()

		var got map[string]any
		err = json.NewDecoder(resp.Body).Decode(&got)
		testhelper.MustNilError(t, err)
		testhelper.MapMustContainKeys(t, got, "total", "data")

		if v, ok := got["data"].([]any); true {
			if !ok {
				t.Fatalf("expected numeric total, but is %T", got["total"])
			}
			if len(v) != 5 {
				t.Fatalf("expected 5 data items, but got %d", len(v))
			}
		}
		if v, ok := got["total"].(float64); true {
			if !ok {
				t.Fatalf("expected numeric total, but is %T", got["total"])
			}
			if v != 5 {
				t.Fatalf("expected 5 total items, but got %f", v)
			}
		}
	})
}

func TestCreateTask(t *testing.T) {
	createTask(t)
}

func TestGetTask(t *testing.T) {
	t.Run("get created task", func(t *testing.T) {
		id := createTask(t)
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/tasks/%d", httpSrvAddress, id))
		testhelper.MustNilError(t, err)
		testhelper.MustEqual(t, resp.StatusCode, http.StatusOK)
		defer func() {
			_ = resp.Body.Close()
		}()

		var got map[string]any
		err = json.NewDecoder(resp.Body).Decode(&got)
		testhelper.MustNilError(t, err)
		testhelper.MapMustContainKeys(t, got, "id", "title", "description", "status", "priority", "due_date", "category", "tags", "assignee_id", "creator_id", "created_at")

	})

	t.Run("not existed task", func(t *testing.T) {
		resp, err := http.Get(httpSrvAddress + "/api/v1/tasks/10000")
		testhelper.MustNilError(t, err)
		testhelper.MustEqual(t, resp.StatusCode, http.StatusNotFound)
		defer resp.Body.Close()
	})
}

func TestDeleteTask(t *testing.T) {
	id := createTask(t)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/tasks/%d", httpSrvAddress, id), http.NoBody)
	testhelper.MustNilError(t, err)

	t.Run("idempotent delete created task", func(t *testing.T) {
		for range 5 {
			resp, err := http.DefaultClient.Do(req)
			testhelper.MustNilError(t, err)
			testhelper.MustEqual(t, resp.StatusCode, http.StatusNoContent)
			defer resp.Body.Close()
		}
	})
}

func createTask(t *testing.T) int64 {
	t.Helper()
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

	resp, err := http.Post(httpSrvAddress+"/api/v1/tasks", "application/json", buf)
	testhelper.MustNilError(t, err)
	testhelper.MustEqual(t, resp.StatusCode, http.StatusCreated)

	defer resp.Body.Close()

	var got map[string]any
	err = json.NewDecoder(resp.Body).Decode(&got)
	testhelper.MustNilError(t, err)
	testhelper.MapMustContainKeys(t, got, "id", "title", "description", "status", "priority", "due_date", "category", "tags", "assignee_id", "creator_id", "created_at")
	v, ok := got["id"].(float64)
	if !ok {
		t.Fatalf("expected if be numeric, but is %T", got["id"])
	}
	if v < 1 {
		t.Fatalf("expected id greter than 0, but is %v", v)
	}
	return int64(v)
}

func TestChangeTaskStatus(t *testing.T) {
	id := createTask(t)

	for range 5 {
		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/tasks/%d/status", httpSrvAddress, id), strings.NewReader(`{"status": "inprogress"}`))
		testhelper.MustNilError(t, err)
		resp, err := http.DefaultClient.Do(req)
		testhelper.MustNilError(t, err)
		testhelper.MustEqual(t, resp.StatusCode, http.StatusForbidden)
		defer resp.Body.Close()
	}
}

func TestChangeTaskAssignee(t *testing.T) {
	id := createTask(t)

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/tasks/%d/assignee", httpSrvAddress, id), strings.NewReader(`{"assignee_id":15}`))
	testhelper.MustNilError(t, err)

	t.Run("idempotent change status task", func(t *testing.T) {
		for range 5 {
			resp, err := http.DefaultClient.Do(req)
			testhelper.MustNilError(t, err)
			testhelper.MustEqual(t, resp.StatusCode, http.StatusNoContent)
			defer resp.Body.Close()
		}
	})
}

func TestChangeTaskPriority(t *testing.T) {
	id := createTask(t)

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/tasks/%d/priority", httpSrvAddress, id), strings.NewReader(`{"priority": "medium"}`))
	testhelper.MustNilError(t, err)

	t.Run("idempotent change status task", func(t *testing.T) {
		for range 5 {
			resp, err := http.DefaultClient.Do(req)
			testhelper.MustNilError(t, err)
			testhelper.MustEqual(t, resp.StatusCode, http.StatusNoContent)
			defer resp.Body.Close()
		}
	})
}
