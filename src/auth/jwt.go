package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/redsuperbat/kafka-commander/src/server"
	"github.com/rs/zerolog/log"
)

const bearer = "Bearer "

func NewJwtMiddleware(pubKey *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" {
			server.SendDefaultErr(c, http.StatusUnauthorized)
			c.Abort()
			return
		}

		if strings.HasPrefix(tokenString, bearer) {
			tokenString = tokenString[len(bearer):]
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return pubKey, nil
		})

		if err != nil {
			log.Info().Msg(err.Error())
			server.SendDefaultErr(c, http.StatusUnauthorized)
			c.Abort()
			return
		}

		if !token.Valid {
			log.Info().Msg("Invalid token tried accessing the API")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			server.SendDefaultErr(c, http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
