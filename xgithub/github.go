package xgithub

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"pbed/cons"
	"time"

	"github.com/google/go-github/v36/github"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
)

type (
	githubConfig struct {
		Repo  string `json:"repo"`
		Owner string `json:"owner"`
		// github write token
		Token  string `json:"token"`
		Branch string `json:"branch"`
		// default is time format, eg: 2015/01/02/
		PrefixFormat string `json:"prefixFormat"`
	}

	GithubBed struct {
		config githubConfig
		client *github.Client
	}
)

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

	conf := githubConfig{}
	err = json.Unmarshal(fs, &conf)
	if err != nil {
		panic(err)
	}

	if conf.PrefixFormat == "" {
		conf.PrefixFormat = time.Now().Format("2006/01/02/")
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

func (g *GithubBed) Upload(fname string) (string, error) {
	fi, err := os.Stat(fname)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(fname, os.O_RDONLY, 0755)
	if err != nil {
		return "", err
	}

	fs, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	conf := g.config

	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("upload by pic bed at " + time.Now().Format("2006-01-02 15:04:05")),
		Content:   fs,
		Branch:    github.String(conf.Branch),
		Committer: &github.CommitAuthor{Name: github.String(conf.Owner), Email: github.String("kazma233@outlook.com")},
	}

	resp, _, err := g.client.Repositories.CreateFile(context.Background(), conf.Owner, conf.Repo, conf.PrefixFormat+fi.Name(), opts)
	if err != nil {
		return "", err
	}

	return resp.Content.GetDownloadURL(), nil
}
