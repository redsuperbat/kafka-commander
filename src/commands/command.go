package commands

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redsuperbat/kafka-commander/src/server"
)

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
		log.Println("Error occurred when writing to kafka", err.Error())
		server.SendDefaultErr(ctx, http.StatusInternalServerError)
		return
	}

	log.Println("Successfully commanded", body.EventType())
}
