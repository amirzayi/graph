package handler

import (
	"fmt"
	"testing"

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
