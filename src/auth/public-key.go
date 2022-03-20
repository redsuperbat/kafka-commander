package auth

import (
	"crypto/rsa"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/redsuperbat/kafka-commander/src/options"
)

func GetPubKey() *rsa.PublicKey {

	key := os.Getenv("JWT_PUB_KEY")
	args := options.GetArgs()

	if args.PubKeyUrl == "" && key == "" {
		log.Fatalln("Please specify a JWT pub key either by using the JWT_PUB_KEY env variable or a url with the --pub-key-url flag")
	}

	if key != "" {
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
		if err != nil {
			log.Fatalln("Invalid RSA public key", err.Error())
		}
		return pubKey
	}

	resp, err := http.Get(args.PubKeyUrl)
	if err != nil {
		log.Fatalln("No public key was found at", args.PubKeyUrl, err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Unable to read response body", err.Error())
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(body)
	if err != nil {
		log.Fatalln("Invalid public key in response body", err.Error())
	}

	return pubKey
}
