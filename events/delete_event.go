package events

import (
	"event-mgt/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteEvent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	eventId, _ := strconv.Atoi(vars["id"])
	userId := req.Header.Get("User-ID")
	id, errs := strconv.Atoi(userId)
	query := "DELETE FROM events WHERE event_id = ? AND organizer_id = ?"
	result, err := database.Db.Exec(query, eventId, id)
	if errs != nil {
		fmt.Println("Error:", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error:", err)
	}

	if rowsAffected == 0 {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Event not found or not owned by the user"))
	} else {
		res.Write([]byte("Deleted!"))
		res.WriteHeader(http.StatusNoContent)
		fmt.Print("Deleted!")
	}
}
