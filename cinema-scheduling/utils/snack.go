package utils

import (
	"mime/multipart"
	"net/url"
	"path/filepath"
)

// SaveSnackImage handles either a file upload or URL and returns a pointer to string
func SaveSnackImage(file *multipart.FileHeader, imageURL string, saveFolder string, saveFunc func(*multipart.FileHeader, string) error) (*string, error) {
	if file != nil {
		savePath := filepath.Join(saveFolder, file.Filename)
		if err := saveFunc(file, savePath); err != nil {
			return nil, err
		}
		return &savePath, nil
	} else if imageURL != "" {
		return &imageURL, nil
	}
	return nil, nil
}

// ConvertSnackImageToPublicURL converts local path to full public URL and returns pointer
func ConvertSnackImageToPublicURL(image *string, baseURL string) *string {
	if image == nil || *image == "" {
		return nil
	}
	u, err := url.JoinPath(baseURL, *image)
	if err != nil {
		return image
	}
	return &u
}
