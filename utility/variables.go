package utility

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/ne3s-diag-handler/diag-snapshot/config"
)

var (
	// Logger ...
	Logger *log.Logger

	// Config ...
	Config *config.AppConfig

	// MinioConnect ...
	MinioConnect *minio.Client
)
