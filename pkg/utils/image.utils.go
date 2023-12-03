package utils

import (
	"bufio"
	"image"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

type Format struct {
	Height int
	Width  int
}

var A4_HEIGHT int = 100
var A4_WIDTH int = 50

func ConvertImageToJpeg(sourceFilePath string, targetFilePath string) error {
	var img image.Image

	slog.Debug("ConvertImageToJpeg : Operation started", slog.String("source", sourceFilePath), slog.String("target", targetFilePath))
	sourceFileExt := strings.ToLower(filepath.Ext(sourceFilePath))
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		slog.Debug("Error in ConvertImageToJpeg 1", slog.String("reason", err.Error()))
		return err
	}
	defer sourceFile.Close()
	r := bufio.NewReader(sourceFile)
	slog.Debug("Reader opened, starting to decode")

	if sourceFileExt == ".jpeg" || sourceFileExt == ".jpg" {
		err = os.Rename(sourceFilePath, targetFilePath)
		img = nil
	}
	if sourceFileExt == ".png" {
		img, err = png.Decode(r)
	}
	if sourceFileExt == ".webp" {
		img, err = webp.Decode(r)
	}
	if sourceFileExt == ".tiff" {
		img, err = tiff.Decode(r)
	}
	if err != nil || img == nil {
		slog.Debug("Error in ConvertImageToJpeg 2", slog.String("reason", err.Error()))
		return err
	}

	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		slog.Debug("Error in ConvertImageToJpeg 3", slog.String("reason", err.Error()))
		return err
	}
	defer targetFile.Close()
	targetFile.Chmod(0755)

	w := bufio.NewWriter(targetFile)
	return jpeg.Encode(w, img, &jpeg.Options{})
}
