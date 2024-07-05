package logger_test

import (
	"errors"
	"fmt"
	"github.com/eugene-vodyanko/orgonaut/pkg/logger"
	"log/slog"
	"testing"
)

func TestLogger_Debug(t *testing.T) {

	teardown, err := logger.SetupDefaultLogger("Debug", "", false)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = teardown()
		if err != nil {
			t.Fatal(err)
		}
	}()

	slog.Debug("this is debug message", "str_param", "str_value", "int_param", 42)
}

func TestLogger_Error(t *testing.T) {

	teardown, err := logger.SetupDefaultLogger("Debug", "", false)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = teardown()
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = errors.New("something happened")
	slog.Error("this is error message", "error", fmt.Errorf("error: %w", err))
	slog.Error("this is yet another error message", slog.Any("error", err))
}
