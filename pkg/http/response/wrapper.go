package response

import (
	"test-task/internal/models"

	"github.com/gin-gonic/gin"
)

type Wrapper interface {
	Write(code int, message string)
	Error(code int, err error)
}

type wrapper struct {
	c *gin.Context
}

func New(c *gin.Context) Wrapper {
	return &wrapper{c: c}
}

// Write writes a JSON response with the provided HTTP status code and message.
//
// It serializes the response into JSON format using the provided code and message.
func (w *wrapper) Write(code int, message string) {
	w.c.JSON(code, models.Response{Code: code, Message: message})
}

// Error writes a JSON response with the provided HTTP status code and error message.
//
// It serializes the response into JSON format using the provided code and the
// error message extracted from the error object.
func (w *wrapper) Error(code int, err error) {
	w.c.JSON(code, models.Response{Code: code, Message: err.Error()})
}
