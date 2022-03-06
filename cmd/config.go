package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// const timeFormat = "2006-01-02T15:04:05"

type rootConfig struct {
	NumWorkers          uint
	BatchSize           uint
	DirsListPath        string
	BatchesDirPath      string
	FilesListPath       string
	Source              string
	Target              string
	ProjectDir          string
	ReferenceTimeString string
	DirValidationMode   string
	// ReferenceTime       *time.Time
}

var cfg rootConfig

func parseTime(timeString string) (*time.Time, error) {
	assertedTime, err := time.Parse(timeDateFormat, timeString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time flag value: %w", err)
	}
	if !assertedTime.Before(time.Now()) {
		return nil, fmt.Errorf("reference time flag value is in the future: %w", err)
	}
	return &assertedTime, nil
}

func readConfigFile() error {
	viper.SetConfigName("rb")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}
