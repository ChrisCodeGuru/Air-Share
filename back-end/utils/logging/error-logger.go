package logging

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
)

func ErrorLogger(errorMsg error, r *http.Request) {
	t := time.Now().UTC()
	dateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	logDetails := map[string]interface{}{
		"URL":    	r.URL.String(),
		"Error":	errorMsg.Error(),
	}

	logDetailsJSON, err := json.Marshal(logDetails)
	if err != nil {
		panic(err)
	}
    
    // Add logs to file
    file, err := os.OpenFile("./log/error-log.log", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()
	_, writeErr := file.WriteString(dateTime + " " + string(logDetailsJSON) + " \n")
    if writeErr != nil {
        fmt.Println(err)
    }
}