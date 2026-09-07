package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	db.AutoMigrate(&models.Task{})

	router.Bootstrap(ginrouter, router.Dependencies{
		TaskService: task.NewService(task.NewSQLRepository(db), audit.NewDiscarded()),
	})

	m.Run()
}

func TestListTask(t *testing.T) {
	resp, err := http.Get(httpSrvAddress + "/api/v1/tasks")
	testhelper.MustNilError(t, err)
	testhelper.MustEqual(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	var got map[string]any
	err = json.NewDecoder(resp.Body).Decode(&got)
	testhelper.MustNilError(t, err)
	testhelper.MapMustContainKeys(t, got, "total", "data")

	total, _ := got["total"]

	if v, ok := total.(float64); true {
		if !ok {
			t.Fatalf("expected numeric total, but is %T", total)
		}
		if v != 0 {
			t.Fatalf("expected zero value total, but is %f", v)
		}
	}
}

func TestCreateTask(t *testing.T) {
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
}

func TestGetTask(t *testing.T) {
	t.Run("get created task", func(t *testing.T) {
		t.Run("create task", TestCreateTask)
		resp, err := http.Get(httpSrvAddress + "/api/v1/tasks/1")
		testhelper.MustNilError(t, err)
		testhelper.MustEqual(t, resp.StatusCode, http.StatusOK)
		defer resp.Body.Close()

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
	t.Run("create task", TestCreateTask)

	req, err := http.NewRequest(http.MethodDelete, httpSrvAddress+"/api/v1/tasks/1", http.NoBody)
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
