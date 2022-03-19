package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return e.Message
}

func NewRespErr(code int, messages ...string) *ResponseError {
	var message string
	if len(messages) == 0 {
		message = http.StatusText(code)
	} else {
		message = strings.Join(messages, ", ")
	}
	return &ResponseError{
		Message: message,
		Code:    code,
	}
}

func SendDefaultErr(c *gin.Context, code int) {
	c.JSON(code, NewRespErr(code))
}
