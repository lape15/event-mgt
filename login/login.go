package login

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"event-mgt/database"
	"event-mgt/shared"
	"event-mgt/signup"

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
	shared.Credential
	Id int
}

func scanDb(email, userName string) (CredentialId, error) {
	var user CredentialId

	row, err := database.Db.QueryRow("SELECT email, username, password,user_id FROM users WHERE email = ? or username = ? LIMIT 1", email, userName)
	row.Scan(&user.Email, &user.Username, &user.Password, &user.Id)
	if err != nil {
		log.Fatal(err)
		return CredentialId{}, err
	}
	return user, nil
}

func Login(res http.ResponseWriter, req *http.Request) {
	var user shared.Credential
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
		tokenString, err := signup.CreateToken(user.Username)
		if err != nil {
			fmt.Print(err.Error())
		}
		response := shared.UserCache{
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
