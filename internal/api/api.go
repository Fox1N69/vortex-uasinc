package api

import (
	"test-task/infra"
	"test-task/internal/manager"
	"test-task/pkg/http/middleware"
	"test-task/pkg/http/request"

	"github.com/gin-gonic/gin"
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

func (c *server) Run() {
	c.gin.Use(c.middleware.CORS())
	c.handlers()
	c.v1()

	c.gin.Run(c.infra.Port())
}

func (c *server) handlers() {
	h := request.DefaultHandler()

	c.gin.NoRoute(h.NoRoute)
	c.gin.GET("/", h.Index)
}

func (c *server) v1() {

	c.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
