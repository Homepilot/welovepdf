package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/elastic/go-sysinfo"
	"github.com/google/uuid"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TODO : https://go.dev/play/p/0yJNk065ftB

type CustomLogger struct {
	logger       *slog.Logger
	lumberjack   *lumberjack.Logger
	logtailToken string
}

func (c *CustomLogger) Close() {
	c.lumberjack.Close()
}

func (c *CustomLogger) sendLogToLogtail(msg string, args []slog.Attr) {
	logTailUrl := "https://in.logs.betterstack.com/"
	bodyStr := getJsonBodyFromArgs(msg, args)
	body := []byte(bodyStr)
	fmt.Printf("Sending to LOGTAIL w/ body : %s", bodyStr)
	req, err := http.NewRequest(http.MethodPost, logTailUrl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Sending to LOGTAIL failed at creating req : %s", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.logtailToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Sending to LOGTAIL failed at sending res: %s", err.Error())
		return
	}
	fmt.Printf("LOGTAIL req status : %s", res.Status)
	defer res.Body.Close()
	fmt.Println("Sending to LOGTAIL done")
}

func (c *CustomLogger) Info(msg string, args ...slog.Attr) {
	c.logger.Info(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(msg, args)
}

func (c *CustomLogger) Error(msg string, args ...slog.Attr) {
	c.logger.Error(msg, args)
	if c.logtailToken == "" {
		return
	}
	c.sendLogToLogtail(msg, args)
}

func enrichLoggerWithSysInfo(logger *slog.Logger) *slog.Logger {
	goInfo := sysinfo.Go()
	host, _ := sysinfo.Host()
	hostInfo := host.Info()
	osInfo := hostInfo.OS

	return logger.With(
		slog.Group("sysinfo",
			slog.String("OS", goInfo.OS),
			slog.String("Arch", goInfo.Arch),
			slog.Int("MaxProcs", goInfo.MaxProcs),
			slog.String("Version", goInfo.Version),
		),
		slog.Group("hostinfo",
			slog.String("Architecture", hostInfo.Architecture),
			slog.String("Hostname", hostInfo.Hostname),
			slog.String("KernelVersion", hostInfo.KernelVersion),
			slog.String("UniqueID", hostInfo.UniqueID),
		),
		slog.Group("osinfo",
			slog.String("Type", osInfo.Type),
			slog.String("Family", osInfo.Family),
			slog.String("Platform", osInfo.Platform),
			slog.String("Name", osInfo.Name),
			slog.String("Version", osInfo.Version),
			slog.Int("Major", osInfo.Major),
			slog.Int("Minor", osInfo.Minor),
			slog.Int("Patch", osInfo.Patch),
			slog.String("Build", osInfo.Build),
		),
	)
}

func getJsonBodyFromArgs(msg string, args []slog.Attr) string {
	jsonStr := `{ "message": "` + msg + `"`
	for i := 0; i < len(args); i = 1 {
		jsonStr += `, "` + args[i].Key + `": "` + args[i].Value.Kind().String() + `"`
		// fmt.Printf("field: %s", args[i])
	}
	return jsonStr + ` }`
}

func InitLogger(tempDir string, logtailToken string) *CustomLogger {
	logsDir := path.Join(tempDir, "..", "logs")
	EnsureDirectory(logsDir)
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
	enrichedLogger := enrichLoggerWithSysInfo(logger)
	slog.SetDefault(enrichedLogger)

	return &CustomLogger{
		logger:       enrichedLogger,
		lumberjack:   lj,
		logtailToken: logtailToken,
	}
}

func RemoveEmptyLogsFiles(tempDir string) {
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
