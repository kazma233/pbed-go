package main

import (
	"flag"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pbed/bed"
	"pbed/cons"
	"pbed/xgithub"
	"strings"

	"github.com/mitchellh/go-homedir"
)

var (
	uploadDirPath  string
	uploadFilePath string
	serverAddr     string
)

func init() {
	flag.StringVar(&uploadDirPath, "d", "", "upload dir files, example: d ./")
	flag.StringVar(&uploadFilePath, "p", "", "upload file path: example: p ./a.txt")
	flag.StringVar(&serverAddr, "server", "", "start a pic bed server, example: server :9000")

	flag.Parse()

	initConfig()
}

func main() {
	b := xgithub.New()

	// start a server
	if serverAddr != "" {
		startServer(b, serverAddr)

		return
	}

	// upload dir
	if uploadDirPath != "" {
		uploadDir(b, uploadDirPath)

		return
	}

	// upload sigle file
	if uploadFilePath != "" {
		upload(b, uploadFilePath)

		return
	}

	log.Println("use -h to get help")
}

func startServer(b bed.Bed, addr string) {
	// index
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("ui/index.html"))
		err := tpl.Execute(rw, "")
		if err != nil {
			log.Printf("write template error: %v", err)
		}
	})
	// load file
	http.HandleFunc("/upload", func(rw http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("ui/index.html"))

		bs, fileName, err := readRequestFile(r)
		if err != nil {
			log.Printf("read file error: %v", err)
			_ = tpl.Execute(rw, err.Error())
			return
		}

		urlPath, err := b.UploadByBytes(bs, fileName)
		if err != nil {
			log.Printf("upload error: %v", err)
			_ = tpl.Execute(rw, err.Error())
			return
		}

		err = tpl.Execute(rw, urlPath)
		if err != nil {
			log.Printf("write template error: %v", err)
		}
	})

	log.Println(http.ListenAndServe(addr, nil))
}

func readRequestFile(r *http.Request) (bs []byte, fileName string, err error) {
	err = r.ParseMultipartForm(1024 * 1024 * 15)
	if err != nil {
		return
	}

	f, h, err := r.FormFile("file")
	if err != nil {
		return
	}

	fileName = h.Filename

	bs, err = ioutil.ReadAll(f)

	return
}

func uploadDir(b bed.Bed, baseDirPath string) {
	filepath.Walk(baseDirPath, func(path string, info fs.FileInfo, err error) error {
		// ignore dir
		if info.IsDir() {
			return nil
		}

		// ignore hidden file
		if strings.HasPrefix(info.Name(), ".") || strings.HasPrefix(path, ".") {
			return nil
		}

		upload(b, path)

		return nil
	})
}

// upload upload one file
func upload(b bed.Bed, filePath string) (url string) {
	url, err := b.UploadByPath(filePath)
	if err != nil {
		log.Printf("[%s]: Upload finish, url -> %s", filePath, url)
	} else {
		log.Printf("[%s]: Upload failed, reason -> %v", filePath, err)
	}

	return
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
