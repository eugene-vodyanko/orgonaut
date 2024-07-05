package util

import (
	"log/slog"
	"time"
)

// Timer returns a function that log (slog.Info()) the name argument and
// the elapsed time between the call to Timer and the call to
// the returned function. The returned function is intended to
// be used in a defer statement:
//
//	defer Timer("foo")()
func Timer(name string) func() {
	start := time.Now()
	return func() {
		slog.Info("timer", "name", name, "elapsed", time.Since(start))
	}
}
