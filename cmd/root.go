package cmd

import (
	"fmt"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cobra"
)

//rb --num-workers <N> --op-mode <initial | recovery> --recovery-reference <dir-path> <src> <target>

// flags
const (
	initial  = "initial"
	recovery = "recovery"
)

type config struct {
	NumWorkers        uint
	OpMode            string
	RecoveryReference string
	Src               string
	Target            string
}

func checkDirectoryPath(path interface{}) error {
	assertedPath, ok := path.(string)
	if !ok {
		return fmt.Errorf("%v must be a valid directory path", path)
	}
	fileInfo, err := os.Stat(assertedPath)
	if err != nil || fileInfo.IsDir() {
		return fmt.Errorf("%s must be a valid directory path", assertedPath)
	}
	return nil
}

func (cfg config) Validate() error {
	return validation.ValidateStruct(&cfg,
		validation.Field(&cfg.NumWorkers, validation.Min(1)),
		validation.Field(&cfg.OpMode, validation.In(initial, recovery)),
		validation.Field(&cfg.RecoveryReference, validation.By(checkDirectoryPath)),
		validation.Field(&cfg.Src, validation.By(checkDirectoryPath)),
		validation.Field(&cfg.Target, validation.By(checkDirectoryPath)),
	)
}

var rootCmd = &cobra.Command{
	Use:   "rb",
	Short: "rb is a tool for backing up files over unreliable network connections",
	Long:  "rb is a tool for backing up files over unreliable network connections",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	var cfg config
	rootCmd.PersistentFlags().StringVar(&cfg.OpMode, "op-mode", "", "--op-mode < initial | recovery >")

	if err := cfg.Validate(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
