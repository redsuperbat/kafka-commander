package auth

import (
	"crypto/rsa"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/redsuperbat/kafka-commander/src/options"
	"github.com/rs/zerolog/log"
)

func GetPubKey() *rsa.PublicKey {

	key := os.Getenv("JWT_PUB_KEY")
	args := options.GetArgs()

	if args.PubKeyUrl == "" && key == "" {
		log.Fatal().Msg("Please specify a JWT pub key either by using the JWT_PUB_KEY env variable or a url with the --pub-key-url flag")
	}

	if key != "" {
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
		if err != nil {
			log.Fatal().Msgf("Invalid RSA public key %s", err.Error())
		}
		return pubKey
	}

	resp, err := http.Get(args.PubKeyUrl)
	if err != nil {
		log.Fatal().Msgf("No public key was found at %s %s", args.PubKeyUrl, err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Msgf("Unable to read response body %s", err.Error())
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(body)
	if err != nil {
		log.Fatal().Msgf("Invalid public key in response body %s", err.Error())
	}

	return pubKey
}
