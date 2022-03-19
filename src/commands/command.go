package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redsuperbat/kafka-commander/src/server"
	"github.com/segmentio/kafka-go"
)

type Command map[string]interface{}

func (c Command) CheckValidity() (*string, *server.ResponseError) {
	value, ok := c["eventType"].(string)
	if ok {
		return &value, nil
	}
	return nil, server.NewRespErr(http.StatusBadRequest, "Invalid command, field eventType is required.")
}

func HandleCommand(producer *kafka.Writer, ctx *gin.Context) {
	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	var body Command
	json.Unmarshal(jsonData, &body)
	eventType, validityErr := body.CheckValidity()

	if validityErr != nil {
		ctx.JSON(validityErr.Code, validityErr.Message)
		return
	}

	msg := kafka.Message{
		Value: jsonData,
	}

	err := producer.WriteMessages(ctx.Request.Context(), msg)
	if err != nil {
		fmt.Println(err.Error())
		server.SendDefaultErr(ctx, http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully sent event", eventType)
}
