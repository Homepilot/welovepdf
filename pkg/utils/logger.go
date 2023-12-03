package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TODO ADD SYS INFO ATTRIBUTES to logs
// TODO REPLACE ATTRIBUTES: https://go.dev/play/p/0yJNk065ftB

type CustomLogger struct {
	lumberjack    *lumberjack.Logger
	logtailToken  string
	logsDirPath   string
	logsToSend    [][]byte
	logsBatchSize int
}

func SetupLogger(logsDir string, logtailToken string, logLevel slog.Level) *CustomLogger {
	fmt.Printf("LOGTAIL TOKEN : %s", logtailToken)
	removeEmptyLogsFiles2(logsDir)
	fileName := strings.Join(strings.Split(time.Now().Local().Format(time.DateTime), " "), "") + ".log"
	logFilePath := path.Join(logsDir, fileName)

	lj := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}

	logTailLogLevel := logLevel
	if logTailLogLevel == slog.LevelDebug {
		logTailLogLevel = slog.LevelInfo
	}

	customLogger := &CustomLogger{
		lumberjack:   lj,
		logtailToken: logtailToken,
		logsDirPath:  logsDir,
	}

	slogHandlers := slogmulti.Fanout(
		slog.NewTextHandler(lj, &slog.HandlerOptions{AddSource: true, Level: logLevel}),
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: logLevel}),
		slog.NewJSONHandler(customLogger, &slog.HandlerOptions{Level: logLevel}),
	)
	if logtailToken == "" {
		slogHandlers = slogmulti.Fanout(
			slog.NewTextHandler(lj, &slog.HandlerOptions{AddSource: true, Level: logLevel}),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: logLevel}),
		)
	}

	slogLogger := slog.New(slogHandlers)
	slog.SetDefault(slogLogger)

	return customLogger
}

// Should only be called by a slog json handler
func (l *CustomLogger) Write(data []byte) (int, error) {
	fmt.Printf("GOT LOG TO SEND !!, data length : %d", len(data))
	l.logsToSend = append(l.logsToSend, data)
	// if len(l.logsToSend) >= l.logsBatchSize {
	// 	_ = l.flush()
	// }
	l.sendToLogtail(data)
	return len(data), nil
}

// func (l *CustomLogger) flush() error {
// 	err := l.sendCurrentBatchToLogtail()
// 	if err != nil {
// 		fmt.Printf("Error sending to logtail : %s", err.Error())
// 	}
// 	return err
// }

func (c *CustomLogger) Close() {
	c.lumberjack.Close()
	removeEmptyLogsFiles2(c.logsDirPath)
}

func (c *CustomLogger) sendToLogtail(body []byte) error {
	logTailUrl := "https://in.logs.betterstack.com/"

	req, err := http.NewRequest(http.MethodPost, logTailUrl, bytes.NewBuffer(body))
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

// func mergeJsonAsBytesArrayToJsonBytes(bytesArrays *[][]byte) ([]byte, error) {
// 	arraysToMerge := *bytesArrays

// 	return []byte{}, nil
// }

// func (c *CustomLogger) getBaseLogBody(msg string, level slog.Level) map[string]any {
// 	logObj := map[string]any{
// 		"message": msg,
// 		"level":   level,
// 		"dt":      time.Now().Unix(),
// 	}
// 	for k, v := range c.sysInfo {
// 		logObj[k] = v
// 	}
// 	return logObj
// }

// func (c *CustomLogger) getJsonBodyFromSlogArgs(msg string, level slog.Level, args []slog.Attr) *bytes.Buffer {
// 	logObj := c.getBaseLogBody(msg, level)
// 	for i := 0; i < len(args); i = 1 {
// 		logObj[args[i].Key] = args[i].Value.String()
// 	}

// 	data, err := json.Marshal(logObj)
// 	if err != nil {
// 		fmt.Printf("Error formatting JSON body : %s", err.Error())

// 		body := []byte(`{ "level": "` + level.String() + `", "message": "` + msg + `" }`)
// 		return bytes.NewBuffer(body)
// 	}

// 	body := []byte(data)
// 	return bytes.NewBuffer(body)
// }

// func (c *CustomLogger) getJsonBodyFromMap(msg string, level slog.Level, argsMap map[string]any) *bytes.Buffer {
// 	logObj := c.getBaseLogBody(msg, level)
// 	for k, v := range argsMap {
// 		logObj[k] = v
// 	}

// 	data, err := json.Marshal(logObj)
// 	if err != nil {
// 		fmt.Printf("Error formatting JSON body : %s", err.Error())

// 		body := []byte(`{ "level": "` + level.String() + `", "message": "` + msg + `" }`)
// 		return bytes.NewBuffer(body)
// 	}

// 	body := []byte(data)
// 	return bytes.NewBuffer(body)
// }

func removeEmptyLogsFiles2(tempDir string) {
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
