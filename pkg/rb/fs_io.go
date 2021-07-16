package rb

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func BackupLog(fileNum, numberOfFiles int, sourcePath, targetPath string) {
	fmt.Printf("file #%d / %d (%s -> %s)\n", fileNum, numberOfFiles, sourcePath, targetPath)
}

func Backup(sourcesLogPath, sourcePathRoot, targetPathRoot string, fileCount int, startTime time.Time) error {
	fileSourcesLog, err := os.Open(sourcesLogPath)
	if err != nil {
		return err
	}
	defer fileSourcesLog.Close()
	var reader = bufio.NewReader(fileSourcesLog)
	var targetsLogPath = fmt.Sprintf("rb_target_%d-%d-%d_%d.%d.%d.csv", startTime.Day(), startTime.Month(), startTime.Year(), startTime.Hour(), startTime.Minute(), startTime.Second())
	fileTargetsLog, err := os.Create(targetsLogPath)
	if err != nil {
		return err
	}
	defer fileTargetsLog.Close()
	var writer = bufio.NewWriter(fileTargetsLog)
	for i := 1; ; i++ {
		data, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		var sourcePath = strings.Replace(string(data), "\n", "", -1)
		if strings.HasPrefix(sourcePath, "ERROR: ") {
			continue
		}
		relativePath, err := filepath.Rel(sourcePathRoot, sourcePath)
		if err != nil {
			return err
		}
		targetPath, err := filepath.Abs(fmt.Sprintf("%s/%s", targetPathRoot, relativePath))
		if err != nil {
			return err
		}
		BackupLog(i, fileCount, sourcePath, targetPath)
		modTime, err := Copy(sourcePath, targetPath, targetPathRoot)
		if err != nil {
			return err
		}
		writer.WriteString(fmt.Sprintf("%s,%s,%s\n", sourcePath, targetPath, modTime))
		writer.Flush()
		fmt.Println("done")
	}
	return nil
}

func Copy(sourcePath, targetPath string, targetPathRoot string) (time.Time, error) {
	fileStatSource, err := os.Stat(sourcePath)
	if err != nil {
		WaitForDirectory(targetPathRoot)
		return Copy(sourcePath, targetPath, targetPathRoot)
	}
	if !fileStatSource.Mode().IsRegular() {
		return time.Unix(0, 0), fmt.Errorf("%s is not a regular file", sourcePath)
	}
	_, err = os.Stat(targetPath)
	if err != nil {
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			WaitForDirectory(targetPathRoot)
			return Copy(sourcePath, targetPath, targetPathRoot)
		}
	}
	src, err := os.Open(sourcePath)
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer src.Close()

	dest, err := os.Create(targetPath)

	if err != nil {
		WaitForDirectory(targetPathRoot)
		return Copy(sourcePath, targetPath, targetPathRoot)
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return fileStatSource.ModTime(), err
}

func WaitForDirectory(path string) {
	fmt.Printf("Waiting for directory %s to be available...\n", path)
	var searching = true
	for searching {
		_, err := os.Stat(path)
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			searching = false
		}
	}
}
