package common

import (
	"github.com/creamdog/gonfig"
	"github.com/go-redis/redis"
)

// Config - Contains struct for reading global configuration values
var Config gonfig.Gonfig

// RedisClient - System client to access Reds
var RedisClient *redis.Client

// SystemConfiguration - ...
var GlobalConfig struct {
	ListenPort      int
	DevelopmentMode bool
}
