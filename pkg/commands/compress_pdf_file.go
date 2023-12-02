package commands

import (
	"log/slog"
	"os"
	"path"
	wlptypes "welovepdf/pkg/types"
	utils "welovepdf/pkg/utils"

	"github.com/google/uuid"
)

// TODO Split by extracting out CompressMultiPagePdfFile
func BuildCompressMultiPagePdfFile(
	compressSinglePageFile func(targetImageQuality int, config *wlptypes.FileToFileOperationConfig) error,
	splitPdfFile wlptypes.FileToDirOperation,
	mergePdfFiles wlptypes.FilesToFileOperation,
	tempDirPath string,
) func(targetImageQuality int, compressionConfig *wlptypes.FileToFileOperationConfig) error {
	mergeAllFilesInDir := buildMergeAllFilesInDir(mergePdfFiles)
	compressAllFilesInDir := buildCompressAllFilesInDir(tempDirPath, compressSinglePageFile)

	return func(
		quality int,
		config *wlptypes.FileToFileOperationConfig) error {
		operationId := uuid.New().String()
		tempDirPath1 := path.Join(tempDirPath, operationId, "compress_jpg")
		tempDirPath2 := path.Join(tempDirPath, operationId, "compress_pdf")
		utils.EnsureDirectory(tempDirPath1)
		utils.EnsureDirectory(tempDirPath2)
		defer os.RemoveAll(tempDirPath1)
		defer os.RemoveAll(tempDirPath2)

		// 1. Split file into 1 file per page
		err := splitPdfFile(&wlptypes.FileToDirOperationConfig{
			SourceFilePath: config.SourceFilePath,
			TargetDirPath:  tempDirPath1,
		})
		if err != nil {
			slog.Error("CompressFile : error splitting file, compression aborted", slog.String("reason", err.Error()))
			return err
		}

		// 2. Compress the splitted files
		err = compressAllFilesInDir(quality, &wlptypes.DirToDirOperationConfig{
			SourceDirPath: tempDirPath1,
			TargetDirPath: tempDirPath2,
		})
		if err != nil {
			slog.Error("CompressFile : error compressing files in dir", slog.String("tempDirPath1", tempDirPath1), slog.String("tempDirPath2", tempDirPath2), slog.String("reason", err.Error()))
			return err
		}

		// 3. Merge all compressed files back into 1
		err = mergeAllFilesInDir(&wlptypes.DirToFileOperationConfig{
			SourceDirPath:  tempDirPath2,
			TargetFilePath: config.TargetFilePath,
		})
		if err != nil {
			slog.Error("CompressFile : error during final merge !", slog.String("reason", err.Error()))
			return err
		}

		slog.Debug("CompressFile operation succeeded")

		return nil
	}
}
