package global

import (
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var Config *Configuration

var DB *gorm.DB

var MinioClient *minio.Client
