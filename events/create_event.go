package events

import (
	"encoding/json"
	"event-mgt/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Event struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
	OrganizerID int       `json:"organizer_id" db:"organizer_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	EventLimit  int       `json:"event_limit"`
}

type Res struct {
	Message string `json:"message"`
	Event   Event  `json:"event"`
}

func CreateEvent(res http.ResponseWriter, req *http.Request) {
	var event Event
	userID := req.Header.Get("User-ID")
	err := json.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	sqlFile, err := os.ReadFile("tables/event.sql")
	if err != nil {
		log.Fatal(err)
	}
	q, err := database.Db.Exec(string(sqlFile))
	if err != nil {
		panic(err.Error())
	}
	num, errn := strconv.Atoi(userID)
	if errn != nil {
		fmt.Println("Error:", errn)
	}
	event.OrganizerID = num
	start := event.Start.Format("2006-01-02 15:04:05")
	end := event.End.Format("2006-01-02 15:04:05")

	insert, err := database.Db.Exec("INSERT INTO events(name,description,start,end,location,event_limit,organizer_id) VALUES (?, ?, ?,?, ?, ?,?)", event.Name, event.Description, start, end, event.Location, event.EventLimit, event.OrganizerID)
	if err != nil {
		panic(err.Error())
	}
	// defer insert.Close()
	// rowsAffected, errs := q.RowsAffected()
	// if err != nil {
	// 	panic(err.Error())
	// }
	rowsAffected, errs := insert.RowsAffected()
	if errs != nil {
		fmt.Print("Error here")
	}
	fmt.Printf("Rows affected: %d\n", rowsAffected)
	response := Res{
		Message: "Created!",
		Event:   event,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Print(q)
	fmt.Printf("%v%s\n", event, userID)
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonResponse)

}
