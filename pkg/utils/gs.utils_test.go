package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnsureGhostScriptSetup tests the EnsureGhostScriptSetup function.
// This test will create a temporary file to simulate the GhostScript binary.
func TestEnsureGhostScriptSetup(t *testing.T) {
	dir, err := os.UserHomeDir()
	assert.Nil(t, err)
	gsBinaryPath := filepath.Join(dir, "Documents/welovepdf/ghostscript")
	os.Remove(gsBinaryPath)
	// defer os.Remove(gsBinaryPath)

	// Already setup
	_, err = os.Create(gsBinaryPath)
	assert.Nil(t, err)

	dummyContent := []byte("dummy content")
	EnsureGhostScriptSetup(gsBinaryPath, dummyContent)
	assert.True(t, IsGhostScriptSetup(gsBinaryPath))

	// Not yet setup
	err = os.Remove(gsBinaryPath)
	assert.Nil(t, err)
	EnsureGhostScriptSetup(gsBinaryPath, dummyContent)

	assert.True(t, IsGhostScriptSetup(gsBinaryPath))
	content, err := os.ReadFile(gsBinaryPath)
	assert.Nil(t, err)
	assert.Equal(t, dummyContent, content)
}

// TestIsGhostScriptSetup tests the IsGhostScriptSetup function.
func TestIsGhostScriptSetup(t *testing.T) {
	gsBinaryPath := filepath.Join(os.TempDir(), "ghostscript")
	os.Create(gsBinaryPath)
	defer os.Remove(gsBinaryPath)

	assert.True(t, IsGhostScriptSetup(gsBinaryPath))
	assert.False(t, IsGhostScriptSetup("non_existent_path"))
}

// // TestConvertToLowQualityJpeg tests the convertToLowQualityJpeg function.
// // This test assumes that the necessary binary and source file exist.
// func TestConvertToLowQualityJpeg(t *testing.T) {
// 	config := &FileToFileOperationConfig{
// 		BinaryPath:     "path_to_ghostscript_binary", // Replace with actual path
// 		SourceFilePath: "path_to_source_file",        // Replace with actual path
// 		TargetFilePath: filepath.Join(os.TempDir(), "output.jpg"),
// 	}
// 	success := convertToLowQualityJpeg(75, config)
// 	assert.True(t, success)
// }

// // TestResizePdfToA4 tests the ResizePdfToA4 function.
// // This test assumes that the necessary binary and source file exist.
// func TestResizePdfToA4(t *testing.T) {
// 	config := &FileToFileOperationConfig{
// 		BinaryPath:     "path_to_ghostscript_binary", // Replace with actual path
// 		SourceFilePath: "path_to_source_file",        // Replace with actual path
// 		TargetFilePath: filepath.Join(os.TempDir(), "resized.pdf"),
// 	}
// 	success := ResizePdfToA4(config)
// 	assert.True(t, success)
// }
