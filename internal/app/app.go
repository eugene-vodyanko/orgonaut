package app

import (
	"context"
	"fmt"
	"github.com/eugene-vodyanko/orgonaut/internal/config"
	"github.com/eugene-vodyanko/orgonaut/internal/controller/task"
	"github.com/eugene-vodyanko/orgonaut/internal/infrastructure/broker"
	"github.com/eugene-vodyanko/orgonaut/internal/infrastructure/repository"
	"github.com/eugene-vodyanko/orgonaut/internal/service"
	"github.com/eugene-vodyanko/orgonaut/pkg/kafka/kafkakit"
	"github.com/eugene-vodyanko/orgonaut/pkg/oracle"
	"github.com/eugene-vodyanko/orgonaut/pkg/runner"
	"github.com/eugene-vodyanko/orgonaut/pkg/util"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) error {

	// Init DB
	ora, err := oracle.New(
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.URL,
		cfg.DB.Schema,
		cfg.DB.Pool.MaxOpenConns,
		cfg.DB.Pool.MaxIdleConns,
		cfg.DB.Pool.MaxLifetime,
		cfg.DB.Pool.MaxIdleTime,
	)
	if err != nil {
		log.Fatal(fmt.Errorf("app - oracle init error: %w", err))
	}

	defer func() {
		err := ora.Close()
		if err != nil {
			log.Fatal(fmt.Errorf("app - oracle close error: %w", err))
		}
	}()

	// Init Kafka writer
	writer, err := kafkakit.NewWriter(cfg.Kafka.Brokers, "",
		cfg.Kafka.Compress,
		cfg.Kafka.BatchSize,
		time.Duration(cfg.Kafka.BatchTimeout)*time.Millisecond,
		cfg.Kafka.RequiredAcks,
		cfg.Kafka.CreateTopic,
		cfg.Kafka.MaxReqSize,
	)

	if err != nil {
		log.Fatal(fmt.Errorf("app - kafka writer init error: %w", err))
	}

	// Init service
	srv := service.New(
		repository.NewRepository(cfg.DB.Schema, ora),
		repository.NewTxManager(ora.Db),
		broker.NewBroker(writer),
	)

	// Init routes
	routes, err := task.NewRoutes(cfg.Tasks, srv)
	if err != nil {
		log.Fatal(fmt.Errorf("app - routes init error: %w", err))
	}

	// Init runner
	r := runner.NewRunner("",
		cfg.Runner.RepeatPolicy.InitialInterval,
		cfg.Runner.RepeatPolicy.MaxInterval,
		cfg.Runner.RepeatPolicy.BackoffCoefficient,
		cfg.Runner.MaxWorkers,
		routes...,
	)

	// Run tasks
	defer util.Timer("uptime")()

	ctx := context.Background()
	err = r.RunTasks(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("app - run tasks error: %w", err))
	}

	// Wait stop signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	r.Stop(ctx)

	return nil
}
