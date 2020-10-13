package controller

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/klauspost/pgzip"
	"github.com/mholt/archiver"
	"github.com/minio/minio-go/v7"
	"github.com/ne3s-diag-handler/diag-snapshot/utility"
	"github.com/spf13/afero"
	"golang.org/x/net/context"
)

var (
	currentTime = time.Now()
	objectsChan = make(chan minio.ObjectInfo)
)

const (
	defaultBlockSize = 1 << 20 * 2
)

type archive interface {
	Create(out io.Writer) error
	Write(f archiver.File) error
	Close() error
}

// FetchObjectKeys ...
func FetchObjectKeys(w http.ResponseWriter, r *http.Request) {
	var objectKeys []string
	ctx := context.Background()
	objectCh := utility.MinioConnect.ListObjects(ctx, utility.Config.Minio.MinioBucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	if objectCh == nil {
		log.Println("No object keys found")
		return
	}

	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err.Error())
			return
		}
		objectKeys = append(objectKeys, object.Key)
	}

	fileName, err := getFileName(ctx, objectKeys[0], utility.Config.Minio.MinioBucket)
	if err != nil {
		log.Println(err.Error())
		return
	}

	FetchObjectContent(ctx, objectKeys, utility.Config.Minio.MinioBucket, fileName)
}

// FetchObjectContent ...
func FetchObjectContent(ctx context.Context, filesToZip []string, bucketName string, fileName string) {
	fmt.Println("FetchObjectContent...", len(filesToZip))
	var objectsData []*minio.Object
	for _, v := range filesToZip {
		r, err := utility.MinioConnect.GetObject(ctx, bucketName, v, minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		objectsData = append(objectsData, r)
	}
	CompressedArchiever(ctx, objectsData, bucketName, fileName)
}

// CompressedArchiever ...
func CompressedArchiever(ctx context.Context, minioObjectsData []*minio.Object, bucketName string, fileName string) {
	fs := afero.NewMemMapFs()
	writer, err := fs.Create(fileName + ".tar.gz")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	defer putMinioObject(ctx, bucketName, fileName, fs, writer)
	defer writer.Close()

	pgzipWriter, err := pgzip.NewWriterLevel(writer, pgzip.BestCompression)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = pgzipWriter.SetConcurrency(defaultBlockSize, runtime.GOMAXPROCS(2))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer pgzipWriter.Close()

	tarWriter := tar.NewWriter(pgzipWriter)
	defer tarWriter.Close()

	for _, mo := range minioObjectsData {
		info, err := mo.Stat()
		if err != nil {
			log.Println(err.Error())
			return
		}

		// ****************************************
		
		//******************************************

		if info.Key == writer.Name() {
			continue
		}

		fmt.Println(info.Key)

		header := &tar.Header{
			Name: info.Key,
			Mode: 0600,
			Size: info.Size,
		}

		err = tarWriter.WriteHeader(header)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if _, err = io.Copy(tarWriter, mo); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func putMinioObject(ctx context.Context, buck string, fileName string, fs afero.Fs, writer afero.File) {
	fopen, err := fs.Open(writer.Name())
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer fopen.Close()

	info, err := fs.Stat(writer.Name())
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(info.Size(), info.Name())

	upInfo, err := utility.MinioConnect.PutObject(ctx, buck, writer.Name(), fopen, info.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully uploaded bytes: ", upInfo)

}

// n days uploaded files should be deleted
func deleteFiles(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	days, err := utility.ConverAtoI(flag.Lookup("days").Value.String())
	if err != nil {
		checkError(err)
	}

	go func() {
		defer close(objectsChan)
		for objectInfo := range utility.MinioConnect.ListObjects(ctx, utility.Config.Minio.MinioBucket, minio.ListObjectsOptions{
			Recursive: true,
		}) {
			fileDays := objectInfo.LastModified.Day() - currentTime.Day()
			fileTime := utility.Abs(fileDays)
			if err != nil {
				checkError(err)
			}

			if fileTime >= days {
				utility.Logger.Printf("Objects whih will be deleted are: %s", objectInfo.Key)
				objectsChan <- objectInfo
			}
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for errChan := range utility.MinioConnect.RemoveObjects(ctx, utility.Config.Minio.MinioBucket, objectsChan, opts) {
		utility.Logger.Printf("ObjectName: %s Err: %s", errChan.ObjectName, errChan.Err)
	}
}

func checkError(e error) {
	utility.Logger.Println(e.Error())
	return
}

func downloadFiles(w http.ResponseWriter, r *http.Request) {
	// Under implementation
}
