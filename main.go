package main

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
	"pbed/cons"
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

	initConfig()
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

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	fn := home + cons.ConfigPath

	_, err = os.Stat(fn)
	if err != nil && os.IsNotExist(err) {
		log.Println("init config")

		f, err := os.Create(fn)
		if err != nil {
			panic(err)
		}

		bs, err := xgithub.ConfigTemplate()
		if err != nil {
			panic(err)
		}

		_, err = f.Write(bs)

		if err != nil {
			panic(err)
		}
	}
}
