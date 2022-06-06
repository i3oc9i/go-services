package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "go.uber.org/automaxprocs" // set GOMAXPROCS based on resource CPU limit
)

var build = "develop"

func doingStuff() {
	log.Println("sales: doing some stuff")
}

func main() {
	nOfCPU := runtime.GOMAXPROCS(0)

	log.Printf("sales: started - build[%s] CPU[%d].\n", build, nOfCPU)
	defer log.Println("sales: terminated.")

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	doingStuff()

	<-shutdown

	log.Println("sales: stopped.")
}
