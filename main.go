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

// var db *sql.DB

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", login.Login).Methods("POST")
	r.HandleFunc("/signup", signup.Signup).Methods("POST")
	r.HandleFunc("/events", protected.ProtectedHandler(handleEventRequest))
	r.HandleFunc("/events/{id}", singleEventHandler)
	r.HandleFunc("/event/rsvp", rsvpEvent)
	// r.HandleFunc("/", handleRequest)
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }
	// cfg := mysql.Config{
	// 	User:                 os.Getenv("DB_USER"),
	// 	Passwd:               os.Getenv("DB_PASS"),
	// 	Net:                  "tcp",
	// 	Addr:                 "sql11.freesqldatabase.com",
	// 	DBName:               os.Getenv("DB_NAME"),
	// 	AllowNativePasswords: true,
	// }

	// var err error
	// db, err = sql.Open("mysql", cfg.FormatDSN())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pingErr := db.Ping()
	// if pingErr != nil {
	// 	log.Fatal(pingErr)
	// }
	// fmt.Println("Connected!")
	database.ConnectDb()
	http.ListenAndServe(":8000", r)
	defer database.Db.CloseDb()
}

func handleEventRequest(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		events.CreateEvent(res, req)
	case http.MethodGet:
		getUserEvents(res, req)
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
