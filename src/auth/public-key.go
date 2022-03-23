package auth

import (
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

const (
	JWT_PUB_KEY     = "JWT_PUB_KEY"
	JWT_PUB_KEY_URL = "JWT_PUB_KEY_URL"
)

func parseKey[T string | []byte](key T) *rsa.PublicKey {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		log.Fatal().Msgf("Invalid RSA public key %s", err.Error())
	}
	return pubKey
}

func GetPubKey() *rsa.PublicKey {
	key := os.Getenv(JWT_PUB_KEY)
	if key != "" {
		return parseKey(key)
	}

	keyUrl := os.Getenv(JWT_PUB_KEY_URL)

	if keyUrl == "" {
		log.Fatal().Msg("Please specify a JWT pub key either by using the 'JWT_PUB_KEY' or 'JWT_PUB_KEY_URL' env variables")
	}

	url, err := url.Parse(keyUrl)

	if err != nil {
		log.Fatal().Msgf("Invalid url %s", url.String())
	}

	resp, err := http.Get(url.String())

	if err != nil || resp.StatusCode >= 400 {
		log.Fatal().Msgf("No public key was found at %s", keyUrl)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Msgf("Unable to read response body %s", err.Error())
	}
	var result map[string]any

	json.Unmarshal(body, &result)

	publicKeyString := result["publicKey"].(string)

	return parseKey(publicKeyString)
}
