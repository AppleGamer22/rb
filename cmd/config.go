package cmd

import (
	"fmt"
	"time"
)

const timeFormat = "2006-01-02T15:04:05"

type rootConfig struct {
	NumWorkers            uint
	Src                   string
	Target                string
	RecoveryReferenceTime time.Time
}
var cfg rootConfig

func parseTime(timeString string) (*time.Time, error) {
	assertedTime, err := time.Parse(timeFormat, timeString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time flag value")
	}
	if !assertedTime.Before(time.Now()) {
		return nil, fmt.Errorf("reference time flag value is in the future")
	}
	return &assertedTime, nil
}

