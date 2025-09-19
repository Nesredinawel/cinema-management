package utils

import (
	"mime/multipart"
	"net/url"
	"path/filepath"
)

// SaveMoviePoster handles either a file upload or URL
func SaveMoviePoster(file *multipart.FileHeader, posterURL string, saveFolder string, saveFunc func(*multipart.FileHeader, string) error) (*string, error) {
	if file != nil {
		savePath := filepath.Join(saveFolder, file.Filename)
		if err := saveFunc(file, savePath); err != nil {
			return nil, err
		}
		return &savePath, nil
	} else if posterURL != "" {
		return &posterURL, nil
	}
	return nil, nil
}

// ConvertPosterToPublicURL converts local path to a full public URL
func ConvertPosterToPublicURL(poster *string, baseURL string) string {
	if poster == nil || *poster == "" {
		return ""
	}
	u, err := url.JoinPath(baseURL, *poster)
	if err != nil {
		return *poster
	}
	return u
}
