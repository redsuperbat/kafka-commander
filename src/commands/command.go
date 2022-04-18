package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redsuperbat/kafka-commander/src/server"
	"github.com/rs/zerolog/log"
)

var allowedEvents = map[string]bool{
	"ChatMessageSentEvent": true,
}

type Command map[string]interface{}

func (c Command) Valid() *server.ResponseError {
	value, ok := c["eventType"].(string)
	if !ok {
		return server.NewRespErr(http.StatusBadRequest, "Invalid command, field eventType must be a string")
	}
	if value == "" {
		return server.NewRespErr(http.StatusBadRequest, "Invalid command, field eventType is required")
	}
	if !strings.HasSuffix(value, "Event") {
		return server.NewRespErr(http.StatusBadRequest, "Invalid command, field eventType requires suffix 'Event'")
	}

	if !allowedEvents[value] {
		errMsg := fmt.Sprintf("Invalid Command, eventType [%s] is not allowed.", value)
		return server.NewRespErr(http.StatusBadRequest, errMsg)
	}

	return nil
}

func (c Command) EventType() string {
	return c["eventType"].(string)
}

func HandleCommand(writeMessageFunc func([]byte) error, ctx *gin.Context) {
	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	var body Command
	json.Unmarshal(jsonData, &body)

	if validityErr := body.Valid(); validityErr != nil {
		ctx.JSON(validityErr.Code, validityErr)
		return
	}

	err := writeMessageFunc(jsonData)

	if err != nil {
		log.Info().Msgf("Error occurred when writing to kafka %s", err.Error())
		server.SendDefaultErr(ctx, http.StatusInternalServerError)
		return
	}

	log.Info().Msgf("Successfully commanded %s", body.EventType())
}
