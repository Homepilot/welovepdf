package models

import (
	"log/slog"
)

type FrontendLogger struct{}

func (l *FrontendLogger) Info(msg string, extras map[string]any) {
	slog.Info(msg, extras)
}

func (l *FrontendLogger) Warn(msg string, extras map[string]any) {
	slog.Warn(msg, extras)
}

func (l *FrontendLogger) Error(msg string, extras map[string]any) {
	slog.Error(msg, extras)
}

func (l *FrontendLogger) Debug(msg string, extras map[string]any) {
	slog.Debug(msg, extras)
}
