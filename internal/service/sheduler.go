package service

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"marketflow/internal/domain/types"
	"marketflow/pkg/logger"
)

// Task представляет задачу, которую нужно периодически выполнять
type Task struct {
	Name string
	Type types.TaskType

	Interval time.Duration

	Handler func(ctx context.Context) error
}

// Scheduler представляет планировщик задач
type Scheduler struct {
	ctx     context.Context
	cancel  context.CancelFunc
	tasks   []Task
	logger  logger.Logger
	wg      sync.WaitGroup
	started bool
	mu      sync.Mutex
}

func NewScheduler(ctx context.Context, logger logger.Logger) *Scheduler {
	schedulerCtx, cancel := context.WithCancel(ctx)

	return &Scheduler{
		ctx:    schedulerCtx,
		cancel: cancel,
		tasks:  make([]Task, 0),
		logger: logger,
	}
}

// Start запускает планировщик задач
func (s *Scheduler) Start() {
	if len(s.tasks) == 0 {
		s.logger.Info(s.ctx, "no tasks to start in sheduler")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		s.logger.Warn(s.ctx, "Scheduler already started")
		return
	}

	s.logger.Info(s.ctx, "Starting scheduler")

	// Запускаем каждую задачу в отдельной горутине
	for _, task := range s.tasks {
		s.wg.Add(1)
		switch task.Type {
		case types.TaskTypeInterval:
			go s.runTask(task)
		default:
			s.logger.Warn(s.ctx, "failed to start task, invalid task type", "name", task.Name, "type", task.Type)
		}
	}

	s.started = true
}

// runTask запускает отдельную задачу в цикле с заданным интервалом
func (s *Scheduler) runTask(task Task) {
	defer s.wg.Done()

	log := s.logger.GetSlogLogger().With("task name", task.Name)

	defer func() {
		if r := recover(); r != nil {
			log.ErrorContext(s.ctx, "Task panicked",
				"panic", r,
				"stack", debug.Stack())
		}
	}()

	log.InfoContext(s.ctx, "Starting task", "interval", task.Interval)

	// Первый запуск задачи сразу после старта
	if err := task.Handler(s.ctx); err != nil {
		log.ErrorContext(s.ctx, "Failed to execute task", "error", err)
	}

	// Создаем тикер для периодического запуска
	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.DebugContext(s.ctx, "Executing task")

			// Создаем тайм-аут для выполнения задачи
			taskCtx, cancel := context.WithTimeout(s.ctx, task.Interval/2)

			if err := task.Handler(taskCtx); err != nil {
				log.ErrorContext(s.ctx, "Failed to execute task", "error", err)
			}

			cancel()

		case <-s.ctx.Done():
			log.InfoContext(s.ctx, "Stopping task")
			return
		}
	}
}

// Close останавливает планировщик задач
func (s *Scheduler) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return
	}

	s.logger.Info(s.ctx, "Stopping scheduler")
	s.cancel()

	// Ожидаем завершения всех задач с тайм-аутом
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info(s.ctx, "All scheduler tasks stopped")
	case <-time.After(5 * time.Second):
		s.logger.Warn(s.ctx, "Scheduler tasks shutdown timed out")
	}

	s.started = false
}

// AddTask добавляет новую задачу в шедулер
// Задача будет запущена только после перезапуска шедулера
func (s *Scheduler) AddTask(name string, taskType types.TaskType, interval time.Duration, handler func(ctx context.Context) error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks = append(s.tasks, Task{
		Name:     name,
		Type:     taskType,
		Interval: interval,
		Handler:  handler,
	})

	s.logger.Info(s.ctx, "Task added", "name", name, "interval", interval)
}
