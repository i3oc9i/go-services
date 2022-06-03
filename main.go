package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

func doingStuff() {
	log.Println("go-service: doing some stuff")
}

func main() {

	// Set the correct number of threads for the service
	// based on what is available by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		log.Fatalln("go-service: %w.\n", err)
	}

	nOfCPU := runtime.GOMAXPROCS(0)

	log.Printf("go-service: started - build[%s] CPU[%d].\n", build, nOfCPU)
	defer log.Println("go-service: terminated.")

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	doingStuff()

	<-shutdown

	log.Println("go-service: stopped.")
}
