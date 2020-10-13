package controller

import (
	"github.com/minio/minio-go/v7"
	"github.com/ne3s-diag-handler/diag-snapshot/utility"
	"golang.org/x/net/context"
)

func listObjects(ctx context.Context) (objectCh <-chan minio.ObjectInfo) {
	objectCh = utility.MinioConnect.ListObjects(ctx, utility.Config.Minio.MinioBucket, minio.ListObjectsOptions{
		Recursive: true,
	})
	return objectCh
}
