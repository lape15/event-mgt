package test

import (
	"bytes"
	"encoding/json"
	"event-mgt/database"
	"event-mgt/shared"
	"event-mgt/signup"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type MockHasher struct{}

func (m MockHasher) HashPassword(password string) string {
	return "mockHashedPassword"
}

func TestSignUp(t *testing.T) {
	var credential = shared.Credential{
		Email:    "star@gmail.com",
		Password: "password",
		Username: "star",
	}

	reqBody, err := json.Marshal(credential)
	if err != nil {
		t.Fatalf("Error marshaling user details: %v", err)
	}
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\?").
		WithArgs(credential.Email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Expect the table creation query to be executed
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(0, 0))

	var oldPasswordHasher = signup.NewPasswordHasher
	signup.NewPasswordHasher = MockHasher{}
	defer func() { signup.NewPasswordHasher = oldPasswordHasher }()

	credential.Password = signup.NewPasswordHasher.HashPassword(credential.Password)

	mock.ExpectExec("INSERT INTO users \\(username, email, password\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs(credential.Username, credential.Email, credential.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tokenString, err := signup.CreateToken(credential.Username)
	if err != nil {
		panic(err.Error())
	}

	oldDB := database.Db
	database.Db.SetDB(mockDB)
	defer func() { database.Db = oldDB }()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		signup.Signup(res, req, "../tables/users.sql")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	response := shared.UserCache{
		Email:    credential.Email,
		Username: credential.Username,
		Token:    tokenString,
	}
	expectedJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Error marshaling expected response: %v", err)
	}

	if rr.Body.String() != string(expectedJSON) {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedJSON))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
