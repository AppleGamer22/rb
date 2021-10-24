package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const defaultBatchSize = 1000

var batchesTargetDirPath string
var batchesErrorsDirPath string
var filesListFilePath string
var batchSize uint

func init() {
	sliceCmd.Flags().StringVarP(&rootDirPath, "project", "p", "", "project root path")
	sliceCmd.Flags().StringVarP(&filesListFilePath, "files-list", "f", "", "files list file path")
	sliceCmd.Flags().UintVarP(&batchSize, "batch-size", "b", defaultBatchSize, "maximum number of files in a batch")
	rootCmd.AddCommand(sliceCmd)

}

var sliceCmd = &cobra.Command{
	Use:   "slice",
	Short: "slice list of files to copy into smaller batches",
	Long:  "create directory skeleton in target",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("arguments mismatch, no argument expected")
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		batchesTargetDirPath = filepath.Join(rootDirPath, sliceBatchesDirName)
		batchesErrorsDirPath = filepath.Join(rootDirPath, sliceBatchesErrorDir)
		return nil
	},
	RunE: sliceRunCommand,
}

func sliceRunCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("src: %v\n", cfg.Src)

	operationLogLine := "slice copy batches start"
	if err := writeOpLog(operationLogLine); err != nil {
		return err
	}

	if err := os.Chdir(skeletonWorkDir); err != nil {
		return err
	}

	inFilesListFile, errorsFile, err := setupForSlice()
	if err != nil {
		return err
	}
	defer func() {
		_ = inFilesListFile.Close()
		_ = errorsFile.Close()
	}()

	scan := bufio.NewScanner(inFilesListFile)
	var batchCounter, lineCounter uint = 0, 0
	for scan.Scan() {
		if lineCounter == 0 {
			batchCounter++
			errorsFileName := fmt.Sprintf(sliceErrorsFileNamePattern, time.Now().Format(timeDateFormat))
			errorsFilePath := filepath.Join(rootDirPath, sliceBatchesErrorDir, errorsFileName)
		}
		text := scan.Text()
		//dirs = append(dirs, text)
		lineCounter = (lineCounter + 1) % batchSize
	}
	if err = scan.Err(); err != nil {

		return nil, err
	}

	operationLogLine = "slice copy batches end"
	if err = writeOpLog(operationLogLine); err != nil {
		return err
	}

	return nil
}

func setupForSlice() (inFilesList, errs *os.File, err error) {
	inFilesList, err = os.Open(filesListFilePath)
	if err != nil {
		return nil, nil, err
	}

	errorsFileName := fmt.Sprintf(sliceErrorsFileNamePattern, time.Now().Format(timeDateFormat))
	errorsFilePath := filepath.Join(rootDirPath, sliceBatchesErrorDir, errorsFileName)
	errs, err = os.Create(errorsFilePath)
	if err != nil {
		return nil, nil, err
	}

	return inFilesList, errs, nil
}
