package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var fullCmd = &cobra.Command{
	Use:   "full",
	Short: "full backup",
	Long:  "with full backup all files and folders are copied from src to target",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("arguments mismatch, expecting 2 arguments")
		}
		cfg.Src = args[0]
		cfg.Target = args[1]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("numWorkers: %v\n", cfg.NumWorkers)
		fmt.Printf("src: %v\n", cfg.Src)
		fmt.Printf("target: %v\n", cfg.Target)
	},
}

func init() {
	fullCmd.PersistentFlags().UintVarP(&cfg.NumWorkers, "workers", "w", 2, "number of concurrent file-copy workers")
	rootCmd.AddCommand(fullCmd)
}
