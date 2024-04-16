package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Rsvp struct {
	EventId int    `json:"event_id"`
	Email   string `json:"email"`
}

func rsvpEvent(res http.ResponseWriter, req *http.Request) {
	var rsvp Rsvp
	userID := req.Header.Get("User-ID")
	err := json.NewDecoder(req.Body).Decode(&rsvp)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	sqlFile, err := os.ReadFile("tables/rsvp.sql")
	if err != nil {
		log.Fatal(err)
	}
}
