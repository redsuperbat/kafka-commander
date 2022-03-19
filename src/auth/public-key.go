package auth

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/redsuperbat/kafka-commander/src/options"
)

func GetPubKey() string {

	key := os.Getenv("JWT_PUB_KEY")
	args := options.GetArgs()

	if args.PubKeyUrl == "" && key == "" {
		log.Fatalln("Please specify a JWT pub key either by using the JWT_PUB_KEY env variable or a url with the --pub-key-url flag")
	}

	if key != "" {
		return key
	}

	resp, err := http.Get(args.PubKeyUrl)
	if err != nil {
		log.Fatalln("No public key was found at", args.PubKeyUrl, err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Unable to read response body", err.Error())
	}

	stringifiedBody := string(body)
	if !strings.HasPrefix(stringifiedBody, "-----BEGIN RSA PUBLIC KEY-----") {
		log.Fatalln("Invalid public key in response body", stringifiedBody)
	}

	return stringifiedBody
}
