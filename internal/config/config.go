package config

import (
	"context"
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config interface {
	Config() *viper.Viper
	SetMode() string
	Port() string
	RedisClient() *redis.Client
}

type config struct {
	configFile string
}

func New(configFile string) Config {
	return &config{configFile: configFile}
}

var (
	vprOnce sync.Once
	vpr     *viper.Viper
)

func (cfg *config) Config() *viper.Viper {
	vprOnce.Do(func() {
		viper.SetConfigFile(cfg.configFile)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Fatalf("[config][Config][viper.ReadInConfig] %v", err)
		}

		vpr = viper.GetViper()
	})

	return vpr
}

var (
	modeOnce    sync.Once
	mode        string
	development = "dev"
	production  = "release"
)

func (cfg *config) SetMode() string {
	modeOnce.Do(func() {
		env := cfg.Config().Sub("environment").GetString("mode")
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

func (cfg *config) Port() string {
	portOnce.Do(func() {
		port = cfg.Config().Sub("server").GetString("port")
	})

	return ":" + port
}

var (
	rdbOnce sync.Once
	rdb     *redis.Client
)

func (cfg *config) RedisClient() *redis.Client {
	rdbOnce.Do(func() {
		config := cfg.Config().Sub("redis")
		addr := config.GetString("addr")
		password := config.GetString("password")
		db := config.GetInt("db")

		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		if _, err := rdb.Ping(context.Background()).Result(); err != nil {
			logrus.Fatalf("[config][RedisClient][rdb.Ping] %v", err)
		}

		logrus.Println("Connected to Redis")
	})

	return rdb
}
