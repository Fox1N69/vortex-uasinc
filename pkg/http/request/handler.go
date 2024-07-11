package request

import (
	"errors"
	"net/http"

	"test-task/pkg/http/response"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	NoRoute(c *gin.Context)
	Index(c *gin.Context)
}

type handler struct {
	// Stuff maybe needed for handler
}

func DefaultHandler() Handler {
	return &handler{}
}

// NoRoute handles requests for routes that are not found.
//
// It returns an HTTP 404 error with a message indicating that the route was not found.
func (h *handler) NoRoute(c *gin.Context) {
	response.New(c).Error(http.StatusNotFound, errors.New("route not found"))
}

// Index handles requests to the root endpoint ("/").
//
// It returns an HTTP 200 OK response with a message indicating that the application is running.
func (h *handler) Index(c *gin.Context) {
	response.New(c).Write(http.StatusOK, "application running")
}
