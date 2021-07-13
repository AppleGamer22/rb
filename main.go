package main

import (
	"flag"
	"log"
	"path/filepath"
)

func main() {
	var sourceFlag = flag.String("source", "", "")
	var targetFlag = flag.String("target", "", "")
	flag.Parse()
	var source, err1 = filepath.Abs(string(*sourceFlag))
	if err1 != nil {
		log.Fatal("source path is invalid")
	}
	var target, err2 = filepath.Abs(string(*targetFlag))
	if err2 != nil {
		log.Fatal("target path is not valid")
	}
	var files, err3 = GetFilePaths(source, target)
	if err3 != nil {
		log.Fatal(err3)
	}
	var err4 = SaveMetadataToFile(files, "files.json")
	if err4 != nil {
		log.Fatal(err4)
	}
	var err5 = MarkAsDone("files.json", 0)
	if err5 != nil {
		log.Fatal(err5)
	}
}
