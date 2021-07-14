package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"
)

func main() {
	var logsFlag = flag.String("logs", "", "")
	flag.Parse()
	var source, err1 = filepath.Abs(flag.Arg(0))
	if err1 != nil {
		log.Fatal("source path is invalid")
	}
	var target, err2 = filepath.Abs(flag.Arg(1))
	if err2 != nil {
		log.Fatal("target path is not valid")
	}
	logsPath := string(*logsFlag)

	var files []FileMetadata
	if logsPath != "" {
		logs, err := GetLogFromFile(logsPath)
		if err != nil {
			log.Fatal(err)
		}
		files, err = GetFilePathsSinceDate(source, target, logs.LastBackupTime)
		if err != nil {
			log.Fatal(err)
		}
		var err4 = SaveMetadataToFile(files, "files.json", 0, false, logs.LastBackupTime)
		if err4 != nil {
			log.Fatal(err4)
		}
	} else {
		files, err2 = GetFilePaths(source, target)
		if err2 != nil {
			log.Fatal(err2)
		}
		var err4 = SaveMetadataToFile(files, "files.json", 0, false, time.Unix(0, 0))
		if err4 != nil {
			log.Fatal(err4)
		}
	}
	var err = Backup("files.json", target)
	if err != nil {
		log.Fatal(err)
	}
}
