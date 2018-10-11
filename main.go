package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/marni/goigc"
	"net/http"
	"strconv"
	"time"
)

var startTime = time.Now()
var trackMAP = make(map[int]trackInfo)
var ids []int

// TODO: Make this print out in ISO 8061 standard
func uptime() time.Duration {
	return time.Since(startTime)
}

type apiInfo struct {
	Uptime  time.Duration `json:"uptime"`
	Info    string        `json:"info"`
	Version string        `json:"version"`
}

type trackInfo struct {
	Date     time.Time `json:"date"`
	Pilot    string    `json:"pilot"`
	Glider   string    `json:"glider"`
	GliderID string    `json:"glider_id"`
	Distance float64   `json:"distance"`
}

func newTrack(url string) {
	newTrack, err := igc.ParseLocation(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	track := trackInfo{
		newTrack.Date,
		newTrack.Pilot,
		newTrack.GliderType,
		newTrack.GliderID,
		newTrack.Task.Distance()}

	trackMAP[len(trackMAP)+1] = track
	ids = append(ids, len(ids)+1)
}

func handlerAPI(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	info := apiInfo{uptime(), "Service for IGC tracks.", "v1"}

	jsresp, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsresp)
}

func handlerIGC(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	switch r.Method {
	case "GET":
		if len(trackMAP) > 0 && len(ids) > 0 {
			jsresp, err := json.Marshal(ids)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(jsresp)
		} else {
			http.Error(w, "No files found", http.StatusNotFound)
			return
		}
	case "POST":
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func handlerID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	if len(trackMAP) > 0 && len(ids) > 0 {
		http.Header.Add(w.Header(), "content-type", "application/json")
		id, err := strconv.Atoi(param[0].Value)
		if err != nil {
			http.Error(w, "Please provide a valid ID.", http.StatusBadRequest)
			return
		} else {
			_, ok := trackMAP[id]
			if ok {
				fmt.Println(id)
				jsresp, err := json.Marshal(trackMAP[id])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(jsresp)
			} else {
				http.Error(w, "Given ID not found", http.StatusNotFound)
			}
		}
	} else {
		http.Error(w, "No files found", http.StatusNotFound)
	}
}

func main() {
	fmt.Println("Running...")
	router := httprouter.New()
	newTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	newTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc")
	newTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc")
	/*http.HandleFunc("/igcinfo/api", handlerAPI)
	http.HandleFunc("/igcinfo/api/igc", handlerIGC)
	http.HandleFunc("/igcinfo/api/igc/:id", handlerID)*/
	router.GET("/igcinfo/api", handlerAPI)
	router.GET("/igcinfo/api/igc", handlerIGC)
	router.GET("/igcinfo/api/igc/:id", handlerID)

	http.ListenAndServe("127.0.0.1:8080", router)
}
