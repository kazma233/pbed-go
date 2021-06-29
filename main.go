package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"pbed/cons"
	"pbed/xgithub"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func InitConfig() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	fn := home + cons.ConfigPath

	_, err = os.Stat(fn)
	if err != nil && os.IsNotExist(err) {
		f, err := os.Create(fn)
		if err != nil {
			panic(err)
		}

		_, err = f.WriteString(`{
	"repo":"",
	"owner":"",
	"token":"",
	"branch":"",
	"prefixFormat":""
}`)

		if err != nil {
			panic(err)
		}
	}
}

func main() {
	init := flag.Bool("init", false, "init config file")
	allDir := flag.Bool("all", false, "upload current dir: ./")
	path := flag.String("p", "", "upload file path")

	flag.Parse()

	if *init {
		log.Println("init config")
		InitConfig()

		return
	}

	b := xgithub.New()

	if *allDir {
		filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
			// ignore dir
			if info.IsDir() {
				return nil
			}

			// ignore hidden file
			if strings.HasPrefix(info.Name(), ".") || strings.HasPrefix(path, ".") {
				return nil
			}

			log.Println(path)

			return nil
		})

		return
	}

	if *path != "" {
		p, err := b.Upload(*path)
		if err != nil {
			panic(err)
		}

		log.Printf("upload url: %s", p)
		return
	}

	log.Println("use -h to get help")
}
