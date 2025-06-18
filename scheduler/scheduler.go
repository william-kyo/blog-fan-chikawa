package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Scheduler struct {
	tasks  map[string]*ScheduledTask
	mutex  sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

type ScheduledTask struct {
	name     string
	interval time.Duration
	fn       func()
	ticker   *time.Ticker
	done     chan bool
}

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:  make(map[string]*ScheduledTask),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Add fixed rate schedule task
func (s *Scheduler) ScheduleAtFixedRate(name string, fn func(), interval time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if existingTask, exists := s.tasks[name]; exists {
		s.stopTask(existingTask)
	}

	task := &ScheduledTask{
		name:     name,
		interval: interval,
		fn:       fn,
		ticker:   time.NewTicker(interval),
		done:     make(chan bool),
	}

	s.tasks[name] = task
	s.startTask(task)
}

func (s *Scheduler) startTask(task *ScheduledTask) {
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-task.done:
				return
			case <-task.ticker.C:
				fmt.Printf("[%s] executing task: %v\n", task.name, time.Now().Format("15:04:05"))
				task.fn()
			}
		}
	}()
}

func (s *Scheduler) stopTask(task *ScheduledTask) {
	if task.ticker != nil {
		task.ticker.Stop()
	}
	close(task.done)
}

// Stop a task
func (s *Scheduler) StopTask(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if task, exists := s.tasks[name]; exists {
		s.stopTask(task)
		delete(s.tasks, name)
	}
}

// Stop all tasks
func (s *Scheduler) Shutdown() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, task := range s.tasks {
		s.stopTask(task)
	}
	s.cancel()
	s.tasks = make(map[string]*ScheduledTask)
}
