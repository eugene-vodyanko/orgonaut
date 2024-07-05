package task

import (
	"context"
	"fmt"
	"github.com/eugene-vodyanko/orgonaut/internal/config"
	"github.com/eugene-vodyanko/orgonaut/internal/model"
	"github.com/eugene-vodyanko/orgonaut/pkg/runner"
	"log/slog"

	"github.com/eugene-vodyanko/orgonaut/internal/service"
)

type router struct {
	srv service.Relayer
}

// NewRoutes sets up handlers for the provided configuration.
// The logic provides an approach: one job for one part of the one task (e.g. for
// each certain part_id: from 0 to part_count - 1).
// In other words, the total number of jobs is equal to the sum of all the part_count jobs.
func NewRoutes(tasks map[string]config.Task, s service.Relayer) ([]runner.Task, error) {
	r := &router{srv: s}

	var task []runner.Task

	for k, v := range tasks {
		for i := 0; i < v.PartCount; i++ {
			t := &model.Task{
				BatchSize: v.BatchSize,
				GroupId:   v.GroupId,
				PartId:    i,
			}

			t.Query.From = v.Query.From
			t.Query.Columns = v.Query.Columns
			t.Query.PkColumn = v.Query.PkColumn
			t.Topic = v.Topic

			err := t.Validate()
			if err != nil {
				return nil, fmt.Errorf("router - task[%s] validation error: %w", k, err)
			}

			task = append(task, r.newRoute(t))
		}
	}
	return task, nil
}
func (r *router) newRoute(task *model.Task) runner.Task {
	tag := fmt.Sprintf("task_%s_%d", task.GroupId, task.PartId)

	return runner.Task{
		Tag:     tag,
		Handler: r.newTaskHandler(task, tag),
	}
}
func (r *router) newTaskHandler(task *model.Task, tag string) runner.TaskHandler {
	return func(ctx context.Context) (bool, error) {
		slog.Debug(fmt.Sprintf("handler[%s] - handle next records", tag))

		amount, err := r.srv.Relay(ctx, task)
		if err != nil {
			return false, fmt.Errorf("handler - processing error: %w", err)
		}

		slog.Debug(fmt.Sprintf("handler[%s] - relay done", tag), "sent_amount", amount)

		return amount > 0, nil
	}
}
