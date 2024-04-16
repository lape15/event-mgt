package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func editEvent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var editedEvent Event
	eventId, _ := strconv.Atoi(vars["id"])
	userId := req.Header.Get("User-ID")
	id, err := strconv.Atoi(userId)
	query := `
    UPDATE events
    SET name = ?, description = ?, start = ?, end = ?, location = ?, event_limit = ?
    WHERE event_id = ? AND organizer_id = ?
    `
	if err != nil {
		fmt.Println("Error:", err)
	}
	errs := json.NewDecoder(req.Body).Decode(&editedEvent)
	if errs != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := db.Exec(query, editedEvent.Name, editedEvent.Description, editedEvent.Start, editedEvent.End, editedEvent.Location, editedEvent.EventLimit, eventId, id)
	if err != nil {
		fmt.Println("Error:", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error:", err)
	}

	if rowsAffected == 0 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Event does not belong to user!"))
	} else {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("ok!"))
	}
}
