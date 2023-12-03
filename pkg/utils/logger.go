package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/elastic/go-sysinfo"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

var SYS_INFO_KEY = "sysinfo"
var HOST_INFO_KEY = "hostinfo"
var OS_INFO_KEY = "osinfo"

type CustomLogger struct {
	sysInfoKey    string
	lumberjack    *lumberjack.Logger
	logtailToken  string
	logsDirPath   string
	logsToSend    []string
	logsBatchSize int
}

func SetupLogger(appConfig *AppConfig) *CustomLogger {
	removeEmptyLogsFiles(appConfig.Logger.LogsDirPath)
	fileName := strings.Join(strings.Split(time.Now().Local().Format(time.DateTime), " "), "") + ".log"
	logFilePath := path.Join(appConfig.Logger.LogsDirPath, fileName)

	lj := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}

	logTailLogLevel := appConfig.Logger.LogLevel
	if logTailLogLevel == slog.LevelDebug {
		logTailLogLevel = slog.LevelInfo
	}

	customLogger := &CustomLogger{
		lumberjack:    lj,
		logtailToken:  appConfig.Logger.LogtailToken,
		logsDirPath:   appConfig.Logger.LogsDirPath,
		logsBatchSize: appConfig.Logger.LogsBatchSize,
	}

	slogHandlers := slogmulti.Fanout(
		slog.NewTextHandler(lj, &slog.HandlerOptions{AddSource: appConfig.DebugMode, Level: appConfig.Logger.LogLevel, ReplaceAttr: removeSysInfoAttr}),
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: appConfig.DebugMode, Level: appConfig.Logger.LogLevel, ReplaceAttr: removeSysInfoAttr}),
		slog.NewJSONHandler(customLogger, &slog.HandlerOptions{Level: logTailLogLevel, ReplaceAttr: replaceLogtailAttr}),
	)
	if appConfig.Logger.LogtailToken == "" {
		slogHandlers = slogmulti.Fanout(
			slog.NewTextHandler(lj, &slog.HandlerOptions{AddSource: appConfig.DebugMode, Level: appConfig.Logger.LogLevel}),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: appConfig.DebugMode, Level: appConfig.Logger.LogLevel}),
		)
	}

	slogLogger := slog.New(slogHandlers)
	enrichedLogger := addSysInfoToLogger(slogLogger)
	slog.SetDefault(enrichedLogger)

	return customLogger
}

// Should only be called by a slog json handler
func (l *CustomLogger) Write(data []byte) (int, error) {
	l.logsToSend = append(l.logsToSend, string(data))
	if len(l.logsToSend) >= l.logsBatchSize {
		err := l.flush()
		if err != nil {
			slog.Debug("Error sending batch to Logtail", slog.String("reason", err.Error()))
		}
	}

	return len(data), nil
}

func (l *CustomLogger) flush() error {
	logsBatch := l.logsToSend
	l.logsToSend = []string{}
	slog.Debug("SENDING LOGS BATCH W/ LENGTH", slog.Int("batchLength", len(logsBatch)))
	logsBatchAsJson, err := mergeLogsArrayToJson(logsBatch)
	if logsBatchAsJson == nil {
		slog.Debug("Error merging json obj", slog.String("reason", err.Error()))
		return err
	}
	err = l.sendToLogtail(logsBatchAsJson)
	if err != nil {
		fmt.Printf("Error sending to logtail : %s", err.Error())
	}
	slog.Debug("Logs batch successfully sent to Logtail")
	return err
}

func (c *CustomLogger) Close() {
	c.lumberjack.Close()
	c.flush()
	removeEmptyLogsFiles(c.logsDirPath)
}

func (c *CustomLogger) sendToLogtail(body []byte) error {
	LOGTAIL_URL := "https://in.logs.betterstack.com/"

	req, err := http.NewRequest(http.MethodPost, LOGTAIL_URL, bytes.NewBuffer(body))
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

func replaceLogtailAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.String("dt", a.Value.Time().Format(time.RFC3339))
	}

	if a.Key == slog.MessageKey {
		return slog.String("message", a.Value.String())
	}

	return a
}

func removeSysInfoAttr(groups []string, a slog.Attr) slog.Attr {
	keysToRemove := map[string]bool{
		SYS_INFO_KEY:  true,
		HOST_INFO_KEY: true,
		OS_INFO_KEY:   true,
	}
	parentGroup := strings.Split(a.Key, ".")[0]
	if keysToRemove[parentGroup] {
		return slog.Attr{}
	}

	return a
}

func addSysInfoToLogger(logger *slog.Logger) *slog.Logger {
	goInfo := sysinfo.Go()
	host, _ := sysinfo.Host()
	hostInfo := host.Info()
	osInfo := hostInfo.OS

	sysInfoAttr := slog.Group(
		SYS_INFO_KEY,
		slog.String("OS", goInfo.OS),
		slog.String("Arch", goInfo.Arch),
		slog.Int("MaxProcs", goInfo.MaxProcs),
		slog.String("Version", goInfo.Version),
	)
	hostInfoAttr := slog.Group(
		HOST_INFO_KEY,
		slog.String("Architecture", hostInfo.Architecture),
		slog.String("Hostname", hostInfo.Hostname),
		slog.String("KernelVersion", hostInfo.KernelVersion),
		slog.String("UniqueID", hostInfo.UniqueID),
	)
	osInfoAttr := slog.Group(
		OS_INFO_KEY,
		slog.String("Type", osInfo.Type),
		slog.String("Family", osInfo.Family),
		slog.String("Platform", osInfo.Platform),
		slog.String("Name", osInfo.Name),
		slog.String("Version", osInfo.Version),
		slog.Int("Major", osInfo.Major),
		slog.Int("Minor", osInfo.Minor),
		slog.Int("Patch", osInfo.Patch),
		slog.String("Build", osInfo.Build),
	)

	return logger.With(sysInfoAttr, hostInfoAttr, osInfoAttr)
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

func mergeLogsArrayToJson(jsonStrArr []string) ([]byte, error) {
	var jsonObjArr []map[string]any
	for _, jsonObjStr := range jsonStrArr {
		jsonObj := []byte(jsonObjStr)
		var dataMap map[string]any
		err := json.Unmarshal(jsonObj, &dataMap)
		if err != nil {
			slog.Debug("Error decoding json obj", slog.String("reason", err.Error()), slog.String("json", jsonObjStr))
			continue
		}
		jsonObjArr = append(jsonObjArr, dataMap)
	}
	mergedJson, err := json.Marshal(jsonObjArr)
	if err != nil {
		return nil, err
	}

	return mergedJson, nil
}
