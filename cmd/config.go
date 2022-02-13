package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// const timeFormat = "2006-01-02T15:04:05"

type rootConfig struct {
	NumWorkers          uint   `ini:"num_workers"`
	BatchSize           uint   `ini:"batch_size"`
	DirsListPath        string `ini:"dir_list_path"`
	BatchesDirPath      string `ini:"batches_dir_path"`
	FilesListPath       string `ini:"file_list_path"`
	Source              string `ini:"source"`
	Target              string `ini:"target"`
	ProjectDir          string `ini:"project_dir"`
	ReferenceTimeString string `ini:"reference_time"`
	DirValidationMode   string `ini:"dir_validation_mode"`
	// ReferenceTime       *time.Time `ini:"-"`
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
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}
