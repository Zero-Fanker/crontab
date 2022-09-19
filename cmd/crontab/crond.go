package main

import (
	"crontab/internal/conf"
	"crontab/internal/job"
	"crontab/internal/logger"
	"crontab/internal/web"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	var port *string = flag.String("port", ":8080", "web port")
	var logPath *string = flag.String("logs", "logs/", "log path")
	var confs *string = flag.String("conf", "crontab.conf", "crontab config")

	logger.Configure(logPath)

	loaded, loadErr := conf.Load(confs)
	if !loaded {
		fmt.Printf("Err %s exit.\n", loadErr)
		os.Exit(1)
	}

	job.Run()

	http.HandleFunc("/set", web.Set)
	http.HandleFunc("/get", web.Get)
	http.HandleFunc("/del", web.Del)
	http.HandleFunc("/log", web.Loger)
	http.HandleFunc("/load", web.Load)
	http.HandleFunc("/stop", web.Stop)
	http.HandleFunc("/start", web.Start)
	http.HandleFunc("/status", web.Status)

	startErr := http.ListenAndServe(*port, nil)
	if startErr != nil {
		fmt.Println("Start server failed.", startErr)
		os.Exit(1)
	}
}
