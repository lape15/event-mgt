package events

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"event-mgt/database"

	"github.com/gorilla/mux"
)

type AEvent struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
	EventLimit  int       `json:"event_limit"`
}

func GetEvent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var event AEvent
	var start, end []byte
	eventId, _ := strconv.Atoi(vars["id"])

	row, err := database.Db.QueryRow("SELECT name, description, start, end, location, event_limit FROM events WHERE event_id = ?", eventId)
	row.Scan(
		&event.Name, &event.Description, &start, &end, &event.Location, &event.EventLimit)
	if err != nil {
		log.Fatal(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	event.Start, err = time.Parse("2006-01-02 15:04:05", string(start))
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	event.End, err = time.Parse("2006-01-02 15:04:05", string(end))
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseJson, err := json.Marshal(event)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(responseJson)
}
