package utils

import (
	"os"
	"path/filepath"
)

func GetImagePaths(dir string) []string {
	_, err := os.Stat(dir)
	if err != nil {
		println("Warning: No such directory ", dir)
		os.Exit(1)
	}

	isDir, err := isDirectory(dir)
	if err != nil || !isDir {
		println("Warning: No such directory ", dir)
		os.Exit(1)
	}

	var imgs []string

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg") {
			imgs = append(imgs, path)
		}

		return nil
	})

	return imgs
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
