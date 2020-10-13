package controller

import (
	"archive/tar"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/ne3s-diag-handler/diag-snapshot/utility"
	"golang.org/x/net/context"
)

// Raml ...
type Raml struct {
	XMLName xml.Name `xml:"raml"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Version string   `xml:"version,attr"`
	CmData  struct {
		Text          string `xml:",chardata"`
		Type          string `xml:"type,attr"`
		ManagedObject struct {
			Text      string `xml:",chardata"`
			Class     string `xml:"class,attr"`
			DistName  string `xml:"distName,attr"`
			Operation string `xml:"operation,attr"`
			Version   string `xml:"version,attr"`
			P         []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"p"`
		} `xml:"managedObject"`
	} `xml:"cmData"`
}

// writer *os.File, err error)
func getFileName(ctx context.Context, obKey string, bucketName string) (fileName string, err error) {
	var ramlData Raml
	r, err := utility.MinioConnect.GetObject(ctx, bucketName, obKey, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tarReader := tar.NewReader(r)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", fmt.Errorf("next failed: %w", err)
		}

		switch header.Typeflag {
		case tar.TypeReg:
			if header.Name != "raml.xml" {
				continue
			}

			byteValue, _ := ioutil.ReadAll(tarReader)

			if err := xml.Unmarshal(byteValue, &ramlData); err != nil {
				fmt.Println(err.Error())
				return "", err
			}

			indexPos := strings.Index(ramlData.CmData.ManagedObject.P[1].Text, ".")
			if indexPos == -1 {
				log.Println("No matching string found")
				return "", fmt.Errorf("No matching string found: %s", ".")
			}
			fileName = ramlData.CmData.ManagedObject.P[1].Text[0:indexPos]

			fmt.Printf("File name: %s", fileName)

		default:
			return "", fmt.Errorf("extract tar: unknown type: %v in %s", header.Typeflag, header.Name)
		}

		// t := archiver.DefaultTarGz.Create()

	}
	return fileName, nil
}
