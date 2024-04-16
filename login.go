package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func comparePassword(hashed, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		}
	}
	return true, nil
}

type CredentialId struct {
	Credential
	Id int
}

func scanDb(email, userName string) (CredentialId, error) {
	var user CredentialId
	// err := db.QueryRow("SELECT email, username, password,user_id  FROM (SELECT * FROM users WHERE email = ? UNION SELECT * FROM users WHERE username = ?) AS combined_results LIMIT 1", email, userName).Scan(&user.Email, &user.Username, &user.Password, &user.Id)
	err := db.QueryRow("SELECT email, username, password,user_id FROM users WHERE email = ? or username = ? LIMIT 1", email, userName).Scan(&user.Email, &user.Username, &user.Password, &user.Id)
	if err != nil {
		log.Fatal(err)
		return CredentialId{}, err
	}
	return user, nil
}

func login(res http.ResponseWriter, req *http.Request) {
	var user Credential
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if user.Email == "" || user.Password == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Username and password are required"))
		return
	}

	result, _ := scanDb(user.Email, user.Username)
	valid, passErr := comparePassword(result.Password, user.Password)

	if passErr != nil {
		fmt.Errorf(err.Error())
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte("Invalid password!"))
		return
	}
	if valid {
		tokenString, err := createToken(user.Username)
		if err != nil {
			fmt.Print(err.Error())
		}
		response := struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Token    string `json:"token"`
			Id       int    `json:"id"`
		}{
			Email:    result.Email,
			Username: result.Username,
			Token:    tokenString,
			Id:       result.Id,
		}
		responseJson, err := json.Marshal(response)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(responseJson)

	}

}
