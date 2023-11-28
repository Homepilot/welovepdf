package models

import (
	"welovepdf/pkg/utils"
)

type Logger struct {
	customLogger *utils.CustomLogger
}

func NewFrontendLogger(customLogger *utils.CustomLogger) *Logger {
	return &Logger{
		customLogger: customLogger,
	}
}

func (l *Logger) Info(msg string, extras map[string]any) {
	l.customLogger.InfoJson(msg, extras)
}

func (l *Logger) Warn(msg string, extras map[string]any) {
	l.customLogger.WarnJson(msg, extras)
}

func (l *Logger) Error(msg string, extras map[string]any) {
	l.customLogger.ErrorJson(msg, extras)
}

func (l *Logger) Debug(msg string, extras map[string]any) {
	l.customLogger.DebugJson(msg, extras)
}
