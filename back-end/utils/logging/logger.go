package logging

import (
	"fmt"
	"net/http"
	"time"
	"os"
	"encoding/json"
)

func Logger(formData string, status int, r *http.Request) {
	t := time.Now().UTC()
	dateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	logDetails := map[string]interface{}{
		"Method":    	r.Method,
		"URL":       	r.URL.String(),
		"Protocol":  	r.Proto,
		"User Agent":	r.Header.Get("User-Agent"),
		"IP Address": 	r.Header.Get("X-Real-Ip"),
		"Status Code":	status,
		"Form Input":	formData,
	}

	logDetailsJSON, err := json.Marshal(logDetails)
	if err != nil {
		panic(err)
	}
    
    // Add logs to file
    file, err := os.OpenFile("./log/general-log.log", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()
	_, writeErr := file.WriteString(dateTime + " " + string(logDetailsJSON) + " \n")
    if err != nil {
        fmt.Println(writeErr)
    }
}
