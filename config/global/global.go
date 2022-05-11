package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var Config *Configuration

var DB *gorm.DB

// RedisDB Redis客户端连接，一个redis客户端连接，默认底层维护10 * runtime.NumCPU()个TCP连接
//TCP连接，目前采用 redis.Client 默认连接参数，后续项目优化再添加其他参数
var RedisDB *redis.Client
