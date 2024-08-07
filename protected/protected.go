package protected

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func verifyToken(token string) error {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !tkn.Valid {
		return fmt.Errorf("token is currently invalid")
	}
	return nil
}

func ProtectedHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Cntent-Type", "application/json")
		token := req.Header.Get("Authorization")
		if token == "" {
			res.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(res, "you are not authorized!")
			return
		}
		token = token[len("Bearer "):]
		err := verifyToken(token)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(res, err.Error())
			return
		}
		next(res, req)
	}
}
