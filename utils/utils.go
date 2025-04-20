package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"parse_photo_go/models"
	"time"
)

func ResultMaker(message string) models.Result {
	return models.Result{Result: "ok", Message: message}
}

func DateConvert(date ...string) string {
	if len(date) > 0 && date[0] != "" {
		parsedDate, err := time.Parse("2006-01-02 15:04:05", date[0])
		if err != nil {
			return time.Now().Format("2006-01-02 15:04:05")
		}
		return parsedDate.Format("2006-01-02 15:04:05")
	}
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetHashByPath(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %s", err.Error())
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %s", err.Error())
	}

	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
}
