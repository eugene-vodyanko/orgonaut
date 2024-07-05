package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

// SetupDefaultLogger setups the logging level and output for default slog logger.
// The logging level is passed as a string (case-insensitive). For example: "Debug", "ERROR".
// If an empty file name is passed, logging is performed in a compact text format in stdout.
// Output to the file is made in a structured format: text or JSON (depending on useJsonFmt).
func SetupDefaultLogger(level string, filename string, useJsonFmt bool) (close func() error, err error) {
	lvl, err := parseLevel(level)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to setup log level: %w", err))
	}

	if filename != "" {
		return setupFileLogger(lvl, filename, useJsonFmt)
	} else {
		setupStdoutLogger(lvl)
	}

	return func() error { return nil }, nil
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))

	return level, err
}

func setupStdoutLogger(level slog.Level) {
	slog.SetLogLoggerLevel(level)
}
func setupFileLogger(level slog.Level, filename string, json bool) (func() error, error) {
	file, closeFile, err := openFile(filename)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if json {
		handler = slog.NewJSONHandler(file, opts)
	} else {
		handler = slog.NewTextHandler(file, opts)
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return closeFile, nil
}

func openFile(filename string) (*os.File, func() error, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup log to file: %w", err)
	}

	closeFile := func() error {
		if file != nil {
			err := file.Close()
			if err != nil {
				return fmt.Errorf("failed to close log file: %w", err)
			}
		}

		return nil
	}

	return file, closeFile, nil
}
