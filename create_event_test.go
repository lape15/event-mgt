package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateEvent(t *testing.T) {
	// Create a request body
	event := Event{
		Name:        "Test Event",
		Description: "This is a test event",
		Start:       time.Now(),
		End:         time.Now().Add(1 * time.Hour),
		Location:    "Test Location",
		EventLimit:  100,
	}

	reqBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Error marshaling event: %v", err)
	}

	req, err := http.NewRequest("POST", "/events", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Set user ID in the header
	req.Header.Set("User-ID", "1")

	// Mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS events").WillReturnResult(sqlmock.NewResult(0, 0))

	start := event.Start.Format("2006-01-02 15:04:05")
	end := event.End.Format("2006-01-02 15:04:05")

	mock.ExpectExec("INSERT INTO events").WithArgs(event.Name, event.Description, start, end, event.Location, event.EventLimit, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	// Replace global db variable with mockDB
	oldDB := db
	db = mockDB
	defer func() { db = oldDB }()

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEvent)

	// Serve the HTTP request to the ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Format the expected response
	expected := Res{
		Message: "Created!",
		Event: Event{
			Name:        event.Name,
			Description: event.Description,
			Start:       event.Start,
			End:         event.End,
			Location:    event.Location,
			EventLimit:  event.EventLimit,
			OrganizerID: 1,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Error marshaling expected response: %v", err)
	}

	// Check the response body
	if rr.Body.String() != string(expectedJSON) {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedJSON))
	}

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
