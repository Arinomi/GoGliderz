package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/marni/goigc"
	"net/http"
	"time"
)

// global variable init
var startTime = time.Now()
var trackMAP = make(map[int]trackInfo)
var ids []int

// struct init
type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

type trackInfo struct {
	Date     time.Time `json:"date"`
	Pilot    string    `json:"pilot"`
	Glider   string    `json:"glider"`
	GliderID string    `json:"glider_id"`
	Distance float64   `json:"distance"`
}

// returning uptime as a string in ISO 8601/RFC3339 format
func uptime() string {
	now := time.Now()
	now.Format(time.RFC3339)
	startTime.Format(time.RFC3339)

	return now.Sub(startTime).String()
}

// creates a new track from the presented url and returns its ID
// returns 0 if the url was invalid
func newTrack(url string) int {
	newTrack, err := igc.ParseLocation(url)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	track := trackInfo{
		newTrack.Date,
		newTrack.Pilot,
		newTrack.GliderType,
		newTrack.GliderID,
		newTrack.Task.Distance()}

	trackMAP[len(trackMAP)+1] = track
	ids = append(ids, len(ids)+1)
	return len(ids)
}

// main function
func main() {
	// router init
	router := httprouter.New()

	// mock data
	_ = newTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	_ = newTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc")
	_ = newTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc")

	// routes init
	router.GET("/igcinfo/api", handlerAPI)
	router.GET("/igcinfo/api/igc", handlerIGC)
	router.GET("/igcinfo/api/igc/:id", handlerID)
	router.POST("/igcinfo/api/igc", handlerIGC)
	router.GET("/igcinfo/api/igc/:id/:field", handlerField)

	// server init
	http.ListenAndServe("127.0.0.1:8080", router)
}
