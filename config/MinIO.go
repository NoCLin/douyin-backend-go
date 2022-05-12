package config

import (
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func initMinIO() error {

	// Initialize minio client object.
	minioClient, err := minio.New(G.Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(G.Config.AccessKey, G.Config.SecretAccessKey, ""),
		Secure: G.Config.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf(G.Config.AccessKey, G.Config.SecretAccessKey)
	//log.Printf("%#v\n", minioClient) // minioClient is now set up

	//return minioClient
	G.MinioClient = minioClient
	return nil
}
