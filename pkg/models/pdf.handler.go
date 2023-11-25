package models

import (
	"log"
	"os"
	"path"

	"welovepdf/pkg/utils"

	"github.com/google/uuid"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfHandler struct {
	outputDir string
	TempDir   string
}

func NewPdfHandler(
	outputDir string,
	tempDir string,
) *PdfHandler {
	return &PdfHandler{
		outputDir: outputDir,
		TempDir:   tempDir,
	}
}

func (p *PdfHandler) MergePdfFiles(targetFilePath string, filePathes []string) bool {
	return utils.MergePdfFiles(targetFilePath, filePathes)
}

func (p *PdfHandler) OptimizePdfFile(filePath string) bool {
	log.Println("OptimizePdfFile: operation starting")

	newFileName := utils.AddSuffixToFileName(utils.GetFileNameFromPath(filePath), "_compressed")
	targetFilePath := path.Join(p.outputDir, newFileName)

	err := pdfcpu.OptimizeFile(filePath, targetFilePath, pdfcpu.LoadConfiguration())
	if err != nil {
		log.Printf("Error retrieving targetPath: %s", err.Error())
		return false
	}

	log.Println("Optimization succeeded")

	return true
}

func (p *PdfHandler) CompressFile(filePath string, targetImageQuality int) bool {
	resultFilePath := path.Join(p.outputDir, utils.GetFileNameWoExtensionFromPath(filePath)+"_compressed.pdf")

	pageCount, err := pdfcpu.PageCountFile(filePath)
	if pageCount == 1 {
		return utils.CompressSinglePageFile(p.TempDir, targetImageQuality, &utils.FileToFileOperationConfig{
			SourceFilePath: filePath,
			TargetFilePath: resultFilePath,
		})
	}

	fileId := uuid.New().String()
	tempDirPath1 := path.Join(p.TempDir, fileId, "compress_jpg")
	tempDirPath2 := path.Join(p.TempDir, fileId, "compress_pdf")
	utils.EnsureDirectory(tempDirPath1)
	utils.EnsureDirectory(tempDirPath2)
	defer os.RemoveAll(tempDirPath1)
	defer os.RemoveAll(tempDirPath2)

	// 1. Split file into 1 file per page
	err = pdfcpu.SplitFile(filePath, tempDirPath1, 1, nil)
	if err != nil {
		log.Printf("Error splitting file, compression aborted, error: %s\n", err.Error())
		return false
	}

	// 2. Compress the splitted files
	isCompressionSuccess := utils.CompressAllFilesInDir(p.TempDir, targetImageQuality, &utils.DirToDirOperationConfig{
		SourceDirPath: tempDirPath1,
		TargetDirPath: tempDirPath2,
	})
	if !isCompressionSuccess {
		log.Printf("Error compressing files in dir : %s to dir %s", tempDirPath1, tempDirPath2)
		return false
	}

	// 3. Merge all compressed files back into 1
	isMergeSuccess := utils.MergeAllFilesInDir(resultFilePath, tempDirPath2)
	if isMergeSuccess != true {
		log.Println("Error during final merge !")
		return false
	}

	log.Printf("File compression successful: %s", resultFilePath)

	return true
}

func (p *PdfHandler) ConvertImageToPdf(filePath string, canResize bool) bool {
	targetFilePath := path.Join(p.outputDir, utils.GetFileNameWoExtensionFromPath(filePath)+"_resized.pdf")
	tempFilePath := targetFilePath
	if canResize {
		tempFilePath = utils.GetNewTempFilePath(p.TempDir, "pdf")
		defer os.Remove(tempFilePath)
	}

	isSuccess := utils.ConvertImageToPdf(filePath, tempFilePath)
	if !isSuccess {
		return false
	}

	if !canResize {
		return true
	}

	isSuccess = utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		TargetFilePath: targetFilePath,
		SourceFilePath: tempFilePath,
	})

	if !isSuccess {
		return true // file is converted, even though not resized
	}

	return true
}

func (p *PdfHandler) ResizePdfFileToA4(filePath string) bool {
	return utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		SourceFilePath: filePath,
		TargetFilePath: path.Join(p.outputDir, utils.AddSuffixToFileName(utils.GetFileNameFromPath(filePath), "_resized")),
	})
}