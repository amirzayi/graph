package handler

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/amirzayi/graph/task"
	"github.com/amirzayi/graph/testhelper"
)

func TestConvertTaskPriorityTextToEnum(t *testing.T) {
	t.Parallel()
	for i, tc := range []struct {
		in       string
		expected task.Priority
	}{
		{" HiGh ", task.PriorityHigh},
		{"high ", task.PriorityHigh},
		{"hi", 0},
		{"t", 0},
		{" medium", task.PriorityMedium},
		{"meDiuM", task.PriorityMedium},
		{"l0w", 0},
		{" lOw", task.PriorityLow},
		{" low ", task.PriorityLow},
	} {
		t.Run(fmt.Sprintf("%s round %d", tc.in, i), func(t *testing.T) {
			testhelper.MustEqual(t, convertTaskPriorityTextToEnum(tc.in), tc.expected)
		})
	}
}

func TestConvertTaskPriorityToText(t *testing.T) {
	t.Parallel()
	for i, tc := range []struct {
		in       string
		expected task.Priority
	}{
		{"HiGh", task.PriorityHigh},
		{"high", task.PriorityHigh},
		{"hi", 0},
		{"t", 0},
		{"medium", task.PriorityMedium},
		{"meDiuM", task.PriorityMedium},
		{"l0w", 0},
		{"lOw", task.PriorityLow},
		{"low", task.PriorityLow},
	} {
		t.Run(fmt.Sprintf("%s round %d", tc.in, i), func(t *testing.T) {
			testhelper.MustNotEqual(t, convertTaskPriorityToText(tc.expected), tc.in)
		})
	}

	for i, tc := range []struct {
		in       string
		expected task.Priority
	}{
		{"High", task.PriorityHigh},
		{"", 0},
		{"Medium", task.PriorityMedium},
		{"Low", task.PriorityLow},
	} {
		t.Run(fmt.Sprintf("%s round %d", tc.in, i), func(t *testing.T) {
			testhelper.MustEqual(t, convertTaskPriorityToText(tc.expected), tc.in)
		})
	}
}

func TestConvertTaskStatusTextToEnum(t *testing.T) {
	t.Parallel()
	for i, tc := range []struct {
		in       string
		expected task.Status
	}{

		{"cancelled", task.StatusCancelled},
		{"high", 0},
		{"hi", 0},
		{"t", 0},
		{"done", task.StatusDone},
		{"meDiuM", 0},
		{"l0w", 0},
		{"inprogress", task.StatusInProgress},
		{"todo", task.StatusTodo},
	} {
		t.Run(fmt.Sprintf("%s round %d", tc.in, i), func(t *testing.T) {
			testhelper.MustEqual(t, convertTaskStatusTextToEnum(tc.in), tc.expected)
		})
	}
}

func TestConvertTaskStatusToText(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		in       string
		expected task.Status
	}{
		{"cancelled", task.StatusCancelled},
		{"high", 0},
		{"hi", 0},
		{"t", 0},
		{"done", task.StatusDone},
		{"meDiuM", 0},
		{"l0w", 0},
		{"inprogress", task.StatusInProgress},
		{"todo", task.StatusTodo},
	} {
		t.Run("", func(t *testing.T) {
			testhelper.MustNotEqual(t, convertTaskStatusToText(tc.expected), tc.in)
		})
	}

	for _, tc := range []struct {
		in       string
		expected task.Status
	}{
		{"Cancelled", task.StatusCancelled},
		{"Done", task.StatusDone},
		{"InProgress", task.StatusInProgress},
		{"Todo", task.StatusTodo},
	} {
		t.Run("", func(t *testing.T) {
			testhelper.MustEqual(t, convertTaskStatusToText(tc.expected), tc.in)
		})
	}
}

func BenchmarkSlice(b *testing.B) {
	tasks := generateTestTasks(100)
	for b.Loop() {
		taskReponse := []taskResponse{}
		for _, t := range tasks {
			taskReponse = append(taskReponse, convertTaskDomainToResponse(t))
		}
	}
}

func BenchmarkPreallocateSlice(b *testing.B) {
	tasks := generateTestTasks(100)
	for b.Loop() {
		taskReponse := make([]taskResponse, 0, len(tasks))
		for _, t := range tasks {
			taskReponse = append(taskReponse, convertTaskDomainToResponse(t))
		}
	}
}

// generateTestTasks creates tasks for use in benchmarks and tests
func generateTestTasks(n int) []task.Task {

	tasks := make([]task.Task, n)
	for i := range n {
		tasks[i] = task.Task{
			ID:          int64(i + 1),
			Title:       fmt.Sprintf("Test Task %d", i+1),
			Description: fmt.Sprintf("Description for test task %d", i+1),
			Status:      statuses[rand.Intn(len(statuses)-1)],
			Priority:    priorities[rand.Intn(len(priorities)-1)],
			DueDate:     time.Now().Add(time.Duration(rand.Intn(30)) * 24 * time.Hour),
			Category:    categories[rand.Intn(len(categories)-1)],
			ParentID:    rand.Int63n(100),
			Tags:        generateTags(rand.Intn(3)),
			AssigneeID:  rand.Intn(100) + 1,
			CreatorID:   rand.Intn(100) + 1,
			CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
		}
	}
	return tasks
}

func generateTags(count int) []string {
	if count == 0 {
		return []string{}
	}
	allTags := []string{"urgent", "important", "backlog", "sprint", "bug", "feature"}
	tags := make([]string, count)
	for i := range count {
		tags[i] = allTags[rand.Intn(len(allTags))]
	}
	return tags
}

var (
	statuses   = []task.Status{task.StatusTodo, task.StatusInProgress, task.StatusDone, task.StatusCancelled}
	priorities = []task.Priority{task.PriorityLow, task.PriorityMedium, task.PriorityHigh}
	categories = []string{"Work", "Personal", "Shopping", "Health", "Education", "Finance"}
)
