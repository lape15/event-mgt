package main

import (
	"fmt"
	"net/http"

	"event-mgt/database"
	"event-mgt/events"
	"event-mgt/login"
	"event-mgt/protected"
	"event-mgt/signup"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", login.Login).Methods("POST")
	r.HandleFunc("/signup", signup.Signup).Methods("POST")
	r.HandleFunc("/events", protected.ProtectedHandler(handleEventRequest))
	r.HandleFunc("/events/{id}", singleEventHandler)
	r.HandleFunc("/event/rsvp", events.RsvpEvent)
	r.HandleFunc("/", handleRequest)

	database.ConnectDb()
	http.ListenAndServe(":8000", r)
	defer database.Db.CloseDb()
}

func handleEventRequest(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		events.CreateEvent(res, req, "tables/event.sql")
	case http.MethodGet:
		events.GetUserEvents(res, req)
	default:
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func singleEventHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPatch:
		events.EditEvent(res, req)
	case http.MethodDelete:
		events.DeleteEvent(res, req)
	case http.MethodGet:
		events.GetEvent(res, req)
	default:
		fmt.Print(http.MethodPatch)
		http.Error(res, "Method not allowed there", http.StatusMethodNotAllowed)
	}
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	default:
		fmt.Print("Default calee\n")
		fmt.Print(http.MethodPatch)
		http.Error(res, "Method not allowed there", http.StatusMethodNotAllowed)
	}
}
