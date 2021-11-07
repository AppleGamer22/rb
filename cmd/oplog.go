package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func writeOpLog(s string) error {
	operationLogLine := fmt.Sprintf("%s %s\n", s, time.Now().Format(time.RFC1123Z))
	absPath := filepath.Join(rootDirPath, operationLogFileName)

	opLog, err := os.OpenFile(absPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.New("failed to open operation log")
	}
	defer func() {
		_ = opLog.Close()
	}()

	if _, err = opLog.WriteString(operationLogLine); err != nil {
		return err
	}

	return nil
}
