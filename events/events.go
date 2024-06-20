package events

import (
	"encoding/json"
	"event-mgt/database"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type QueriedEvent struct {
	EventId     int       `json:"event_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Start       string    `json:"start"`
	End         string    `json:"end"`
	Location    string    `json:"location"`
	OrganizerID int       `json:"organizer_id" db:"organizer_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	EventLimit  int       `json:"event_limit"`
}

func getEvents(id int) ([]QueriedEvent, error) {
	var event QueriedEvent
	var events []QueriedEvent
	query := "SELECT event_id, name, description, start, end, location, event_limit FROM events WHERE organizer_id = ?"

	rows, err := database.Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&event.EventId, &event.Name, &event.Description, &event.Start, &event.End, &event.Location, &event.EventLimit)
		if err != nil {
			return nil, err
		}
		events = append(events, event)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func GetUserEvents(res http.ResponseWriter, req *http.Request) {

	userId := req.Header.Get("User-ID")
	num, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Println("Error:", err)
	}
	events, err := getEvents(num)
	if err != nil {
		panic(err.Error())
	}

	response, err := json.Marshal(events)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	if len(events) == 0 {
		emptyArray := []byte("[]")
		res.Write(emptyArray)
	} else {
		res.Write(response)
	}
}
