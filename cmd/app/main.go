package main

import (
	"flag"
	"github.com/eugene-vodyanko/orgonaut/internal/app"
	"github.com/eugene-vodyanko/orgonaut/internal/config"
	"github.com/eugene-vodyanko/orgonaut/pkg/logger"
	"github.com/eugene-vodyanko/orgonaut/pkg/util"

	"log"
)

var configPath string

func main() {
	flag.Parse()
	log.Println("main - cfg path:", configPath)

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("main - cfg error: %s", err)
	}

	teardown, err := logger.SetupDefaultLogger(
		cfg.Logger.Root.LogLevel,
		cfg.Logger.File.Name,
		cfg.Logger.File.Format == "JSON",
	)

	if err != nil {
		log.Fatalf("main - log setup error: %s", err)
	}

	defer func() {
		err = teardown()
		if err != nil {
			log.Fatalf("main - log teardown error: %s", err)
		}
	}()

	if err = app.Run(cfg); err != nil {
		log.Fatal(err)
	}

	util.PrintResourceUsage()
}

func init() {
	flag.StringVar(&configPath, "config-path", "configs/application.yml", "path to config file")
}
