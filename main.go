package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var build = "develop"

func doing_stuff() {
	log.Println("go-service: doing some stuff")
}

func main() {
	log.Println("go-service: started", build)
	defer log.Println("my-service: ended")

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	doing_stuff()

	<-shutdown

	log.Println("go-service: stopped")
}
