package main

import (
	"solod.dev/so/log/slog"
	"solod.dev/so/os"
)

func main() {
	logger()
	defaults()
}

func logger() {
	// Logger writing to stdout.
	h := slog.NewTextHandler(os.Stdout, slog.LevelInfo)
	l := slog.New(&h)

	// Enabled check.
	if l.Enabled(slog.LevelDebug) {
		panic("debug should not be enabled")
	}
	if !l.Enabled(slog.LevelInfo) {
		panic("info should be enabled")
	}

	// Log at info - should appear.
	l.Info("hello world", slog.String("user", "john"), slog.Int("count", 42))

	// Log at debug - should be filtered.
	l.Debug("hidden")

	// Log with no attrs.
	l.Warn("caution")

	// Log with float and bool attrs.
	l.Error("failure", slog.Float64("elapsed", 1.5), slog.Bool("retry", true))

	// Log with string that needs quoting.
	l.Info("test quoting", slog.String("msg", "hello world"))
}

func defaults() {
	// Default logger should be usable.
	slog.Info("default test", slog.Int("port", 8080))
}
