package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var timeString string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "differential backup",
	Long:  "with differential backup only new or modified files and directories are being backed-up",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("arguments mismatch, expecting 2 arguments")
		}
		cfg.Src = args[0]
		cfg.Target = args[1]

		assertedTime, err := time.Parse(timeFormat, timeString)
		if err != nil {
			return fmt.Errorf("failed to parse time flag value")
		}
		if !assertedTime.Before(time.Now()) {
			return fmt.Errorf("reference time flag value is in the future")
		}
		cfg.RecoveryReferenceTime = assertedTime

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("numWorkers: %v\n", cfg.NumWorkers)
		fmt.Printf("src: %v\n", cfg.Src)
		fmt.Printf("target: %v\n", cfg.Target)
	},
}

func init() {
	diffCmd.PersistentFlags().UintVarP(&cfg.NumWorkers, "workers", "w", 2, "number of concurrent file-copy workers")
	diffCmd.PersistentFlags().StringVarP(&timeString, "time", "t", "", "reference time with format: 2006-01-02T15:04:05")
	rootCmd.AddCommand(diffCmd)
}
