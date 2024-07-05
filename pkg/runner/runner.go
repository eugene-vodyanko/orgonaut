package runner

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

const (
	_defaultInitialInterval    = 100
	_defaultMaxInterval        = 1000
	_defaultBackoffCoefficient = 2
	_defaultMaxWorkers         = 1
)

type TaskHandler func(ctx context.Context) (bool, error)
type Task struct {
	Tag     string
	Handler TaskHandler
}
type Runner struct {
	name               string
	initialInterval    int
	maxInterval        int
	backoffCoefficient int
	maxWorkers         int
	tasks              []Task

	wg     *sync.WaitGroup
	sema   chan struct{}
	cancel context.CancelFunc
}

// NewRunner creates an instance of the task executor.
// The label can be used to separate different performers in the logs.
// Task handlers are started and stopped using the RunTasks and Stop methods.
func NewRunner(label string, initialInterval, maxInterval, backoffCoefficient, maxWorkers int, tasks ...Task) *Runner {

	name := "runner"
	if label != "" {
		name = name + ":" + label
	}

	if initialInterval <= 0 {
		initialInterval = _defaultInitialInterval
	}

	if maxInterval <= 0 {
		maxInterval = _defaultMaxInterval
	}

	if backoffCoefficient <= 0 {
		backoffCoefficient = _defaultBackoffCoefficient
	}

	if maxWorkers <= 0 {
		maxWorkers = _defaultMaxWorkers
	}

	return &Runner{
		name:               name,
		initialInterval:    initialInterval,
		maxInterval:        maxInterval,
		backoffCoefficient: backoffCoefficient,
		maxWorkers:         maxWorkers,
		tasks:              tasks,
		sema:               make(chan struct{}, maxWorkers),
	}
}

// RunTasks run each task in a separate goroutine and returns control.
// To control the degree of parallelism with a large number of tasks,
// a semaphore of the size "maxWorkers" is used.
// Each task is started immediately if the handler returns true,
// and is sent to wait (with increasing interval) if the handler returns false
func (r *Runner) RunTasks(ctx context.Context) error {

	slog.Info(fmt.Sprintf("%s - run tasks", r.name),
		"initial_interval", r.initialInterval,
		"max_interval", r.maxInterval,
		"backoff_coefficient", r.backoffCoefficient,
		"max_workers", r.maxWorkers,
	)

	ctx, cancel := context.WithCancel(ctx)
	r.cancel = cancel

	var wg sync.WaitGroup
	r.wg = &wg

	for _, v := range r.tasks {
		v := v
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.runTask(ctx, v.Tag, v.Handler)

		}()
	}

	return nil
}

func (r *Runner) runTask(ctx context.Context, tag string, handler TaskHandler) {

	slog.Info(fmt.Sprintf("%s - run", r.name), "task", tag)

	var timeout int
	timeout = 0
	for {
		select {
		case <-ctx.Done():
			slog.Debug(fmt.Sprintf("%s[%s] - cancel signal has been received", r.name, tag))
			return
		default:
			if timeout != 0 {
				wait(ctx, timeout)
			}

			success, err := r.boundedHandler(ctx, handler)
			if err != nil {
				slog.Error(
					fmt.Sprintf("%s[%s] - call handler error", r.name, tag), "err", err,
				)
			}

			if success {
				timeout = 0
			} else {
				if timeout == 0 {
					timeout = r.initialInterval
				} else {
					timeout = min(timeout*r.backoffCoefficient, r.maxInterval)
				}
				slog.Debug(fmt.Sprintf("%s[%s] - increasing timeout up to %d ms", r.name, tag, timeout))
			}
		}
	}
}

func wait(ctx context.Context, timeout int) {
	if timeout < 1 {
		timeout = 1
	}

	select {
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		return
	case <-ctx.Done():
		return
	}

}

func (r *Runner) boundedHandler(ctx context.Context, handler TaskHandler) (bool, error) {

	select {
	case r.sema <- struct{}{}:
		defer func() {
			<-r.sema
		}()
	case <-ctx.Done():

		return false, nil
	}
	return handler(context.Background())
}

// Stop stops all task worker
func (r *Runner) Stop(ctx context.Context) {
	slog.Info(fmt.Sprintf("%s - stop, releasing resources", r.name))
	r.cancel()
	r.wg.Wait()
	slog.Info(fmt.Sprintf("%s - stop, done", r.name))
}
