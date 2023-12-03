package utils

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func getFileExtensionFromPath(inputFilePath string) string {
	pathParts := strings.Split(inputFilePath, ".")
	ext := pathParts[len(pathParts)-1]
	if len(ext) < 2 || len(ext) > 3 {
		return ""
	}
	return ext
}

func EnsureDirectory(dirPath string) error {
	stats, err := os.Stat(dirPath)
	if err == nil && stats.IsDir() {
		return nil
	}

	if !os.IsNotExist((err)) {
		return err
	}

	creationErr := os.MkdirAll(dirPath, 0755)

	if creationErr != nil {
		return err
	}

	return nil
}

func WriteContentToFileIfNotExists(filePath string, content []byte) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Chmod(0755)
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	return err
}

func GetFileNameFromPath(inputFilePath string) string {
	pathParts := strings.Split(inputFilePath, "/")
	return pathParts[len(pathParts)-1]
}

func GetFileNameWoExtensionFromPath(inputFilePath string) string {
	pathParts := strings.Split(GetFileNameFromPath(inputFilePath), ".")
	return strings.Join(pathParts[:len(pathParts)-1], ".")
}

func GetNewTempFilePath(tempDirPath string, extension string) string {
	newId := uuid.New().String()
	return path.Join(tempDirPath, newId+"."+extension)
}

func AddSuffixToFileNameInPath(fileName string, suffix string) string {
	nameParts := strings.Split(GetFileNameFromPath(fileName), ".")
	if len(nameParts) < 2 {
		return fileName + suffix
	}
	nameParts[len(nameParts)-2] = nameParts[len(nameParts)-2] + suffix
	return strings.Join(nameParts, ".")
}

func getCurrentDateStr() string {
	currentTime := time.Now()
	dateStr := strings.Split(currentTime.String(), " ")[0]
	return strings.Join(strings.Split(dateStr, "-"), "")
}

func getTodaysOutputDir(userHomeDir string) string {
	return path.Join(userHomeDir, "Documents", "welovepdf", getCurrentDateStr())
}

type SearchFileConfig struct {
	RootDirPath        string
	Filename           string
	FileSize           int
	FileLastModifiedAt int
	LastModifiedStr    string
}

func isSameFile(dirEntry os.DirEntry, searchConfig *SearchFileConfig) bool {
	if dirEntry.IsDir() || strings.HasPrefix(dirEntry.Name(), ".") {
		return false
	}
	if dirEntry.Name() == searchConfig.Filename {
		fileInfo, fileInfoErr := dirEntry.Info()
		if fileInfoErr != nil {
			slog.Debug(fmt.Sprintf("Error reading file w/ matching name : %s", dirEntry.Name()), slog.String("reason", fileInfoErr.Error()))
			return false
		}
		modifiedStr := fmt.Sprintf("%d", fileInfo.ModTime().Unix())
		slog.Debug("Found file w/ matching name", slog.String("searched", searchConfig.Filename), slog.String("actual", fileInfo.Name()), slog.Int("size", int(fileInfo.Size())), slog.String("modified", modifiedStr))

		if fileInfo.Size() == int64(searchConfig.FileSize) &&
			strings.Compare(modifiedStr, searchConfig.LastModifiedStr[:len(modifiedStr)]) == 0 {
			return true
		}
		return false
	}

	return false
}

func SearchFileInDirectoryTree(config *SearchFileConfig) string {
	config.LastModifiedStr = fmt.Sprintf("%d", config.FileLastModifiedAt)
	slog.Debug(fmt.Sprintf("config.LastModifiedStr = %s", config.LastModifiedStr))
	matchingFilePath := ""

	slog.Debug(fmt.Sprintf("Looking for file %s in dir %s", config.Filename, config.RootDirPath), slog.String("searched", config.Filename), slog.Int("size", int(config.FileSize)), slog.String("modified", config.LastModifiedStr))

	err := filepath.WalkDir(config.RootDirPath, func(currentPath string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			slog.Debug("error walking through directory", slog.String("reason", err.Error()))
			return nil
		}
		if dirEntry.IsDir() || strings.HasPrefix(dirEntry.Name(), ".") {
			return nil
		}

		shouldReturnValue := isSameFile(dirEntry, config)
		if !shouldReturnValue {
			return nil
		}

		matchingFilePath = currentPath
		return io.EOF
	})

	if err == io.EOF {
		return matchingFilePath
	}

	if err != nil {
		slog.Debug(fmt.Sprintf("Error looking for file %s in dir %s", config.Filename, config.RootDirPath), slog.String("reason", err.Error()))
		return ""
	}

	if matchingFilePath == "" {
		slog.Debug(fmt.Sprintf("No file matching %s in dir %s", config.Filename, config.RootDirPath))
		return ""
	}

	slog.Debug("found matching file", slog.String("filepath", matchingFilePath))
	return matchingFilePath
}

func removeBadCharacters(input string, dictionary []string) string {

	temp := input

	for _, badChar := range dictionary {
		temp = strings.Replace(temp, badChar, "", -1)
	}
	return temp
}

func SanitizeFilePath(path string) string {
	var badCharacters = []string{
		"../",
		"<!--",
		"-->",
		"<",
		">",
		"'",
		"\"",
		"&",
		"$",
		"#",
		"{", "}", "[", "]", "=",
		";", "?", "%20", "%22",
		"%3c",   // <
		"%253c", // <
		"%3e",   // >
		"",      // > -- fill in with % 0 e - without spaces in between
		"%28",   // (
		"%29",   // )
		"%2528", // (
		"%26",   // &
		"%24",   // $
		"%3f",   // ?
		"%3b",   // ;
		"%3d",   // =
	}

	if path == "" {
		return path
	}

	// trim(remove)white space
	trimmed := strings.TrimSpace(path)

	// trim(remove) white space in between characters
	trimmed = strings.Replace(trimmed, " ", "", -1)

	// remove bad characters from filename
	trimmed = removeBadCharacters(trimmed, badCharacters)

	stripped := strings.Replace(trimmed, "\\", "", -1)

	return stripped
}

func ComputeTargetFilePath(outputDir string, originalPath string, extension string, suffix string) string {
	fileName := GetFileNameWoExtensionFromPath(originalPath)
	formattedFileName := SanitizeFilePath(AddSuffixToFileNameInPath(fileName+"."+extension, suffix))

	for ctr := 0; ctr < 1000; ctr += 1 {
		curFileName := formattedFileName
		if ctr > 0 {
			curFileName = AddSuffixToFileNameInPath(curFileName, "("+fmt.Sprintf("%d", ctr)+")")
		}

		_, err := os.Stat(curFileName)
		if err != nil && os.IsNotExist(err) {
			return path.Join(outputDir, formattedFileName)
		}
	}
	return ""
}
