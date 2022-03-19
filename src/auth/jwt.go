package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/redsuperbat/kafka-commander/src/server"
)

type User struct {
	Username string
}

const bearer = "Bearer "

func NewJwtValidator(pubKey string) func(string) (*User, *server.ResponseError) {
	return func(tokenString string) (*User, *server.ResponseError) {
		if tokenString == "" {
			return nil, server.NewRespErr(http.StatusUnauthorized)
		}

		if strings.HasPrefix(tokenString, bearer) {
			tokenString = tokenString[len(bearer):]
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return pubKey, nil
		})

		if err != nil {
			log.Println(err.Error())
			return nil, server.NewRespErr(http.StatusUnauthorized)
		}

		if !token.Valid {
			log.Println("Invalid token tried accessing the API")
			return nil, server.NewRespErr(http.StatusUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, server.NewRespErr(http.StatusUnauthorized)
		}

		username, uOk := claims["username"].(string)
		if !uOk {
			return nil, server.NewRespErr(http.StatusUnauthorized)
		}

		return &User{Username: username}, nil
	}
}
