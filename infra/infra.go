package infra

import (
	"context"
	"errors"
	"sync"
	"test-task/infra/k8s"
	"test-task/pkg/util/logger"
	"test-task/storage/postgres"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Infra interface {
	Config() *viper.Viper
	GetLogger() logger.Logger
	SetMode() string
	Port() string
	RedisClient() *redis.Client
	PSQLClient() *postgres.PSQLClient
	RunSQLMigrations()
	KubernetesDeployer() k8s.KubernetesDeployer
}

type infra struct {
	configFile string
}

func New(configFile string) Infra {
	return &infra{configFile: configFile}
}

var (
	vprOnce sync.Once
	vpr     *viper.Viper
)

// Config returns the Viper configuration instance.
// It reads and initializes configuration from the specified configFile path.
func (i *infra) Config() *viper.Viper {
	vprOnce.Do(func() {
		viper.SetConfigFile(i.configFile)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Fatalf("[infra][Config][viper.ReadInConfig] %v", err)
		}

		vpr = viper.GetViper()
	})

	return vpr
}

// GetLogger returns the application-wide logger instance.
func (i *infra) GetLogger() logger.Logger {
	log := logger.GetLogger()
	return log
}

var (
	modeOnce    sync.Once
	mode        string
	development = "dev"
	production  = "release"
)

// SetMode sets the application mode based on the environment configuration.
// It retrieves the mode from the environment settings and configures Gin framework accordingly.
func (i *infra) SetMode() string {
	modeOnce.Do(func() {
		env := i.Config().Sub("environment").GetString("mode")
		if env == development {
			mode = gin.DebugMode
		} else if env == production {
			mode = gin.ReleaseMode
		} else {
			logrus.Fatalf("[infa][SetMode] %v", errors.New("environment not setup"))
		}

		gin.SetMode(mode)
	})

	return mode
}

var (
	portOnce sync.Once
	port     string
)

// Port retrieves the server port from the configuration.
// It initializes the port once and returns it prefixed with ':'.
func (i *infra) Port() string {
	portOnce.Do(func() {
		port = i.Config().Sub("server").GetString("port")
	})

	return ":" + port
}

var (
	rdbOnce sync.Once
	rdb     *redis.Client
)

// RedisClient returns a Redis client instance configured based on the environment settings.
// It initializes the Redis client once using the configured address, password, and database.
func (i *infra) RedisClient() *redis.Client {
	rdbOnce.Do(func() {
		config := i.Config().Sub("redis")
		addr := config.GetString("addr")
		password := config.GetString("password")
		db := config.GetInt("db")

		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		if _, err := rdb.Ping(context.Background()).Result(); err != nil {
			logrus.Fatalf("[infra][RedisClient][rdb.Ping] %v", err)
		}

		logrus.Println("Connected to Redis")
	})

	return rdb
}

// PSQLClient returns a PostgreSQL client instance initialized with the configuration settings.
// It creates a new PostgreSQL client and establishes a connection using provided credentials.
func (i *infra) PSQLClient() *postgres.PSQLClient {
	config := i.Config().Sub("database")
	user := config.GetString("user")
	pass := config.GetString("pass")
	host := config.GetString("host")
	port := config.GetString("port")
	name := config.GetString("name")

	psqlClient := postgres.NewPSQLClient()
	psqlClient.Connect(user, pass, host, port, name)

	return psqlClient
}

// RunSQLMigrations runs SQL migrations on the configured PostgreSQL database.
// It triggers SQL migrations using the initialized PostgreSQL client.
func (i *infra) RunSQLMigrations() {
	i.PSQLClient().SqlMigrate()
}

// KubernetesDeployer returns a new instance of KubernetesDeployer.
// It initializes a Kubernetes deployer used for managing deployments.
func (i *infra) KubernetesDeployer() k8s.KubernetesDeployer {
	return k8s.NewKubernetesDeployer()
}
