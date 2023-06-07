package fileutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"promotion/pkg/failure"
)

func CSVReadFromPath(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to open file=%s: %w", filePath, err))
	}
	defer f.Close()

	records, err := CSVRead(f)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to read file=%s: %w", filePath, err))
	}
	return records, nil
}

func CSVMultipartRead(fh *multipart.FileHeader) ([][]string, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to open file: %w", err))
	}
	defer f.Close()

	records, err := CSVRead(f)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to read file: %w", err))
	}
	return records, nil
}

func CSVRead(f io.Reader) ([][]string, error) {
	csvReader := csv.NewReader(f)
	csvReader.LazyQuotes = true
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	return records, nil
}

func CSVWrite(filePath string, codes [][]string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if errWrite := w.WriteAll(codes); errWrite != nil {
		return failure.ErrWithTrace(errWrite)
	}
	return nil
}
