package api

import (
	"test-task/infra"
	"test-task/internal/api/algosync"
	"test-task/internal/manager"
	"test-task/pkg/http/middleware"
	"test-task/pkg/http/request"
	"test-task/pkg/util/logger"

	"github.com/gin-gonic/gin"
	swaggerFile "github.com/swaggo/files"
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

// Run starts the server and initializes necessary middleware and handlers.
// It sets up rate limiting based on the configured RPS limit,
// enables CORS middleware, registers application handlers, and API routes.
// It also starts a background service to synchronize algorithm statuses.
// Finally, it logs the start of algorithm synchronization and listens on the configured port.
func (c *server) Run() {
	c.gin.Use(c.middleware.RPSLimit(c.infra.Config().GetInt("rps_limit")))

	c.gin.Use(c.middleware.CORS())
	c.handlers()
	c.v1()

	go c.startAlgorithmSync()

	log := logger.GetLogger()
	log.Info("Start algorithm sync")

	c.gin.Run(c.infra.Port())
}

func (c *server) startAlgorithmSync() {
	c.service.ClientService().StartAlgorithmSync()
}

// handlers sets up custom route handlers for specific routes on the server.
// It assigns default handlers for handling unknown routes and an index route.
func (c *server) handlers() {
	h := request.DefaultHandler()

	c.gin.NoRoute(h.NoRoute)
	c.gin.GET("/", h.Index)
}

// v1 configures versioned API endpoints (v1) for client operations.
// It sets up routes for client management operations such as adding, updating, deleting clients,
// and updating algorithm statuses associated with clients.
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
		}
	}

	c.gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
}
