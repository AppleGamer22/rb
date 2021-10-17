package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)



func writeOpLog(s string) error {
	operationLogLine := []byte(fmt.Sprintf("%s %s\n", s, time.Now().String()))
	absPath := filepath.Join(rootDirPath, operationLogFileName)
	if err := os.WriteFile(absPath, operationLogLine, defaultPerm); err != nil {
		return errors.New("failed to write to operation log")
	}
	return nil
}
