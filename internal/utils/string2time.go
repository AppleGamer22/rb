package utils

import (
	"bufio"
	"os"
	"strings"
	"time"
)

// Reads the first valid date from a given CSV's 3rd column
func GetLastBackupExecutionTime(previousExecutionLogPath string) (time.Time, error) {
	file, err := os.Open(previousExecutionLogPath)
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return time.Unix(0, 0), err
		} else if strings.HasPrefix(line, "ERROR: ") {
			continue
		}
		timezone := time.Now().Location()
		timeColumn := strings.Split(line, ",")[2]
		timeString := strings.Split(timeColumn, " +")[0]
		t, _ := time.Parse("2006-01-02 15:04:05", timeString)
		adjustedTime := t.In(timezone)
		return adjustedTime, nil
	}
}