package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "rb",
	Short: "rb backup tool",
	Long:  "rb is a tool for backing up files over unreliable network connections",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}


//func checkRefTime(t interface{}) error {
//	assertedStrVal, ok := t.(string)
//	if !ok {
//		return fmt.Errorf("%v must be a time string value", t)
//	}
//	assertedTime, err := time.Parse(timeFormat, assertedStrVal)
//	if err != nil {
//		return fmt.Errorf("%v could not be parsed into a time.Time value", assertedStrVal)
//	}
//	if !assertedTime.Before(time.Now()) {
//		return fmt.Errorf("%v is not before now", assertedTime)
//	}
//	return nil
//}

//func checkMinUintValue(n interface{}) error {
//	assertedValue, ok := n.(uint)
//	if !ok {
//		return fmt.Errorf("%v must be a uint value", n)
//	}
//	if assertedValue < 1 {
//		return fmt.Errorf("%v must be greater than 0", n)
//	}
//	return nil
//}
//
//func checkDirectoryPath(path interface{}) error {
//	assertedPath, ok := path.(string)
//	if !ok {
//		return fmt.Errorf("%v must be a valid directory path", path)
//	}
//	fileInfo, err := os.Stat(assertedPath)
//	if err != nil || fileInfo.IsDir() {
//		return fmt.Errorf("%s must be a valid directory path", assertedPath)
//	}
//	return nil

