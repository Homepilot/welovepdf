package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	pdfApi "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfUtils struct{}

func NewPdfUtils() *PdfUtils {
	return &PdfUtils{}
}

func (p *PdfUtils) MergePdfFiles(targetFileName string, filePathes []string) error {
	fmt.Println("MergePdfFiles: operation starting")
	// EnsureTargetDirPath()

	err := pdfApi.MergeCreateFile(filePathes, targetFileName+".pdf", pdfApi.LoadConfiguration())
	if err != nil {
		fmt.Printf("Error retrieving targetPath: %s", err.Error())
		return err
	}

	fmt.Println("Operation succeeded, opening target folder")

	cmd := exec.Command("open /Users/gregoire/Documents")
	openErr := cmd.Run()
	if openErr != nil {
		fmt.Printf("Error opening target folder: %s", openErr.Error())
	}
	directory, openErr2 := os.Open("/Users/gregoire/Documents")

	fmt.Printf("result is here: %s", directory.Name())

	if openErr2 != nil {
		fmt.Printf("Error opening target folder w/ Open: %s", openErr.Error())
	}

	return nil
}

func (p *PdfUtils) OptimizePdfFile(filePath string) error {
	fmt.Println("OptimizePdfFile: operation starting")
	// EnsureTargetDirPath()

	nameParts := strings.Split(GetFileNameFromPath(filePath), ".")
	nameParts[len(nameParts)-2] = nameParts[len(nameParts)-2] + "_compressed"
	targetFilePath := GetTargetDirectoryPath() + "/" + strings.Join(nameParts, ".")

	err := pdfApi.OptimizeFile(filePath, targetFilePath, pdfApi.LoadConfiguration())
	if err != nil {
		fmt.Printf("Error retrieving targetPath: %s", err.Error())
		return err
	}

	fmt.Println("Operation succeeded, opening target folder")

	return nil
}

func (p *PdfUtils) ConvertImageToPdf(filePath string) error {
	fmt.Println("ConvertImageToPdf: operation starting")

	originalFileName := GetFileNameFromPath(filePath)
	fileNameParts := strings.Split(originalFileName, ".")
	fileNameParts[len(fileNameParts)-1] = "pdf"
	targetFileName := strings.Join(fileNameParts, ".")

	conversionError := pdfApi.ImportImagesFile([]string{filePath}, GetTargetDirectoryPath()+"/"+targetFileName, nil, nil)

	if conversionError != nil {
		fmt.Printf("Error importing image: %s", conversionError.Error())
		return conversionError
	}

	fmt.Println("Operation succeeded, opening target folder")

	return nil
}

func (p *PdfUtils) CompressFile(filePath string) error {
	// fmt.Println=4 -dPDFSETTINGS=/screen -dNOPAUSE -dQUIET -dBATCH -sOutputFile=output.pdf input.pdf
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
		fmt.Printf("Error compressing file: %s", err.Error())
		return err
	}

	fmt.Printf("Success compressing file: %s", out)

	fmt.Printf("Operation succeeded, opening target folder")
	return nil
}
