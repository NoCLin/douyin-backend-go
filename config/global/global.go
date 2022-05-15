package global

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var Config *Configuration

var DB *gorm.DB
var DBMock sqlmock.Sqlmock

var RedisDB *redis.Client
var RedisMock redismock.ClientMock

var MinioClient *minio.Client

var TokenSecret = []byte("tokenSecret")
