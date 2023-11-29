package utils

import (
	"bufio"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

func convertImageToJpeg(sourceFilePath string, targetFilePath string) error {
	var img image.Image

	sourceFileExt := strings.ToLower(filepath.Ext(sourceFilePath))
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	r := bufio.NewReader(sourceFile)

	if sourceFileExt == "png" {
		img, err = png.Decode(r)
	}
	if sourceFileExt == "webp" {
		img, err = webp.Decode(r)
	}
	if sourceFileExt == "tiff" {
		img, err = tiff.Decode(r)
	}
	if err != nil || img == nil {
		return err
	}

	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	targetFile.Chmod(0755)
	w := bufio.NewWriter(targetFile)

	return jpeg.Encode(w, img, &jpeg.Options{})
}
