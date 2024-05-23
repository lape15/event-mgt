package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.Claims
}

type ErrorResponse struct {
	Err string
}

type Error interface {
	Error() string
}

type PasswordHasher interface {
	HashPassword(password string) string
}

var passwordHasher PasswordHasher = DefaultHasher{}

type DefaultHasher struct{}

func (h DefaultHasher) HashPassword(password string) string {
	return generatePasswordHash(password)
}

var secretKey = []byte("secret-key")

func generatePasswordHash(password string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		err := &ErrorResponse{
			Err: "Error decoding string",
		}
		log.Fatal(err)
	}
	return string(pass)
}

func createToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func signup(res http.ResponseWriter, req *http.Request) {

	var credential Credential
	err := json.NewDecoder(req.Body).Decode(&credential)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var emailCount int

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", credential.Email).Scan(&emailCount)

	if err != nil {
		fmt.Println("Error querying database:", err)
		return
	}
	fmt.Print(emailCount)
	if emailCount > 0 {
		fmt.Println("Email already exists")
		res.Write([]byte("Email already exists"))
		return
	}

	credential.Password = passwordHasher.HashPassword(credential.Password)
	sqlFile, err := os.ReadFile("tables/users.sql")

	if err != nil {
		log.Fatal(err)
	}
	q, err := db.Exec(string(sqlFile))
	if err != nil {
		panic(err.Error())
	}
	insert, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", credential.Username, credential.Email, credential.Password)
	if err != nil {
		panic(err.Error())
	}
	// defer insert.Close()
	// rowsAffected, errs := q.RowsAffected()
	rowsAffected, errs := insert.RowsAffected()
	if errs != nil {
		fmt.Print("Error here")
	}

	fmt.Printf("Rows affected: %d\n", rowsAffected)
	tokenString, err := createToken(credential.Username)

	if err != nil {

		panic(err.Error())
	}
	response := UserCache{
		Email:    credential.Email,
		Username: credential.Username,
		Token:    tokenString,
		// Id:       credential.Id,
	}
	fmt.Print(q)
	responseJson, err := json.Marshal(response)
	c.update(response.Email, response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(responseJson)

}
