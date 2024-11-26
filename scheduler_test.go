package scheduler

import (
	"errors"
	"testing"
	"time"
)

// MockTask is a mock implementation of the Task interface for testing
type MockTask struct {
	ExecutedCount int
	ShouldFail    bool
}

func (m *MockTask) Execute() error {
	m.ExecutedCount++
	if m.ShouldFail {
		return errors.New("mock task failure")
	}
	return nil
}

func TestSchedulerSuite(t *testing.T) {
	t.Run("AddTask", TestScheduler_AddTask)
	t.Run("AddDuplicateTask", TestScheduler_AddDuplicateTask)
	t.Run("StartAndStop", TestScheduler_StartAndStop)
	t.Run("OneTimeTask", TestScheduler_OneTimeTask)
	t.Run("RemoveTask", TestScheduler_RemoveTask)
	t.Run("HandleTaskFailure", TestScheduler_HandleTaskFailure)
}

func TestScheduler_AddTask(t *testing.T) {
	s := NewScheduler()
	taskName := "test_task"
	mockTask := &MockTask{}

	err := s.AddTask(ScheduledTask{
		Interval: time.Second,
		Task:     mockTask,
		Repeat:   true,
	}, taskName)

	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	if _, exists := s.tasks[taskName]; !exists {
		t.Errorf("Task '%s' was not added", taskName)
	}
}

func TestScheduler_AddDuplicateTask(t *testing.T) {
	s := NewScheduler()
	taskName := "duplicate_task"
	mockTask := &MockTask{}

	err := s.AddTask(ScheduledTask{
		Interval: time.Second,
		Task:     mockTask,
		Repeat:   true,
	}, taskName)

	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	err = s.AddTask(ScheduledTask{
		Interval: time.Second,
		Task:     mockTask,
		Repeat:   true,
	}, taskName)

	if err == nil {
		t.Errorf("Expected error when adding duplicate task, got nil")
	}
}

func TestScheduler_StartAndStop(t *testing.T) {
	s := NewScheduler()
	mockTask := &MockTask{}
	taskName := "test_task"

	err := s.AddTask(ScheduledTask{
		Interval: 100 * time.Millisecond,
		Task:     mockTask,
		Repeat:   true,
	}, taskName)

	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	go s.Start()
	time.Sleep(350 * time.Millisecond) // Allow the task to run a few times
	s.Stop()

	if mockTask.ExecutedCount < 2 {
		t.Errorf("Expected task to be executed at least twice, but got %d", mockTask.ExecutedCount)
	}
}

func TestScheduler_OneTimeTask(t *testing.T) {
	s := NewScheduler()
	mockTask := &MockTask{}
	taskName := "one_time_task"

	err := s.AddTask(ScheduledTask{
		Interval: 100 * time.Millisecond,
		Task:     mockTask,
		Repeat:   false,
	}, taskName)

	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	go s.Start()
	time.Sleep(200 * time.Millisecond) // Allow the task to execute once
	s.Stop()

	if mockTask.ExecutedCount != 1 {
		t.Errorf("Expected one-time task to execute exactly once, but got %d", mockTask.ExecutedCount)
	}

	if _, exists := s.tasks[taskName]; exists {
		t.Errorf("Expected one-time task '%s' to be removed after execution", taskName)
	}
}

func TestScheduler_RemoveTask(t *testing.T) {
	s := NewScheduler()
	mockTask := &MockTask{}
	taskName := "removable_task"

	err := s.AddTask(ScheduledTask{
		Interval: time.Second,
		Task:     mockTask,
		Repeat:   true,
	}, taskName)

	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	s.RemoveTask(taskName)

	if _, exists := s.tasks[taskName]; exists {
		t.Errorf("Task '%s' was not removed", taskName)
	}
}

func TestScheduler_HandleTaskFailure(t *testing.T) {
	s := NewScheduler()
	mockTask := &MockTask{ShouldFail: true}
	taskName := "failing_task"

	err := s.AddTask(ScheduledTask{
		Interval: 100 * time.Millisecond,
		Task:     mockTask,
		Repeat:   false,
	}, taskName)

	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	go s.Start()
	time.Sleep(200 * time.Millisecond) // Allow the task to execute once
	s.Stop()

	if mockTask.ExecutedCount != 1 {
		t.Errorf("Expected failing task to execute exactly once, but got %d", mockTask.ExecutedCount)
	}
}
