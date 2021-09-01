package main

import (
	"encoding/json"
	"log"
	"os"
	"pbed/cons"
	"pbed/xgithub"

	"github.com/mitchellh/go-homedir"
)

var (
	BedType = []string{"github"}
)

type PbedConfig struct {
	Github xgithub.GithubConfig `json:"github"`
}

func init() {
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

		bs, err := json.MarshalIndent(PbedConfig{}, "", "\t")
		if err != nil {
			panic(err)
		}

		_, err = f.Write(bs)

		if err != nil {
			panic(err)
		}
	}
}
