package helpers

import (
	"parse_photo_go/models"
	"time"
)

func ResultMaker(message string) models.Result {
	return models.Result{Result: "ok", Message: message}
}

func dateConvert(date ...string) string {
	if len(date) > 0 && date[0] != "" {
		parsedDate, err := time.Parse("2006-01-02 15:04:05", date[0])
		if err != nil {
			return time.Now().Format("2006-01-02 15:04:05")
		}
		return parsedDate.Format("2006-01-02 15:04:05")
	}
	return time.Now().Format("2006-01-02 15:04:05")
}
