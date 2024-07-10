package api

import (
	"test-task/infra"
	"test-task/internal/api/algosync"
	"test-task/internal/manager"
	"test-task/pkg/http/middleware"
	"test-task/pkg/http/request"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server interface {
	Run()
}

type server struct {
	infra      infra.Infra
	gin        *gin.Engine
	service    manager.ServiceManager
	middleware middleware.Middleware
}

func NewServer(infra infra.Infra) Server {
	return &server{
		infra:      infra,
		gin:        gin.Default(),
		service:    manager.NewServiceManager(infra),
		middleware: middleware.NewMiddleware(infra.Config().GetString("secret.key")),
	}
}

// Run api server
func (c *server) Run() {
	c.gin.Use(c.middleware.CORS())
	c.handlers()
	c.v1()
	c.service.ClientService().StartAlgorithmSync()
	logrus.Info("Start algorithm sync")

	c.gin.Run(c.infra.Port())
}

func (c *server) handlers() {
	h := request.DefaultHandler()

	c.gin.NoRoute(h.NoRoute)
	c.gin.GET("/", h.Index)
}

func (c *server) v1() {
	clientHandler := algosync.NewClientHandler(c.service.ClientService())

	api := c.gin.Group("/api")
	{
		client := api.Group("/client")
		{
			client.POST("/add", clientHandler.AddClient)
			client.PATCH("/:id", clientHandler.UpdateClient)
			client.DELETE("/:id", clientHandler.DeleteClient)
			client.PATCH("/algorithm/:id", clientHandler.UpdateAlgorithmStatus)
			client.POST("/algorithm/create", clientHandler.CreateAlgorithm)
		}
	}

	c.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
