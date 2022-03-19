package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redsuperbat/kafka-commander/src/auth"
	"github.com/redsuperbat/kafka-commander/src/commands"
	"github.com/redsuperbat/kafka-commander/src/options"
	"github.com/redsuperbat/kafka-commander/src/producing"
)

func main() {
	log.Println("Starting kafka-commander...")
	args := options.GetArgs()
	log.Println("Parsed args.")

	log.Println("Getting pubkey...")
	pubKey := auth.GetPubKey()
	log.Println("Successfully got public key")
	validateJwt := auth.NewJwtValidator(pubKey)

	router := gin.Default()
	producer := producing.NewKafkaWriter(args.KafkaBroker, args.KafkaTopic)
	defer producer.Close()

	log.Println("Mapping", args.CommandPath, "to handle the commands")
	router.POST(args.CommandPath, func(ctx *gin.Context) {
		bearerToken := ctx.Request.Header.Get("Authorization")
		user, err := validateJwt(bearerToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err)
			return
		}
		log.Println("User", user.Username, "tried issuing a command")

		commands.HandleCommand(producer, ctx)
	})
	log.Println("Server starting on port", args.ServerPort)
	router.Run(":" + fmt.Sprint(args.ServerPort))
}
