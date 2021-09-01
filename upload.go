package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"pbed/bed"
	"strings"
)

// uploadDir upload dir file
func uploadDir(b bed.Bed, baseDirPath string) (res []string) {
	filepath.Walk(baseDirPath, func(path string, info fs.FileInfo, err error) error {
		// ignore dir
		if info.IsDir() {
			return nil
		}

		// ignore hidden file
		if strings.HasPrefix(info.Name(), ".") || strings.HasPrefix(path, ".") {
			return nil
		}

		res = append(res, upload(b, path))

		return nil
	})

	return
}

// upload one file
func upload(b bed.Bed, filePath string) (url string) {
	url, err := b.UploadByPath(filePath)
	if err != nil {
		log.Printf("[%s]: Upload failed, reason -> %v", filePath, err)
	} else {
		log.Printf("[%s]: Upload finished, url -> %s", filePath, url)
	}

	return
}
