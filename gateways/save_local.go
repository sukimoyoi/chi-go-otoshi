package gateways

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type SaveLocalRepository struct{}

func (slr *SaveLocalRepository) CreateTitleFolder(folderPath string) error {
	return os.MkdirAll(folderPath, os.ModePerm)
}

func (slr *SaveLocalRepository) CreateNumberedFolder(folderParentPath string) (folderPath string, err error) {
	number := 1
	for {
		folderPath = filepath.Join(folderParentPath, fmt.Sprintf("%02d", number))
		if err = os.Mkdir(folderPath, os.ModePerm); err != nil {
			if !os.IsExist(err) {
				return "", err
			}
		} else {
			break
		}
		number++
	}
	return folderPath, nil
}

func (slr *SaveLocalRepository) SaveFileReader(filePath string, r io.Reader) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	return err
}

func (slr *SaveLocalRepository) ReadFile(filePath, str string) ([]byte, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
