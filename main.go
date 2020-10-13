package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/ne3s-diag-handler/diag-snapshot/config"
	"github.com/ne3s-diag-handler/diag-snapshot/controller"
	"github.com/ne3s-diag-handler/diag-snapshot/utility"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
)

func init() {
	flag.String("minio-file-count", "", "Files count to remove from mminio")
	flag.String("days", "", "No of days")
}

func main() {

	utility.Logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
	utility.Config = config.LoadEnv()

	// PROFILING
	debug := http.Server{
		Addr:           ":6000",
		Handler:        http.DefaultServeMux,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("main: Debug listening %s", "6000")
		log.Printf("main: Debug listener closed: %v", debug.ListenAndServe())
	}()
	// PROFILING

	controller.RunController(utility.Config)
}
