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

	"github.com/disintegration/imaging"
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

	if sourceFileExt == ".jpeg" {
		img, err = jpeg.Decode(r)
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

	imgToUse := img
	isLandscapeOriented, err := isLandscape(img)
	if err != nil {
		slog.Debug("Error in ConvertImageToJpeg 3", slog.String("reason", err.Error()))
		return err

	}
	if isLandscapeOriented {
		imgToUse = imaging.Rotate90(img)
	}

	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		slog.Debug("Error in ConvertImageToJpeg 3", slog.String("reason", err.Error()))
		return err
	}
	defer targetFile.Close()
	targetFile.Chmod(0755)

	w := bufio.NewWriter(targetFile)
	return jpeg.Encode(w, imgToUse, &jpeg.Options{})
}

func isLandscape(img image.Image) (bool, error) {
	srcBounds := img.Bounds()
	imgWidth := srcBounds.Dx()
	imgHeight := srcBounds.Dy()

	return imgHeight < imgWidth, nil
}
