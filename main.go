package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/events", protectedHandler(handleEventRequest))
	r.HandleFunc("/events/{id}", singleEventHandler)
	r.HandleFunc("/event/rsvp", rsvpEvent)
	r.HandleFunc("/", handleRequest)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 "sql11.freesqldatabase.com",
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	http.ListenAndServe(":8000", r)
	defer db.Close()
}

func handleEventRequest(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		createEvent(res, req)
	case http.MethodGet:
		getUserEvents(res, req)
	default:
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func singleEventHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPatch:
		editEvent(res, req)
	case http.MethodDelete:
		deleteEvent(res, req)
	case http.MethodGet:
		getEvent(res, req)
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
