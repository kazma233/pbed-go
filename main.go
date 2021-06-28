package main

import (
	"log"
	"os"
	"pbed/cons"
	"pbed/xgithub"

	"github.com/mitchellh/go-homedir"
)

func init() {
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
	b := xgithub.New()

	p, err := b.Upload("any file")
	if err != nil {
		panic(err)
	}

	log.Println(p)
}
