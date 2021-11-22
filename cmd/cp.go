package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/AppleGamer22/recursive-backup/internal/manager"
	"github.com/AppleGamer22/recursive-backup/internal/tasks"
	"github.com/spf13/cobra"
)

var batchesDirPath string
var copyQueueLen uint
var responseChan chan tasks.BackupFileResponse
var digitsRE = regexp.MustCompile("[[:digit:]]+")

func init() {
	cpCmd.Flags().StringVarP(&rootDirPath, "project", "p", "", "mandatory flag: project root path")
	cpCmd.Flags().StringVarP(&batchesDirPath, "batches-dir-path", "b", "", "mandatory flag: copy batches directory path")
	cpCmd.Flags().UintVarP(&copyQueueLen, "copy-queue-len", "q", 5, "copy queue length")
	rootCmd.AddCommand(cpCmd)
}

var cpCmd = &cobra.Command{
	Use:   "cp [source-dir-path] [target-dir-path]",
	Short: "copy files",
	Long:  "copy files recursively from source to target dir",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("arguments mismatch, expecting 2 arguments: [source-dir-path] [target-dir-path]")
		}
		cfg.Src = args[0]
		cfg.Target = args[1]

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(rootDirPath) == 0 {
			return errors.New("rootDirPath must be specified")
		}
		if len(batchesDirPath) == 0 {
			return errors.New("batchesDirPath must be specified")
		}
		batchesToDoDirPath = filepath.Join(batchesDirPath, sliceBatchesToDoDirName)
		batchesDoneDirPath = filepath.Join(batchesDirPath, sliceBatchesDoneDirName)
		return nil
	},
	RunE: cpRunCommand,
	PostRunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return nil
		}
		doneDirWalkFunc := func(path string, d fs.DirEntry, err error) error {
			toBeRemovedPath := filepath.Join(batchesToDoDirPath, d.Name())
			if err := os.Remove(toBeRemovedPath); err != nil {
				writeOpLog(fmt.Sprintf("%s (%+v)", toBeRemovedPath, err))
			}
			return nil
		}

		_ = filepath.WalkDir(batchesDoneDirPath, doneDirWalkFunc)

		return nil
	},
}

func cpRunCommand(_ *cobra.Command, _ []string) error {
	_ = writeOpLog(fmt.Sprintf("cp start for batches in %s", batchesDirPath))
	err := filepath.WalkDir(batchesToDoDirPath, walkDirFunc)
	_ = writeOpLog("cp finished for all batches")
	return err
}

func walkDirFunc(path string, d fs.DirEntry, err error) error {
	in := manager.ServiceInitInput{
		SourceRootDir: cfg.Src,
		TargetRootDir: cfg.Target,
	}
	service := manager.NewService(in)
	switch {
	case err != nil:
		_ = writeOpLog(fmt.Sprintf("error with dir entry. path: %s. error: %s", path, err.Error()))
		return err
	case d.Type().IsDir():
		return nil
	case d.Type().IsRegular():
		_ = writeOpLog(fmt.Sprintf("cp start for batch %s", path))
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		batchFileBasePath := filepath.Base(path)
		batchIDString := digitsRE.FindString(batchFileBasePath)
		batchID, err := strconv.Atoi(batchIDString)
		if err != nil {
			return fmt.Errorf("failed to extract batch number from %s", path)
		}

		now := time.Now().Format(timeDateFormat)
		copyLogDirName := fmt.Sprintf(copyLogDirPattern, now)
		copyLogDirPath := filepath.Join(rootDirPath, copyLogDirName)
		if err = os.MkdirAll(copyLogDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create copy log Dir. Error: %v", err)
		}
		copyLogFileName := fmt.Sprintf(copyBatchLogFileNamePattern, batchID)
		copyLogFilePath := filepath.Join(copyLogDirPath, copyLogFileName)
		copyLogFile, err := os.Create(copyLogFilePath)
		if err != nil {
			return fmt.Errorf("failed to create copy log file. Error: %v", err)
		}
		fmt.Println(copyLogFilePath)

		responseChan = make(chan tasks.BackupFileResponse, copyQueueLen)
		go service.HandleFilesCopyResponse(copyLogFile, responseChan)
		service.RequestFilesCopy(file, responseChan)
		close(responseChan)
		_ = file.Close()
		donePath := filepath.Join(batchesDoneDirPath, batchFileBasePath)

		if err = os.Rename(path, donePath); err != nil {
			_ = writeOpLog(fmt.Sprintf("failed to move batch file to done dir %s (%v)", path, err))
		} else {
			fmt.Printf("%s -> %s\n", path, donePath)
		}
		_ = writeOpLog(fmt.Sprintf("cp finished for batch %s\n", path))
	default:
		return nil
	}

	return nil
}
