package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileType string

const (
	JPEG FileType = "JPEG"
	PNG  FileType = "PNG"
)

func GetNewUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func GetFileType(filename string) (FileType, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpeg", ".jpg":
		return JPEG, nil
	case ".png":
		return PNG, nil
	default:
		return "", errors.New("unsupported file type")
	}
}

func GenerateKeyForS3(filename string) (string, error) {
	fileType, err := GetFileType(filename)
	if err != nil {
		return "", err
	}
	uuidKey := GetNewUUID()

	key := fmt.Sprintf("%s/%s/%s", fileType, filename, uuidKey)

	return key, nil
}
