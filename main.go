package main

import (
	"flag"
	"log"
	"pbed/xgithub"
)

var (
	uploadDirPath  string
	uploadFilePath string
	gui            bool
)

func init() {
	flag.StringVar(&uploadDirPath, "d", "", "upload dir files, example: d ./")
	flag.StringVar(&uploadFilePath, "p", "", "upload file path: example: p ./a.txt")
	flag.BoolVar(&gui, "gui", false, "start a pic bed gui")

	flag.Parse()
}

func main() {
	// start a server
	if gui {
		var xb = xgithub.New()
		startGUI(xb)

		return
	}

	// upload dir
	if uploadDirPath != "" {
		var xb = xgithub.New()
		uploadDir(xb, uploadDirPath)

		return
	}

	// upload sigle file
	if uploadFilePath != "" {
		var xb = xgithub.New()
		upload(xb, uploadFilePath)

		return
	}

	log.Println("use -h to get help")
}
