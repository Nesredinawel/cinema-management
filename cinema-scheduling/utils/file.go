package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveUploadedFile saves the uploaded file to the given directory and returns the relative path.
func SaveUploadedFile(file *multipart.FileHeader, folder string) (string, error) {
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create folder: %w", err)
	}

	savePath := filepath.Join(folder, file.Filename)
	if err := saveFile(file, savePath); err != nil {
		return "", err
	}

	return savePath, nil
}

// helper to save the file

func saveFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()
	if _, err = io.Copy(out, src); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}
