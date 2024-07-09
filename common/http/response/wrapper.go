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

func (w *wrapper) Write(code int, message string) {
	w.c.JSON(code, models.Response{Code: code, Message: message})
}

func (w *wrapper) Error(code int, err error) {
	w.c.JSON(code, models.Response{Code: code, Message: err.Error()})
}
