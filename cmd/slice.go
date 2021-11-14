package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const defaultBatchSize = 1000

var batchesSourceDirPath string
var batchesToDoDirPath string
var batchesDoneDirPath string
var batchesErrorsDirPath string
var filesListFilePath string
var batchSize uint

func init() {
	sliceCmd.Flags().StringVarP(&rootDirPath, "project", "p", "", "mandatory flag: project root path")
	sliceCmd.Flags().StringVarP(&filesListFilePath, "files-list-file-path", "f", "", "mandatory flag: files list file path")
	sliceCmd.Flags().UintVarP(&batchSize, "batch-size", "s", defaultBatchSize, "maximum number of files in a batch")
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
		if len(rootDirPath) == 0 {
			return errors.New("project root path flag must be specified")
		}
		if len(filesListFilePath) == 0 {
			return errors.New("files-list-file-path flag must be specified")
		}

		now := time.Now().Format(timeDateFormat)
		batchesDirName := fmt.Sprintf(sliceBatchesDirNamePattern, now)
		batchesSourceDirPath = filepath.Join(rootDirPath, batchesDirName)
		batchesDoneDirPath = filepath.Join(batchesSourceDirPath, sliceBatchesDoneDirName)
		if err := os.MkdirAll(batchesDoneDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create batches target dir. %s", err.Error())
		}

		batchesToDoDirPath = filepath.Join(batchesSourceDirPath, sliceBatchesToDoDirName)
		if err := os.MkdirAll(batchesToDoDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create batches source dir. %s", err.Error())
		}

		batchesErrorDirName := fmt.Sprintf(sliceBatchesErrorDirPattern, now)
		batchesErrorsDirPath = filepath.Join(rootDirPath, batchesErrorDirName)
		if err := os.MkdirAll(batchesErrorsDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create batches errors dir. %s", err.Error())
		}
		return nil
	},
	RunE: sliceRunCommand,
}

func sliceRunCommand(_ *cobra.Command, _ []string) error {
	operationLogLine := "slice create batches start"
	if err := writeOpLog(operationLogLine); err != nil {
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

	err = sliceFileCopyBatches(inFilesListFile, errorsFile)
	if err != nil {
		return err
	}

	operationLogLine = "slice copy batches end"
	if err = writeOpLog(operationLogLine); err != nil {
		return err
	}

	return nil
}

func sliceFileCopyBatches(inFilesListFile *os.File, errorsFile *os.File) error {
	var batchCounter, lineCounter uint
	var batchFile *os.File
	var writer *bufio.Writer
	var err error
	scanner := bufio.NewScanner(inFilesListFile)
	for scanner.Scan() {
		if lineCounter == 0 {
			batchCounter++
			if batchCounter > 1 {
				_ = writer.Flush()
				_ = batchFile.Close()
			}
			batchFileName := fmt.Sprintf(sliceBatchFileNamePattern, batchCounter)
			batchFilePath := filepath.Join(batchesToDoDirPath, batchFileName)
			batchFile, err = os.Create(batchFilePath)
			if err != nil {
				_, _ = fmt.Fprintf(errorsFile, "failed to create batch file. batch_number: %d\n", batchCounter)
				lineCounter = (lineCounter + 1) % batchSize
				continue
			}
			fmt.Println(batchFilePath)
			writer = bufio.NewWriter(batchFile)
		}

		line := scanner.Text()
		if _, err = fmt.Fprintln(writer, line); err != nil {
			_, _ = fmt.Fprintf(errorsFile, "failed to write line. line: %s, error: %s\n", line, err.Error())
			lineCounter = (lineCounter + 1) % batchSize
			continue
		}

		lineCounter = (lineCounter + 1) % batchSize
	}
	_ = writer.Flush()
	_ = batchFile.Close()
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("files list scanner failed. Error:  %v", err)
	}
	return nil
}

func setupForSlice() (inFilesList, sliceErrorsFile *os.File, err error) {
	inFilesList, err = os.Open(filesListFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open input list file. %s", err)
	}

	errorsFileName := fmt.Sprintf(sliceErrorsFileNamePattern, time.Now().Format(timeDateFormat))
	errorsFilePath := filepath.Join(batchesErrorsDirPath, errorsFileName)
	sliceErrorsFile, err = os.Create(errorsFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create slice errors file. %s", err)
	}
	fmt.Println(errorsFilePath)

	return inFilesList, sliceErrorsFile, nil
}
