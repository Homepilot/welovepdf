package ghostscript

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
	wlptypes "welovepdf/pkg/types"
)

type GhoscriptClient struct {
	binaryPath         string
	viewJpegScriptPath string
}

func NewGhostscriptClient(
	binaryPath string, viewJpegScriptPath string) *GhoscriptClient {
	return &GhoscriptClient{
		binaryPath:         binaryPath,
		viewJpegScriptPath: viewJpegScriptPath,
	}
}

func (c *GhoscriptClient) convertToLowQualityJpeg(targetImageQuality int, config *wlptypes.FileToFileOperationConfig) error {
	log.Printf("converting w/ GS using quality %d, binaryPath '%s', source '%s', target '%s'", targetImageQuality, c.binaryPath, config.SourceFilePath, config.TargetFilePath)
	convertToLowQualityJpegCmd := exec.Command(
		c.binaryPath,
		"-sDEVICE=jpeg",
		"-o",
		config.TargetFilePath,
		"-dJPEGQ="+fmt.Sprintf("%d", targetImageQuality),
		"-dNOPAUSE",
		"-dBATCH",
		"-dUseCropBox",
		"-dTextAlphaBits=4",
		"-dGraphicsAlphaBits=4",
		"-r140",
		config.SourceFilePath)
	err := convertToLowQualityJpegCmd.Run()
	return err
}

func (c *GhoscriptClient) ConvertJpegToPdf(config *wlptypes.FileToFileOperationConfig) error {
	convertCmd := exec.Command(
		c.binaryPath,
		"-dNOSAFER",
		"-sDEVICE=pdfwrite",
		"-o",
		config.TargetFilePath,
		c.viewJpegScriptPath,
		"-c",
		"("+config.SourceFilePath+")",
		"viewJPEG",
	)
	slog.Info("the printed string", slog.String("the string", convertCmd.String()))

	err := convertCmd.Run()
	return err
}

func (c *GhoscriptClient) ResizePdfToA4(config *wlptypes.FileToFileOperationConfig) error {
	resizePdfToA4Cmd := exec.Command(
		c.binaryPath,
		"-o",
		config.TargetFilePath,
		"-sDEVICE=pdfwrite",
		"-sPAPERSIZE=a4",
		"-dFIXEDMEDIA",
		"-dPDFFitPage",
		"-dCompatibilityLevel=1.4",
		config.SourceFilePath)

	err := resizePdfToA4Cmd.Run()
	return err
}

func (c *GhoscriptClient) MergePdfFiles(config *wlptypes.FilesToFileOperationConfig) error {
	mergePdfFilesCmd := exec.Command(
		c.binaryPath,
		"-dNOPAUSE",
		"-sDEVICE=pdfwrite",
		"-sOUTPUTFILE="+config.TargetFilePath,
		"-dBATCH",
	)
	mergePdfFilesCmd.Args = append(mergePdfFilesCmd.Args, config.SourceFilesPathes...)

	err := mergePdfFilesCmd.Run()
	return err
}

func (c *GhoscriptClient) MergeAllFilesInDir(config *wlptypes.DirToFileOperationConfig) error {
	filesToMerge, err := os.ReadDir(config.SourceDirPath)
	if err != nil {
		log.Printf("Error reading temp dir to merge: %s", err.Error())
		return err
	}
	if len(filesToMerge) < 1 {
		log.Println("No files to merge, aborting")
		return nil
	}

	log.Printf("found %d compressed files to merge", len(filesToMerge))
	filesPathesToMerge := []string{}
	for _, file := range filesToMerge {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			filesPathesToMerge = append(filesPathesToMerge, path.Join(config.SourceDirPath, file.Name()))
		}
	}

	return c.MergePdfFiles(&wlptypes.FilesToFileOperationConfig{
		SourceFilesPathes: filesPathesToMerge,
		TargetFilePath:    config.TargetFilePath,
	})
}

func (c *GhoscriptClient) SplitPdfFile(config *wlptypes.FileToDirOperationConfig) error {
	splitPdfFileCmd := exec.Command(
		c.binaryPath,
		"-sDEVICE=pdfwrite",
		"-dSAFER",
		"-o",
		path.Join(config.TargetDirPath, "outfile.%d.pdf"),
		config.SourceFilePath,
	)

	err := splitPdfFileCmd.Run()
	return err
}
