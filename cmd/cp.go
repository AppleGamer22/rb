package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/AppleGamer22/recursive-backup/internal/workers"

	"github.com/AppleGamer22/recursive-backup/internal/manager"
	"github.com/AppleGamer22/recursive-backup/internal/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var generalRequestChannel chan tasks.GeneralRequest
var responseChan chan tasks.BackupFileResponse
var wgCopyWorkerQuitConfirmation sync.WaitGroup
var digitsRE = regexp.MustCompile("[[:digit:]]+")
var in manager.ServiceInitInput
var service manager.API

func UpdateOnQuit() {
	wgCopyWorkerQuitConfirmation.Done()
}

func init() {
	cpCmd.Flags().StringVarP(&cfg.ProjectDir, "project", "p", "", "mandatory flag: project root path")
	cpCmd.Flags().StringVarP(&cfg.BatchesDirPath, "batches-dir-path", "b", "", "mandatory flag: copy batches directory path")
	cpCmd.Flags().UintVarP(&cfg.NumWorkers, "copy-queue-len", "q", 200, "copy queue length")

	viper.BindPFlag("ProjectDir", cpCmd.Flags().Lookup("project"))
	viper.BindPFlag("BatchesDirPath", cpCmd.Flags().Lookup("batches-dir-path"))
	viper.BindPFlag("NumWorkers", cpCmd.Flags().Lookup("copy-queue-len"))
	rootCmd.AddCommand(cpCmd)
}

var cpCmd = &cobra.Command{
	Use:   "cp [source-dir-path] [target-dir-path]",
	Short: "copy files",
	Long:  "copy files recursively from source to target dir",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			if len(cfg.Source) == 0 && len(cfg.Target) == 0 {
				return errors.New("arguments mismatch, expecting 2 arguments: [source-dir-path] [target-dir-path]")
			} else if (len(cfg.Source) > 0 || len(cfg.Target) > 0) && !(len(cfg.Source) > 0 && len(cfg.Target) > 0) {
				return errors.New("arguments mismatch, expecting 2 arguments: [source-dir-path] [target-dir-path]")
			}
		} else {
			cfg.Source = args[0]
			cfg.Target = args[1]
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(cfg.ProjectDir) == 0 {
			return errors.New("cfg.ProjectDir must be specified")
		}
		if len(cfg.BatchesDirPath) == 0 {
			return errors.New("cfg.BatchesDirPath must be specified")
		}
		batchesToDoDirPath = filepath.Join(cfg.BatchesDirPath, sliceBatchesToDoDirName)
		batchesDoneDirPath = filepath.Join(cfg.BatchesDirPath, sliceBatchesDoneDirName)
		in = manager.ServiceInitInput{
			SourceRootDir: cfg.Source,
			TargetRootDir: cfg.Target,
		}
		service = manager.NewService(in)

		return nil
	},
	RunE: cpRunCommand,
}

func cpRunCommand(_ *cobra.Command, _ []string) error {
	_ = writeOpLog(fmt.Sprintf("cp start for batches in %s", cfg.BatchesDirPath))

	responseChan = make(chan tasks.BackupFileResponse, cfg.NumWorkers)
	defer func() {
		service.WaitForAllResponses()
		close(responseChan)
	}()
	generalRequestChannel = make(chan tasks.GeneralRequest, cfg.NumWorkers)
	defer func() {
		wgCopyWorkerQuitConfirmation.Wait()
		close(generalRequestChannel)
	}()
	for i := 1; i <= int(cfg.NumWorkers); i++ {
		wgCopyWorkerQuitConfirmation.Add(1)
		workers.NewCopyWorker(uint(i), cfg.Source, cfg.Target, generalRequestChannel, UpdateOnQuit)
	}

	err := filepath.WalkDir(batchesToDoDirPath, walkDirFunc)
	_ = writeOpLog("cp finished for all batches")

	for i := 0; i < int(cfg.NumWorkers); i++ {
		generalRequestChannel <- tasks.QuitRequest{}
	}
	return err
}

func walkDirFunc(path string, d fs.DirEntry, err error) error {
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
		copyLogDirPath := filepath.Join(cfg.ProjectDir, copyLogDirName)
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

		go service.HandleFilesCopyResponse(copyLogFile, responseChan)
		service.RequestFilesCopy(file, uint(batchID), generalRequestChannel, responseChan)
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
