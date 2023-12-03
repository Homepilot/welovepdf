package models

import (
	"encoding/json"
	"log/slog"
)

type FrontendLogger struct{}

func (l *FrontendLogger) Info(msg string, extraJson string) {
	if extraJson == "" {
		slog.Info(msg)
		return
	}
	slog.Info(msg, formatJsonStringToSlogArgs(extraJson)...)
}

func (l *FrontendLogger) Warn(msg string, extraJson string) {
	if extraJson == "" {
		slog.Warn(msg)
		return
	}
	slog.Warn(msg, formatJsonStringToSlogArgs(extraJson)...)
}

func (l *FrontendLogger) Error(msg string, extraJson string) {
	if extraJson == "" {
		slog.Error(msg)
		return
	}
	slog.Error(msg, formatJsonStringToSlogArgs(extraJson)...)
}

func (l *FrontendLogger) Debug(msg string, extraJson string) {
	if extraJson == "" {
		slog.Debug(msg)
		return
	}
	slog.Debug(msg, formatJsonStringToSlogArgs(extraJson)...)
}

func formatJsonStringToSlogArgs(jsonString string) []any {
	var dataMap map[string]any
	err := json.Unmarshal([]byte(jsonString), &dataMap)
	if err != nil {
		slog.Debug("Error parsing json string", slog.String("reason", err.Error()))
		return []any{slog.String("extraJson", jsonString)}
	}

	attr := []any{}
	for key, value := range dataMap {
		safeStringValue, isString := value.(string)
		if isString {
			attr = append(attr, slog.String(key, safeStringValue))
		}
		safeIntValue, isInt := value.(int)
		if isInt {
			attr = append(attr, slog.Int(key, safeIntValue))
		}
		safeFloatValue, isFloat := value.(float64)
		if isFloat {
			attr = append(attr, slog.Float64(key, safeFloatValue))
		}
		safeBoolValue, isBool := value.(bool)
		if isBool {
			attr = append(attr, slog.Bool(key, safeBoolValue))
		}
	}

	return attr
}
