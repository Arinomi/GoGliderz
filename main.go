package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var startTime = time.Now()

// TODO: Make this print out in ISO 8061 standard
func uptime() time.Duration {
	return time.Since(startTime)
}

type apiInfo struct {
	Uptime  time.Duration `json:"uptime"`
	Info    string        `json:"info"`
	Version string        `json:"version"`
}

func handlerAPI(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	info := apiInfo{uptime(), "Service for IGC tracks.", "v1"}

	jsresp, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsresp)
}

func main() {
	http.HandleFunc("/api/", handlerAPI)
	http.ListenAndServe(":8080", nil)
}
