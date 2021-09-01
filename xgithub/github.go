package xgithub

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"pbed/cons"
	"time"

	"github.com/google/go-github/v36/github"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
)

type (
	Author struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	GithubConfig struct {
		Repo  string `json:"repo"`
		Owner string `json:"owner"`
		// github write token
		Token  string `json:"token"`
		Branch string `json:"branch"`
		// default is time format, eg: 2015/01/02/
		PrefixFormat  string `json:"prefixFormat"`
		MessageFormat string `json:"messageFormat"`
		Author        Author `json:"author"`
	}

	GithubBed struct {
		config GithubConfig
		client *github.Client
	}
)

// New get github file operate obj
func New() *GithubBed {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(home+cons.ConfigPath, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	conf := GithubConfig{}
	err = json.Unmarshal(fs, &conf)
	if err != nil {
		panic(err)
	}

	if conf.PrefixFormat == "" {
		conf.PrefixFormat = time.Now().Format("2006/01/02/")
	}

	if conf.MessageFormat == "" {
		conf.MessageFormat = "upload by pic bed at: %s"
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &GithubBed{
		config: conf,
		client: client,
	}
}

// UploadByPath impl bed.Bed
func (g *GithubBed) UploadByPath(filePath string) (string, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		return "", err
	}

	fs, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return g.UploadByBytes(fs, fi.Name())
}

func (g *GithubBed) UploadByBytes(bs []byte, fileName string) (string, error) {
	conf := g.config

	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(fmt.Sprintf(conf.MessageFormat, time.Now().Format("2006-01-02 15:04:05"))),
		Content:   bs,
		Branch:    github.String(conf.Branch),
		Committer: &github.CommitAuthor{Name: github.String(conf.Author.Name), Email: github.String(conf.Author.Email)},
	}

	resp, _, err := g.client.Repositories.CreateFile(context.Background(), conf.Owner, conf.Repo, conf.PrefixFormat+fileName, opts)
	if err != nil {
		return "", err
	}

	return resp.Content.GetDownloadURL(), nil
}
