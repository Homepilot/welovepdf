package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileExtensionFromPath(t *testing.T) {
	assert.Equal(t, "txt", getFileExtensionFromPath("example.txt"))
	assert.Equal(t, "jpg", getFileExtensionFromPath("my.folder/example.jpg"))
	assert.Equal(t, "", getFileExtensionFromPath("example"))
}

func TestEnsureDirectory(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), "testdir")
	defer os.RemoveAll(tempDirPath)

	err := EnsureDirectory(tempDirPath)
	assert.NoError(t, err)

	dir, err := os.Stat(tempDirPath)
	assert.NoError(t, err)
	assert.True(t, dir.IsDir())
}

func TestGetFileNameFromPath(t *testing.T) {
	assert.Equal(t, "example.txt", GetFileNameFromPath("/tmp/example.txt"))
	assert.Equal(t, "example.txt", GetFileNameFromPath("example.txt"))
}

func TestGetFileNameWoExtensionFromPath(t *testing.T) {
	assert.Equal(t, "example", GetFileNameWoExtensionFromPath("/tmp/example.txt"))
	assert.Equal(t, "example", GetFileNameWoExtensionFromPath("example.txt"))
}

func TestGetNewTempFilePath(t *testing.T) {
	tempDir := os.TempDir()
	filePath := GetNewTempFilePath(tempDir, "txt")
	assert.Contains(t, filePath, tempDir)
	assert.Contains(t, filePath, ".txt")
}

func TestAddSuffixToFileName(t *testing.T) {
	assert.Equal(t, "example_suffix.txt", AddSuffixToFileName("example.txt", "_suffix"))
	assert.Equal(t, "example_suffix", AddSuffixToFileName("example", "_suffix"))
}

func TestGetTodaysOutputDir(t *testing.T) {
	userHomeDir, _ := os.UserHomeDir()
	outputDir := GetTodaysOutputDir(userHomeDir)
	assert.Contains(t, outputDir, userHomeDir)
	assert.Contains(t, outputDir, "Documents/welovepdf")
}
