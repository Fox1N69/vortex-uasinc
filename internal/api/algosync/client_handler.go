package algosync

import (
	"strconv"
	"test-task/internal/models"
	service "test-task/internal/services"
	"test-task/pkg/http/response"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type ClientHandler interface {
	AddClient(c *gin.Context)
	UpdateClient(c *gin.Context)
	DeleteClient(c *gin.Context)
	UpdateAlgorithmStatus(c *gin.Context)
}

type clientHandler struct {
	service service.ClientService
}

func NewClientHandler(clientService service.ClientService) ClientHandler {
	return &clientHandler{service: clientService}
}

// @Summary Add new client to the database
// @Description AddClient creates a new client with the provided data.
// @Accept json
// @Produce json
// @Param body body models.Client true "Client object that needs to be added"
// @Success 200 {object} models.Client "Successfully created client"
// @Failure 400 {object} models.Response "error"
// @Failure 501 {object} models.Response "error"
// @Router /api/client/add [post]
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

// @Summary UpdateClient an existing client
// @Description UpdateClient updates the specified client with new data.
// @Accept json
// @Produce json
// @Param id path int true "Client ID to update"
// @Param body body map[string]interface{} true "Updated client data"
// @Success 200 {object} models.Client
// @Failure 400 {object} models.Response "error"
// @Router /api/client/{id} [patch]
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

	var client models.Client
	if err := mapstructure.Decode(updateParams, &client); err != nil {
		response.Error(400, err)
		return
	}

	c.JSON(200, gin.H{
		"id":      clientID,
		"message": "client update success",
	})
}

// @Summary Delete a client
// @Description DeleteClient deletes the client with the specified ID.
// @Accept json
// @Produce json
// @Param id path int true "Client ID to delete"
// @Success 200 {object} models.Client "Successfully deleted client"
// @Failure 400 {object} models.Response "error"
// @Failure 501 {object} models.Response "error"
// @Router /api/client/{id} [delete]
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

// @Summary Update algorithm status
// @Description UpdateAlgorithmStatus updates the algorithm status for the specified client.
// @Accept json
// @Produce json
// @Param id path int true "Algorithm ID to update"
// @Param body body map[string]interface{} true "Updated algorithm status data"
// @Success 200 {object} models.Client "Successfully updated algorithm status"
// @Failure 400 {object} models.Response "error"
// @Failure 501 {object} models.Response "error"
// @Router /api/client/algorithm/{id} [patch]
func (ch *clientHandler) UpdateAlgorithmStatus(c *gin.Context) {
	response := response.New(c)

	algorithmID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(400, err)
		return
	}

	var statusParams map[string]interface{}
	if err := c.ShouldBindJSON(&statusParams); err != nil {
		response.Error(400, err)
		return
	}

	if err := ch.service.UpdateAlgorithmStatus(algorithmID, statusParams); err != nil {
		response.Error(501, err)
		return
	}

	c.JSON(200, gin.H{
		"id":      algorithmID,
		"message": "algorithm updated success",
	})
}
