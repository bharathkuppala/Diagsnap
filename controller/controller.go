package controller

import (
	"net/http"

	miniov1 "github.com/ne3s-diag-handler/diag-snapshot/minio"

	"github.com/ne3s-diag-handler/diag-snapshot/config"

	"github.com/ne3s-diag-handler/diag-snapshot/utility"
)

// RunController ...
func RunController(appConfig *config.AppConfig) {
	utility.MinioConnect = miniov1.SetupMinio(appConfig)
	if utility.MinioConnect == nil {
		utility.Logger.Println("problem in setting up minio client")
		return
	}
	router := http.NewServeMux()

	router.HandleFunc("/api/v1/view/files", FetchObjectKeys) // serving zip files
	// router.HandleFunc("/api/v1/download/files", nil)   // pull files from minio and zip
	router.HandleFunc("/api/v1/delete/files", deleteFiles) // delete from served zip files

	server := &http.Server{
		Addr:         ":" + appConfig.Server.Port,
		Handler:      router,
		WriteTimeout: appConfig.Server.TimeoutWrite,
		ReadTimeout:  appConfig.Server.TimeoutWrite,
	}

	if err := server.ListenAndServe(); err != nil {
		utility.Logger.Println("error in spinning up the server: " + err.Error())
		return
	}
}
