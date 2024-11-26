package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Task interface {
	Execute() error
}

type ScheduledTask struct {
	Interval time.Duration
	Task     Task
	Repeat   bool
}

type Scheduler struct {
	tasks   map[string]ScheduledTask
	mu      sync.Mutex
	wg      sync.WaitGroup
	stopCtx context.Context
	cancel  context.CancelFunc
}

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:   make(map[string]ScheduledTask),
		stopCtx: ctx,
		cancel:  cancel,
	}
}

// handleSignals listens for OS signals and initiates graceful shutdown
func (s *Scheduler) handleSignals() {
	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan) // Ensure cleanup of signal notifications
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-signalChan
	log.Printf("Received signal: %s. Initiating graceful shutdown...", sig)
	s.Stop()
}

// Stop halts the scheduler
func (s *Scheduler) Stop() {
	s.cancel()  // Cancel the context to signal all goroutines to stop
	s.wg.Wait() // Wait for all goroutines to complete
	log.Println("Scheduler stopped gracefully")
}

// AddTask adds a new task to the scheduler
func (s *Scheduler) AddTask(sk ScheduledTask, taskName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[taskName]; exists {
		return errors.New(fmt.Errorf("Task with name '%s' already exists", taskName).Error())
	}

	s.tasks[taskName] = sk
	return nil
}

// Start begins executing scheduled tasks
func (s *Scheduler) Start() {
	go s.handleSignals()

	s.mu.Lock()
	defer s.mu.Unlock()

	for taskName, scheduledTask := range s.tasks {
		s.wg.Add(1)

		go func(name string, t ScheduledTask) {
			defer s.wg.Done()
			ticker := time.NewTicker(t.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := t.Task.Execute(); err != nil {
						log.Printf("Error executing task '%s': %v", name, err)
					}
					if !t.Repeat {
						s.RemoveTask(name)
						return
					}
				case <-s.stopCtx.Done():
					log.Printf("Task '%s' stopped gracefully", name)
					return
				}
			}
		}(taskName, scheduledTask)
	}
}

// RemoveTask removes a task from the scheduler
func (s *Scheduler) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
}
