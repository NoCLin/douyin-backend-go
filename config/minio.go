package config

import (
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func initMinIO() *minio.Client {

	// Initialize minio client object.
	minioClient, err := minio.New(G.Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(G.Config.AccessKey, G.Config.SecretAccessKey, ""),
		Secure: G.Config.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}

func initTestMinio() *minio.Client {
	client, err := minio.New("a", &minio.Options{
		Creds: credentials.NewStaticV4("ss", "ss", ""),
	})
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
