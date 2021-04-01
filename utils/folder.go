package utils

import (
	"os"
	"path/filepath"
)

func WalkFolderFiles(folderPath string) (filePaths []string, err error) {
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	return
}
