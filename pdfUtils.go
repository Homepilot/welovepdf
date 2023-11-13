package main

import (
	"os"
	"os/exec"
	"strings"

	pdfApi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) MergePdfFiles(targetFileName string, filePathes []string) bool {
	runtime.LogInfo(a.ctx, "MergePdfFiles: operation starting")
	a.EnsureTargetDirPath()

	err := pdfApi.MergeCreateFile(filePathes, targetFileName+".pdf", pdfApi.LoadConfiguration())
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", err.Error())
		return false
	}

	runtime.LogInfo(a.ctx, "Operation succeeded, opening target folder")
	cmd := exec.Command("open /Users/gregoire/Documents")
	openErr := cmd.Run()
	if openErr != nil {
		runtime.LogErrorf(a.ctx, "Error opening target folder: %s", openErr.Error())
	}
	directory, openErr2 := os.Open("/Users/gregoire/Documents")

	runtime.LogInfof(a.ctx, "result is here: %s", directory.Name())

	if openErr2 != nil {
		runtime.LogErrorf(a.ctx, "Error opening target folder w/ Open: %s", openErr.Error())
	}
	return true
}

func (a *App) CompressPdfFile(filePath string) bool {
	runtime.LogInfo(a.ctx, "CompressPdfFile: operation starting")
	a.EnsureTargetDirPath()

	pathParts := strings.Split(filePath, ".")
	pathParts[len(pathParts)-2] = pathParts[len(pathParts)-2] + "_compressed"
	targetFilePath := strings.Join(pathParts, ".")

	err := pdfApi.OptimizeFile(filePath, targetFilePath, pdfApi.LoadConfiguration())
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", err.Error())
		return false
	}

	runtime.LogInfo(a.ctx, "Operation succeeded, opening target folder")
	os.Open(a.targetDirPath)
	return true
}

func (a *App) ConvertImageToPdf(filePath string) bool {
	runtime.LogInfo(a.ctx, "ConvertImageToPdf: operation starting")
	a.EnsureTargetDirPath()

	originalFileName := GetFileNameFromPath(filePath)
	fileNameParts := strings.Split(originalFileName, ".")
	fileNameParts[len(fileNameParts)-1] = "pdf"
	targetFileName := strings.Join(fileNameParts, ".")

	conversionError := pdfApi.ImportImagesFile([]string{filePath}, a.targetDirPath+"/"+targetFileName, nil, nil)

	if conversionError != nil {
		runtime.LogErrorf(a.ctx, "Error importing image: %s", conversionError.Error())
		return false
	}

	runtime.LogInfo(a.ctx, "Operation succeeded, opening target folder")
	os.Open(a.targetDirPath)
	return true
}

func (a *App) CompressFile(filePath string) bool {
	// gs -sDEVICE=pdfwrite -dCompatibilityLevel=1.4 -dPDFSETTINGS=/screen -dNOPAUSE -dQUIET -dBATCH -sOutputFile=output.pdf input.pdf
	cmd := exec.Command("gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dSubsetFonts=true",
		"-dUseFlateCompression=true",
		"-dOptimize=true",
		"-dProcessColorModel=/DeviceRGB",
		"-dDownsampleGrayImages=true ",
		"-dGrayImageDownsampleType=/Bicubic",
		"-dGrayImageResolution=75",
		"-dAutoFilterGrayImages=false",
		"-dDownsampleMonoImages=true",
		"-dMonoImageDownsampleType=/Bicubic",
		"-dCompressPages=true",
		"-dMonoImageResolution=75",
		"-dDownsampleColorImages=true ",
		"-dCompressStreams=true ",
		"-dColorImageDownsampleType=/Bicubic",
		"-dColorImageResolution=75",
		"-dImageQuality=80",
		"-dImageUpperPPI=100",
		"-dAutoFilterColorImages=false",
		"-dPDFSETTINGS=/default",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-dSAFER",
		"-sOutputFile="+filePath+"_compressed.pdf",
		"-dCompressFonts=true",
		"-r150",
		filePath,
	)
	// cmd := exec.Command("ps2pdf",
	// 	"-dPDFSETTINGS=/screen",
	// 	filePath,
	// 	filePath+"_compressed.pdf",
	// )

	out, err := cmd.Output()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error compressing file: %s", err.Error())
		return false
	}

	runtime.LogInfof(a.ctx, "Success compressing file: %s", out)

	runtime.LogInfo(a.ctx, "Operation succeeded, opening target folder")
	os.Open(a.targetDirPath)
	return true
}

func GetFileNameFromPath(filePath string) string {
	pathParts := strings.Split(filePath, "/")
	return pathParts[len(pathParts)-1]
}
