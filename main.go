package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
)

func main() {
	var sourceFlag = flag.String("source", "", "")
	var targetFlag = flag.String("target", "", "")
	var logsFlag = flag.String("logs", "", "")
	flag.Parse()
	var source, err1 = filepath.Abs(string(*sourceFlag))
	if err1 != nil {
		log.Fatal("source path is invalid")
	}
	var target, err2 = filepath.Abs(string(*targetFlag))
	if err2 != nil {
		log.Fatal("target path is not valid")
	}
	logs, err := filepath.Abs(string(*logsFlag))
	fmt.Println(logs)
	if err != nil {
		log.Fatal("logs file path is not valid")
	}
	var files, err3 = GetFilePaths(source, target)
	if err3 != nil {
		log.Fatal(err3)
	}
	var err4 = SaveMetadataToFile(files, "files.json", 0, false)
	if err4 != nil {
		log.Fatal(err4)
	}
	err = Backup("files.json", target)
	if err != nil {
		log.Fatal(err)
	}
}
