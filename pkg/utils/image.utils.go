package utils

import (
	"bufio"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
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

func resizeJpegToFitA4(sourceFilePath string, targetFilePath string, tempDirPath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	src, err := jpeg.Decode(sourceFile)
	srcBounds := src.Bounds()
	imgWidth := srcBounds.Dx()
	imgHeight := srcBounds.Dy()

	needResize := imgHeight > A4_HEIGHT || imgWidth <= A4_WIDTH
	if !needResize {
		return os.Rename(sourceFilePath, targetFilePath)
	}

	isLandscape := imgHeight < imgWidth
	// Check if should rotate 90°
	if isLandscape {
		tempFilePath := GetNewTempFilePath(tempDirPath, "jpg")
		_ = sourceFile.Close()
		err = RotateImageClockwise90(sourceFilePath, tempFilePath)
		if err != nil {
			return err
		}
		return resizeJpegToFitA4(tempFilePath, targetFilePath, tempDirPath)
	}

	outputFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	// Set the target size
	targetHeight, targetWidth := getResizeToFormatDimensions(imgHeight, imgWidth, &Format{
		Height: A4_HEIGHT,
		Width:  A4_WIDTH,
	})
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Resize
	draw.ApproxBiLinear.Scale(dst, dst.Rect, src, srcBounds, draw.Over, nil)
	return jpeg.Encode(outputFile, dst, &jpeg.Options{})
}

func getResizeToFormatDimensions(imgHeight int, imgWidth int, format *Format) (int, int) {
	var targetHeight, targetWidth float64
	srcHeightFloat := float64(imgHeight)
	srcWidthFloat := float64(imgWidth)
	// Determine scale factor from dimensions
	heightScaleFactor := float64(format.Height) / srcHeightFloat
	widthScaleFactor := float64(format.Width) / srcWidthFloat

	if heightScaleFactor >= 1 && widthScaleFactor >= 1 {
		return imgHeight, imgWidth
	}

	if widthScaleFactor > heightScaleFactor {
		targetHeight = widthScaleFactor * srcHeightFloat
		targetWidth = widthScaleFactor * srcWidthFloat
	} else {
		targetHeight = heightScaleFactor * srcHeightFloat
		targetWidth = heightScaleFactor * srcWidthFloat

	}
	return int(math.Round(targetHeight)), int(math.Round(targetWidth))
}

func RotateImageClockwise90(sourceFilePath string, targetFilePath string) error {
	slog.Debug("RotateImageClockwise90 : Operation started")
	slog.Debug("RotateImageClockwise90")
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	srcImage, _, err := image.Decode(sourceFile)
	slog.Debug("RotateImageClockwise90 2")
	if err != nil {
		return err
	}
	slog.Debug("RotateImageClockwise90 3")
	rotatedImage := tensorToImage(rotateMatrixBy90Degrees[color.Color](imageToTensor(srcImage)))
	slog.Debug("RotateImageClockwise90 3")

	targetDim := rotatedImage.Bounds()
	dstImage := image.NewRGBA(image.Rect(0, 0, targetDim.Dx(), targetDim.Dy()))

	slog.Debug("RotateImageClockwise90 4")
	newImage, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer newImage.Close()
	slog.Debug("RotateImageClockwise90 5")
	return jpeg.Encode(newImage, dstImage, &jpeg.Options{})
}

func imageToTensor(img image.Image) *[][]color.Color {
	slog.Debug("imageToTensor")
	size := img.Bounds().Size()
	var pixels [][]color.Color
	// put pixels into two three two dimensional array
	for i := 0; i < size.X; i++ {
		var y []color.Color
		for j := 0; j < size.Y; j++ {
			y = append(y, img.At(i, j))
		}
		pixels = append(pixels, y)
	}

	return &pixels
}

func tensorToImage(pixels *[][]color.Color) image.Image {
	slog.Debug("tensorToImage")
	ppixels := *pixels
	rect := image.Rect(0, 0, len(ppixels), len(ppixels[0]))
	nImg := image.NewRGBA(rect)

	for x := 0; x < len(ppixels); x++ {
		for y := 0; y < len(ppixels[0]); y++ {
			q := ppixels[x]
			if q == nil {
				continue
			}
			p := ppixels[x][y]
			if p == nil {
				continue
			}
			original, ok := color.RGBAModel.Convert(p).(color.RGBA)
			if ok {
				nImg.Set(x, y, original)
			}
		}
	}

	return rect
}

func rotateMatrixBy90Degrees[T any](matrix *[][]T) *[][]T {
	slog.Debug("rotateMatrixBy90Degrees")
	pmatrix := *matrix
	newMatrix := [][]T{}
	// To rotate by 90°, simply set newX = oldY && newY = oldX
	for i := 0; i < len(pmatrix); i++ {
		newArr := []T{}
		for j := 0; j < len(pmatrix[i]); j++ {
			newArr = append(newArr, pmatrix[i][j])
		}
		newMatrix = append(newMatrix, newArr)
	}
	return &newMatrix
}
