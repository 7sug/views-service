package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"views-servive/config"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	settings := config.Read()
	mainCtx, cancelMainCtx := context.WithCancel(context.Background())

	defer cancelMainCtx()

	app := NewApp(mainCtx, settings)

	if err := app.InitServices(); err != nil {
		log.Println("Failed to initialized app", err)
		return
	}

	app.Start()
	log.Println("service up")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	if err := app.Stop(context.WithTimeout); err != nil {
		log.Println("service down", err)
		return
	}
}
