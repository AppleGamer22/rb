package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var ltCmd = &cobra.Command{
	Use:   "lt",
	Short: "list recent modifications in source",
	Long:  "list source recursively for recent modifications in directories and files",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("arguments mismatch, expecting 1 argument")
		}
		cfg.Src = args[0]

		assertedTime, err := parseTime(timeString)
		if err != nil {
			return err
		}
		cfg.RecoveryReferenceTime = *assertedTime

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("numWorkers: %v\n", cfg.NumWorkers)
		fmt.Printf("src: %v\n", cfg.Src)
		fmt.Printf("target: %v\n", cfg.Target)
	},
}

func init() {
	ltCmd.PersistentFlags().StringVarP(&timeString, "time", "t", "", "reference time with format: 2006-01-02T15:04:05")
	rootCmd.AddCommand(ltCmd)
}
