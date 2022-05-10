package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var Config *Configuration

var DB *gorm.DB
var RedisDB *redis.Client
