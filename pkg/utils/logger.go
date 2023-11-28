package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TODO : https://go.dev/play/p/0yJNk065ftB

type CustomLogger struct {
	logger       *slog.Logger
	lumberjack   *lumberjack.Logger
	logtailToken string
	logsDirPath  string
	sysInfo      map[string]any
}

func NewLogger(logsDir string, logtailToken string) *CustomLogger {
	removeEmptyLogsFiles(logsDir)
	logFilePath := path.Join(logsDir, uuid.NewString()+".log")

	lj := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}

	logger := slog.New(slogmulti.Fanout(
		slog.NewJSONHandler(lj, &slog.HandlerOptions{AddSource: true}),
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}),
	))

	slog.SetDefault(logger)

	return &CustomLogger{
		logger:       logger,
		lumberjack:   lj,
		logtailToken: logtailToken,
		logsDirPath:  logsDir,
		sysInfo:      getSysInfo(),
	}
}

func (c *CustomLogger) Close() {
	c.lumberjack.Close()
	removeEmptyLogsFiles(c.logsDirPath)
}

func (c *CustomLogger) InfoJson(msg string, argsMap map[string]any) {
	c.logger.Info(msg, []slog.Attr{})
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(c.getJsonBodyFromMap(msg, slog.LevelInfo, argsMap))
}

func (c *CustomLogger) WarnJson(msg string, argsMap map[string]any) {
	c.logger.Warn(msg, []slog.Attr{})
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(c.getJsonBodyFromMap(msg, slog.LevelWarn, argsMap))
}

func (c *CustomLogger) ErrorJson(msg string, argsMap map[string]any) {
	c.logger.Error(msg, []slog.Attr{})
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(c.getJsonBodyFromMap(msg, slog.LevelError, argsMap))
}

func (c *CustomLogger) DebugJson(msg string, argsMap map[string]any) {
	c.logger.Debug(msg, []slog.Attr{})
}

func (c *CustomLogger) Info(msg string, args ...slog.Attr) {
	c.logger.Info(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(c.getJsonBodyFromSlogArgs(msg, slog.LevelInfo, args))
}

func (c *CustomLogger) Warn(msg string, args ...slog.Attr) {
	c.logger.Warn(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(c.getJsonBodyFromSlogArgs(msg, slog.LevelWarn, args))
}

func (c *CustomLogger) Error(msg string, args ...slog.Attr) {
	c.logger.Error(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(c.getJsonBodyFromSlogArgs(msg, slog.LevelError, args))
}

func (c *CustomLogger) Debug(msg string, args ...slog.Attr) {
	c.logger.Debug(msg, args)
}

func (c *CustomLogger) log(msg string, level slog.Level, args []slog.Attr) {
	if c.logtailToken != "" {
		_ = c.sendLogToLogtail(c.getJsonBodyFromSlogArgs(msg, slog.LevelDebug, args))
	}

	switch level {
	case slog.LevelError:
		c.logger.Error(msg, args)
	case slog.LevelWarn:
		c.logger.Warn(msg, args)
	case slog.LevelInfo:
		c.logger.Info(msg, args)
	default:
		c.logger.Debug(msg, args)
	}

}

func (c *CustomLogger) sendLogToLogtail(body *bytes.Buffer) error {
	logTailUrl := "https://in.logs.betterstack.com/"

	req, err := http.NewRequest(http.MethodPost, logTailUrl, body)
	if err != nil {
		fmt.Printf("Sending to LOGTAIL failed at creating req : %s", err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.logtailToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Sending to LOGTAIL failed at sending res: %s", err.Error())
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 400 {
		return nil
	}
	return fmt.Errorf("Response has status : %s", res.Status)
}

func (c *CustomLogger) getBaseLogBody(msg string, level slog.Level) map[string]any {
	logObj := map[string]any{
		"message": msg,
		"level":   level,
		"dt":      time.Now().Unix(),
	}
	for k, v := range c.sysInfo {
		logObj[k] = v
	}
	return logObj
}

func (c *CustomLogger) getJsonBodyFromSlogArgs(msg string, level slog.Level, args []slog.Attr) *bytes.Buffer {
	logObj := c.getBaseLogBody(msg, level)
	for i := 0; i < len(args); i = 1 {
		logObj[args[i].Key] = args[i].Value.String()
	}

	data, err := json.Marshal(logObj)
	if err != nil {
		fmt.Printf("Error formatting JSON body : %s", err.Error())

		body := []byte(`{ "level": "` + level.String() + `", "message": "` + msg + `" }`)
		return bytes.NewBuffer(body)
	}

	body := []byte(data)
	return bytes.NewBuffer(body)
}

func (c *CustomLogger) getJsonBodyFromMap(msg string, level slog.Level, argsMap map[string]any) *bytes.Buffer {
	logObj := c.getBaseLogBody(msg, level)
	for k, v := range argsMap {
		logObj[k] = v
	}

	data, err := json.Marshal(logObj)
	if err != nil {
		fmt.Printf("Error formatting JSON body : %s", err.Error())

		body := []byte(`{ "level": "` + level.String() + `", "message": "` + msg + `" }`)
		return bytes.NewBuffer(body)
	}

	body := []byte(data)
	return bytes.NewBuffer(body)
}

func removeEmptyLogsFiles(tempDir string) {
	logsDir := path.Join(tempDir, "..", "logs")
	logsFiles, err := os.ReadDir(logsDir)
	if err != nil {
		return
	}
	for i := 0; i < len(logsFiles); i += 1 {
		logFile := logsFiles[i]
		fileInfo, err := logFile.Info()
		if err == nil && !logFile.IsDir() && fileInfo.Size() == 0 {
			_ = os.RemoveAll(path.Join(logsDir, fileInfo.Name()))
		}
	}
}
