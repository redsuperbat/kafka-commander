package main

import (
	"fmt"
	"log"

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
	jwtMiddleware := auth.NewJwtMiddleware(pubKey)

	producer := producing.NewKafkaWriter(args.KafkaBroker, args.KafkaTopic)
	defer producer.Close()

	router := gin.Default()
	router.Use(jwtMiddleware)

	log.Println("Mapping", args.CommandPath, "to handle the commands")
	router.POST(args.CommandPath, func(ctx *gin.Context) {
		commands.HandleCommand(producer, ctx)
	})
	log.Println("Server starting on port", args.ServerPort)
	router.Run(":" + fmt.Sprint(args.ServerPort))
}
