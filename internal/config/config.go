package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type (
	Config struct {
		Logger Logger          `yaml:"logging"`
		DB     Datasource      `yaml:"datasource"`
		Kafka  Kafka           `yaml:"kafka"`
		Runner Runner          `yaml:"runner"`
		Tasks  map[string]Task `yaml:"tasks"`
	}

	Logger struct {
		Root struct {
			LogLevel string `yaml:"root"`
		} `yaml:"level"`
		File struct {
			Name   string `yaml:"name"`
			Format string `yaml:"format"`
		} `yaml:"file"`
	}

	Datasource struct {
		URL      string `yaml:"url"`
		Schema   string `yaml:"schema"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`

		Pool struct {
			MaxOpenConns int `yaml:"max_open_conns"`
			MaxIdleConns int `yaml:"max_idle_conns"`
			MaxLifetime  int `yaml:"max_life_time"`
			MaxIdleTime  int `yaml:"max_idle_time"`
		} `yaml:"connection_pool"`
	}

	Kafka struct {
		Brokers      []string `yaml:"brokers,flow"`
		BatchSize    int      `yaml:"batch_size"`
		BatchTimeout int      `yaml:"batch_timeout"`
		RequiredAcks string   `yaml:"required_acks"`
		Compress     bool     `yaml:"compress"`
		CreateTopic  bool     `yaml:"topic_auto_create"`
		MaxReqSize   int64    `yaml:"max_request_size"`
	}

	Runner struct {
		MaxWorkers int `yaml:"max_workers"`

		RepeatPolicy struct {
			MaxInterval        int `yaml:"max_interval"`
			InitialInterval    int `yaml:"initial_interval"`
			BackoffCoefficient int `yaml:"backoff_coefficient"`
		} `yaml:"repeat_policy"`
	}

	Task struct {
		GroupId   string `yaml:"group_id"`
		PartCount int    `yaml:"part_count"`
		BatchSize int    `yaml:"batch_size"`
		Topic     string `yaml:"topic"`

		Query struct {
			Columns  string `yaml:"columns"`
			From     string `yaml:"from"`
			PkColumn string `yaml:"pk_column"`
		} `yaml:"query"`
	}
)

// NewConfig read configuration
func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("config file open error: %w", err)
	}

	defer func() {
		err = configFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	decoder := yaml.NewDecoder(configFile)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config decoding error: %w", err)
	}

	return cfg, nil
}
