package commands

import (
	"log/slog"
	"os"
	wlptypes "welovepdf/pkg/types"
	utils "welovepdf/pkg/utils"
)

func BuildCompressSinglePageFile(
	convertPdfToJpeg func(targetImageQuality int, config *wlptypes.FileToFileOperationConfig) error,
	convertJpegToPdf wlptypes.FileToFileOperation,
	tempDirPath string,
) func(targetImageQuality int, compressionConfig *wlptypes.FileToFileOperationConfig) error {
	return func(
		targetImageQuality int,
		compressionConfig *wlptypes.FileToFileOperationConfig) error {
		tempFilePath := utils.GetNewTempFilePath(tempDirPath, "jpg")
		defer os.Remove(tempFilePath)

		err := convertPdfToJpeg(targetImageQuality, &wlptypes.FileToFileOperationConfig{
			SourceFilePath: compressionConfig.SourceFilePath,
			TargetFilePath: tempFilePath,
		})

		if err != nil {
			slog.Error("CompressSinglePageFile : Error converting file to JPG: %s", slog.String("filepath", tempFilePath))
			return err
		}

		return convertJpegToPdf(&wlptypes.FileToFileOperationConfig{
			TargetFilePath: compressionConfig.TargetFilePath,
			SourceFilePath: tempFilePath,
		})
	}
}
