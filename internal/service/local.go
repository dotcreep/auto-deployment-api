package service

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

func CopyFile(src, dst string) error {
	dstDir := filepath.Dir(dst)
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
			return err
		}
	}
	for {
		if _, err := os.Stat(src); os.IsNotExist(err) {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}

func CreateDir(path, name, area string) error {
	fullPath := filepath.Join(path, name, area)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
