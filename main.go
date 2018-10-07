package main

import (
	"fmt"
	"net/http"
	"time"
)

var startTime = time.Now()

// TODO: Make this print out in ISO 8061 standard
func uptime() time.Duration {
	return time.Since(startTime)
}

/*type apiInfo struct {
	Uptime		time.Duration 	`json:"uptime"`
	Info		string 			`json:"info"`
	Version		string 			`json:"version"`
}*/

func handlerAPI(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	fmt.Fprintln(w, "Uptime:", uptime())
}

func main() {
	http.HandleFunc("/api/", handlerAPI)
	http.ListenAndServe(":8080", nil)
}
