package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redsuperbat/kafka-commander/src/commands"
	"github.com/redsuperbat/kafka-commander/src/options"
	"github.com/redsuperbat/kafka-commander/src/producing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	args := options.GetArgs()
	zerolog.SetGlobalLevel(args.LogLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("Getting pubkey...")
	// pubKey := auth.GetPubKey()
	log.Info().Msg("Successfully got public key")
	// jwtMiddleware := auth.NewJwtMiddleware(pubKey)

	writeMessage, close := producing.NewConn()
	defer close()

	router := gin.Default()
	// router.Use(jwtMiddleware)

	log.Info().Msgf("Mapping %s to handle the commands", args.CommandPath)
	router.POST(args.CommandPath, func(ctx *gin.Context) {
		commands.HandleCommand(writeMessage, ctx)
	})
	log.Info().Msgf("Server starting on port %v", args.ServerPort)
	router.Run(":" + fmt.Sprint(args.ServerPort))
}
