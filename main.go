package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redsuperbat/kafka-commander/src/auth"
	"github.com/redsuperbat/kafka-commander/src/commands"
	"github.com/redsuperbat/kafka-commander/src/producing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LOG_LEVEL    = "LOG_LEVEL"
	COMMAND_PATH = "COMMAND_PATH"
	SERVER_PORT  = "SERVER_PORT"
)

func initializeLogger() {
	inputLogLevel := os.Getenv(LOG_LEVEL)
	if inputLogLevel == "" {
		configZerolog(zerolog.DebugLevel)
		return
	}

	logLevel, err := zerolog.ParseLevel(inputLogLevel)
	if err != nil {
		log.Error().Msg(err.Error())
		logLevel = zerolog.TraceLevel
	}
	configZerolog(logLevel)
}

func configZerolog(logLevel zerolog.Level) {
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	initializeLogger()
	log.Info().Msg("Getting pubkey...")
	pubKey := auth.GetPubKey()
	log.Info().Msg("Successfully got public key")
	jwtMiddleware := auth.NewJwtMiddleware(pubKey)

	writeMessage, close := producing.NewConn()
	defer close()

	router := gin.Default()
	router.Use(jwtMiddleware)

	commandPath := os.Getenv(COMMAND_PATH)
	if commandPath == "" {
		commandPath = "/"
	}
	log.Info().Msgf("Mapping %s to handle the commands", commandPath)
	router.POST(commandPath, func(ctx *gin.Context) {
		commands.HandleCommand(writeMessage, ctx)
	})

	serverPort, err := strconv.Atoi(os.Getenv(SERVER_PORT))
	if err != nil {
		serverPort = 8887
	}

	log.Info().Msgf("Server starting on port %v", serverPort)
	router.Run(":" + fmt.Sprint(serverPort))
}
