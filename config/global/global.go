package global

import (
	"github.com/go-redis/redis"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var Config *Configuration

var DB *gorm.DB
var RedisDB *redis.Client
var MinioClient *minio.Client

var TokenSecret = []byte("tokenSecret")
