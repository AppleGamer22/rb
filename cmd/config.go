package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const timeFormat = "2006-01-02T15:04:05"

type rootConfig struct {
	NumWorkers     uint       `yaml:"num_workers"`
	BatchSize      uint       `yaml:"batch_size"`
	DirsListPath   string     `yaml:"dir_list_path"`
	BatchesDirPath string     `yaml:"batches_dir_path"`
	FilesListPath  string     `yaml:"file_list_path"`
	Source         string     `yaml:"source"`
	Target         string     `yaml:"target"`
	ProjectDir     string     `yaml:"project_dir"`
	ReferenceTime  *time.Time `yaml:"reference_time"`
}

var cfg rootConfig

func parseTime(timeString string) (*time.Time, error) {
	assertedTime, err := time.Parse(timeDateFormat, timeString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time flag value")
	}
	if !assertedTime.Before(time.Now()) {
		return nil, fmt.Errorf("reference time flag value is in the future")
	}
	return &assertedTime, nil
}

func readConfigFile() error {
	viper.SetConfigName("rb")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}
