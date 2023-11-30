package utils

import (
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
)

type LoggerConfig struct {
	LogsDirPath   string
	LogLevel      slog.Level
	LogtailToken  string
	LogsBatchSize int
}

type AppConfig struct {
	DebugMode          bool
	Logger             *LoggerConfig
	OutputDirPath      string
	LocalAssetsDirPath string
	TempDirPath        string
	UserHomeDir        string
}

func getDefaultConfig() *AppConfig {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Error retrieving the user's home directory : %s", err.Error()))
	}
	localAssetsDirPath := path.Join(userHomeDir, ".welovepdf")

	return &AppConfig{
		OutputDirPath:      getTodaysOutputDir(userHomeDir),
		LocalAssetsDirPath: localAssetsDirPath,
		TempDirPath:        path.Join(localAssetsDirPath, "temp"),
		UserHomeDir:        userHomeDir,
		DebugMode:          false,
		Logger: &LoggerConfig{
			LogLevel:      slog.LevelInfo,
			LogtailToken:  "",
			LogsBatchSize: 10,
			LogsDirPath:   path.Join(localAssetsDirPath, "logs"),
		},
	}
}

func GetAppConfigFromAssetsDir(assetsDir embed.FS) *AppConfig {
	var configFilePath = "assets/config/config.json"
	newConfig := getDefaultConfig()

	jsonObj, err := assetsDir.ReadFile(configFilePath)
	if err != nil {
		slog.Warn("Error reading config file, returining default config")
		return newConfig
	}
	jsonValue := map[string]any{}

	err = json.Unmarshal(jsonObj, &jsonValue)
	if err != nil {
		slog.Warn("Error reading config from file content")
		return newConfig
	}

	debugMode, debugModeOk := jsonValue["debugMode"].(bool)
	if debugModeOk {
		newConfig.DebugMode = debugMode
	}
	logsBatchSize, logsBatchSizeOk := jsonValue["logsBatchSize"].(int)
	if logsBatchSizeOk {
		newConfig.Logger.LogsBatchSize = logsBatchSize
	}
	logLevel, logLevelOk := jsonValue["logLevel"].(slog.Level)
	if logLevelOk {
		newConfig.Logger.LogLevel = logLevel
	}
	logtailToken, logtailTokenOk := jsonValue["logtailToken"].(string)
	if logtailTokenOk {
		newConfig.Logger.LogtailToken = logtailToken
	}

	return newConfig
}
