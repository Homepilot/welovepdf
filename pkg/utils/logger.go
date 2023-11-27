package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"

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

func (c *CustomLogger) Info(msg string, args ...slog.Attr) {
	c.logger.Info(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(msg, slog.LevelInfo, args)
}

func (c *CustomLogger) Warn(msg string, args ...slog.Attr) {
	c.logger.Warn(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(msg, slog.LevelWarn, args)
}

func (c *CustomLogger) Error(msg string, args ...slog.Attr) {
	c.logger.Error(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(msg, slog.LevelError, args)
}

func (c *CustomLogger) Debug(msg string, args ...slog.Attr) {
	c.logger.Debug(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(msg, slog.LevelDebug, args)
}

func (c *CustomLogger) log(msg string, level slog.Level, args []slog.Attr) {
	if c.logtailToken != "" {
		_ = c.sendLogToLogtail(msg, slog.LevelDebug, args)
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

func (c *CustomLogger) sendLogToLogtail(msg string, level slog.Level, args []slog.Attr) error {
	logTailUrl := "https://in.logs.betterstack.com/"

	req, err := http.NewRequest(http.MethodPost, logTailUrl, c.getJsonBodyFromArgs(msg, level, args))
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
		fmt.Printf("LOGTAIL req status : %s", res.Status)
		return nil
	}
	return fmt.Errorf("Response has status : %s", res.Status)
}

func (c *CustomLogger) getJsonBodyFromArgs(msg string, level slog.Level, args []slog.Attr) *bytes.Buffer {
	logObj := map[string]any{
		"message": msg,
		"level":   level,
	}
	for k, v := range c.sysInfo {
		logObj[k] = v
	}

	for i := 0; i < len(args); i = 1 {
		logObj[args[i].Key] = args[i].Value.String()
		fmt.Printf("adding args w/ name : %s and value : %s", args[i].Key, args[i].String())
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
