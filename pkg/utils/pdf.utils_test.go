package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompressAllFilesInDir tests the CompressAllFilesInDir function.
// This test assumes that the directory and files for compression exist.
// You might want to set up a test directory with files before running this test.
func TestCompressAllFilesInDir(t *testing.T) {
	sourceDirPath := ("../../test/assets/pdf")
	tempDirPath := filepath.Join(os.TempDir(), "test_compress")
	targetDirPath := filepath.Join(tempDirPath, "target")
	os.MkdirAll(targetDirPath, 0755)
	defer os.RemoveAll(targetDirPath)

	config := &DirToDirOperationConfig{
		SourceDirPath: sourceDirPath,
		TargetDirPath: targetDirPath,
		BinaryPath:    "test",
	}

	err := CompressAllFilesInDir(tempDirPath, 75, "toto.ps", config)
	assert.Nil(t, err)
}

// TestMergeAllFilesInDir tests the MergeAllFilesInDir function.
// This test assumes that the directory for merging exists and contains PDF files.
// func TestMergeAllFilesInDir(t *testing.T) {
// 	sourceDirPath := "../../test/assets/pdf"
// 	targetFilePath := filepath.Join(os.TempDir(), "merged.pdf")
// 	os.MkdirAll(sourceDirPath, 0755)
// 	defer os.Remove(targetFilePath)

// 	success := MergeAllFilesInDir(sourceDirPath, targetFilePath)
// 	assert.True(t, success)
// 	stats, err := os.Stat(targetFilePath)
// 	assert.Nil(t, err)
// 	assert.False(t, stats.IsDir())
// }

// TestConvertImageToPdf tests the ConvertImageToPdf function.
// This test assumes that the source image file exists.
// func TestConvertImageToPdf(t *testing.T) {
// 	sourceFilePath := "../../test/assets/img/test-image-1.png" // Replace with actual test image path
// 	targetFilePath := filepath.Join(os.TempDir(), "converted.pdf")
// 	defer os.Remove(targetFilePath)

// 	success := ConvertImageToPdf(sourceFilePath, targetFilePath)
// 	assert.True(t, success)
// 	stats, err := os.Stat(targetFilePath)
// 	assert.Nil(t, err)
// 	assert.False(t, stats.IsDir())
// }

// TestCompressSinglePageFile tests the CompressSinglePageFile function.
// This test assumes that the source file for compression exists.
// func TestCompressSinglePageFile(t *testing.T) {
// 	tempDirPath := os.TempDir()
// 	sourceFilePath := "../../test/assets/pdf/test-document-1.pdf" // Replace with actual test file path
// 	targetFilePath := filepath.Join(tempDirPath, "compressed.pdf")
// 	defer os.Remove(targetFilePath)

// 	config := &FileToFileOperationConfig{
// 		SourceFilePath: sourceFilePath,
// 		TargetFilePath: targetFilePath,
// 	}

// 	success := CompressSinglePageFile(tempDirPath, 20, config)
// 	assert.True(t, success)
// }

// TestMergePdfFiles tests the MergePdfFiles function.
// This test assumes that the PDF files to be merged exist.
// func TestMergePdfFiles(t *testing.T) {
// 	targetFilePath := filepath.Join(os.TempDir(), "merged.pdf")
// 	filePaths := []string{"../../test/assets/pdf/test-document-1.pdf", "../../test/assets/pdf/test-document-2.pdf"} // Replace with actual file paths

// 	success := MergePdfFiles(targetFilePath, filePaths)
// 	assert.True(t, success)
// 	stats, err := os.Stat(targetFilePath)
// 	assert.Nil(t, err)
// 	assert.False(t, stats.IsDir())
// }
