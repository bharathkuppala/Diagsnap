package miniov1

import (
	"log"

	"github.com/ne3s-diag-handler/diag-snapshot/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// SetupMinio ...
func SetupMinio(appConfig *config.AppConfig) *minio.Client {
	minioClient, err := newMinioClient(appConfig.Minio.MinioEndpoint, appConfig.Minio.MinioAccessKey, appConfig.Minio.MinioSecretKey)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return minioClient
}

// newMinioClient ...
func newMinioClient(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	mc, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return mc, nil
}
