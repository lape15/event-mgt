package events

import (
	"encoding/json"
	"event-mgt/database"
	"event-mgt/shared"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Rsvp struct {
	EventId int    `json:"event_id"`
	Email   string `json:"email"`
}

func RsvpEvent(res http.ResponseWriter, req *http.Request) {
	var rsvp Rsvp
	userID := req.Header.Get("User-ID")
	err := json.NewDecoder(req.Body).Decode(&rsvp)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if shared.IsUserHostOfEvent(userID, rsvp.EventId) {
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte("You cannot RSVP to your own event!"))
		return
	}

	var exists = shared.DoesUserExist(rsvp.Email)
	if !exists {
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte("User must have an account to rsvp an event!"))
		return
	}
	sqlFile, err := os.ReadFile("tables/rsvp.sql")
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
	rsvp.EventId = num
	insert, err := database.Db.Query("INSERT INTO event_rsvps(event_id,user_id)", rsvp.EventId, userID)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	rowsAffected, errs := q.RowsAffected()
	if errs != nil {
		fmt.Print("Error here")
	}
	fmt.Printf("Rows affected: %d\n", rowsAffected)
}
