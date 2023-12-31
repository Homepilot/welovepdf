package ghostscript

import (
	"fmt"
	"os/exec"
	"path"
	wlptypes "welovepdf/pkg/types"
)

type GhostScriptCommander struct {
	binaryPath         string
	viewJpegScriptPath string
}

func NewGhostscriptClient(
	binaryPath string, viewJpegScriptPath string) *GhostScriptCommander {
	return &GhostScriptCommander{
		binaryPath:         binaryPath,
		viewJpegScriptPath: viewJpegScriptPath,
	}
}

func (c *GhostScriptCommander) ConvertPdfToJpeg(targetImageQuality int, config *wlptypes.FileToFileOperationConfig) error {
	convertPdfToJpegCmd := exec.Command(
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
	err := convertPdfToJpegCmd.Run()
	return err
}

func (c *GhostScriptCommander) ConvertJpegToPdf(config *wlptypes.FileToFileOperationConfig) error {
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

	err := convertCmd.Run()
	return err
}

func (c *GhostScriptCommander) ResizePdfToA4(config *wlptypes.FileToFileOperationConfig) error {
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

func (c *GhostScriptCommander) MergePdfFiles(config *wlptypes.FilesToFileOperationConfig) error {
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

func (c *GhostScriptCommander) SplitPdfFile(config *wlptypes.FileToDirOperationConfig) error {
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
