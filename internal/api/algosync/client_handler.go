package algosync

import (
	"strconv"
	"test-task/internal/models"
	service "test-task/internal/services"
	"test-task/pkg/http/response"

	"github.com/gin-gonic/gin"
)

type ClientHandler interface {
	AddClient(c *gin.Context)
	UpdateClient(c *gin.Context)
	DeleteClient(c *gin.Context)
}

type clientHandler struct {
	service service.ClientService
}

func NewClientHandler(clientService service.ClientService) ClientHandler {
	return &clientHandler{service: clientService}
}

func (ch *clientHandler) AddClient(c *gin.Context) {
	response := response.New(c)
	var client models.Client

	if err := c.ShouldBindJSON(&client); err != nil {
		response.Error(400, err)
		return
	}

	id, err := ch.service.Create(&client)
	if err != nil {
		response.Error(501, err)
		return
	}

	c.JSON(200, gin.H{"message:": "create client success", "id": id})
}

func (ch *clientHandler) UpdateClient(c *gin.Context) {
	response := response.New(c)

	clientID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(400, err)
		return
	}

	var updateParams map[string]interface{}
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		response.Error(400, err)
		return
	}

	if err := ch.service.Update(clientID, updateParams); err != nil {
		response.Error(501, err)
		return
	}

	c.JSON(200, gin.H{
		"id":      clientID,
		"message": "client update success",
	})
}

func (ch *clientHandler) DeleteClient(c *gin.Context) {
	response := response.New(c)

	clientID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(400, err)
		return
	}

	if err := ch.service.Delete(clientID); err != nil {
		response.Error(501, err)
		return
	}

	c.JSON(200, gin.H{
		"id":      clientID,
		"message": "client deleted success",
	})
}