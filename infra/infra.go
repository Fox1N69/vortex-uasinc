package infra

import (
	"context"
	"errors"
	"sync"
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

var (
	logOnce sync.Once
	log     logger.Logger
)

func (i *infra) GetLogger() logger.Logger {
	logOnce.Do(func() {
		logger.Init(i.SetMode())
		log = logger.GetLogger()
	})
	return log
}

var (
	modeOnce    sync.Once
	mode        string
	development = "dev"
	production  = "release"
)

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

func (i *infra) RunSQLMigrations() {
	i.PSQLClient().SqlMigrate()
}
