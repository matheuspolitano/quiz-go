package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/matheuspolitano/quiz-go/backend/internal/api"
	"github.com/matheuspolitano/quiz-go/backend/internal/config"
	"github.com/matheuspolitano/quiz-go/backend/internal/memdb"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	store, err := memdb.NewDBManager()
	if err != nil {
		log.Fatal(err)
	}
	svc, err := api.New(cfg, store)
	if err != nil {
		log.Fatal(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, interruptSignals...)

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Starting Quiz server on part %s", cfg.ApiPort)
		serverErrors <- svc.Start()
	}()

	select {
	case err := <-serverErrors:
		log.Printf("Could not start the server %v", err)
	case quitSignal := <-quit:
		log.Printf("Received signal %s. Initiating graceful shutdown...", quitSignal)
	}

	if err := svc.Shutdown(); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

}
