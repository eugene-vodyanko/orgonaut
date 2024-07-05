package runner_test

import (
	"context"
	"math/rand"

	"github.com/eugene-vodyanko/orgonaut/pkg/runner"
	"testing"
	"time"
)

func TestServer_Run(t *testing.T) {
	checkTimeout := 100
	maxTimeout := 300
	backoffCoefficient := 2
	maxWorkers := 3
	label := "test"

	handler := func(ctx context.Context) (bool, error) {
		dur := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(dur)
		return false, nil
	}

	s := runner.NewRunner(label,
		checkTimeout,
		maxTimeout,
		backoffCoefficient,
		maxWorkers,
		runner.Task{Tag: "task1", Handler: handler},
		runner.Task{Tag: "task2", Handler: handler},
		runner.Task{Tag: "task3", Handler: handler},
		runner.Task{Tag: "task4", Handler: handler},
		runner.Task{Tag: "task5", Handler: handler},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//go func() {
	//	time.Sleep(2000 * time.Millisecond)
	//	cancel()
	//}()

	err := s.RunTasks(ctx)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5000 * time.Millisecond)
	s.Stop(ctx)
}
