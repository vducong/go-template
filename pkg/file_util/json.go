package fileutil

import (
	"io"
	"os"
	"promotion/pkg/failure"
)

func JSONReader(filePath string) ([]byte, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	defer jsonFile.Close()

	byteData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	return byteData, nil
}
