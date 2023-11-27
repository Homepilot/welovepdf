package utils

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/elastic/go-sysinfo"
	"github.com/google/uuid"
)

// TODO : https://go.dev/play/p/0yJNk065ftB

func getLogWriters(tempDir string) io.Writer {
	logsDir := path.Join(tempDir, "..", "logs")
	EnsureDirectory(logsDir)
	logFilePath := path.Join(logsDir, uuid.NewString()+".log")
	tempFile, _ := os.Create(logFilePath)
	defer tempFile.Close()
	tempFile.Chmod(0755)
	logFileWriter := bufio.NewWriter(tempFile)

	// return io.MultiWriter(logFileWriter, os.Stdout)
	return logFileWriter
}

func InitLogger(tempDir string) *slog.Logger {
	handler := slog.NewJSONHandler(getLogWriters(tempDir), nil)
	logger := slog.New(handler)

	goInfo := sysinfo.Go()
	host, _ := sysinfo.Host()
	hostInfo := host.Info()
	osInfo := hostInfo.OS

	enrichedLogger := logger.With(
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

	slog.SetDefault(enrichedLogger)

	return enrichedLogger
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