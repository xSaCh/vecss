package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"vecss/internal/mq"
	storageaws "vecss/internal/storage/aws"

	"vecss/internal/vus"
)

func main() {
	emitter, err := mq.NewRabbitMqEmitter("guest", "guest", "localhost")

	if err != nil {
		panic(err)
	}
	defer emitter.Connection.Close()
	emitter.Setup()

	s := storageaws.NewS3Repository()
	s.HandleBucket()

	svr := http.Server{
		Addr:    ":8080",
		Handler: vus.InitRouters(s, emitter),
	}

	// Start the server in a separate goroutine
	go func() {
		log.Printf("Server is running on port %s", svr.Addr)
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
			return
		}
	}()

	// Wait for an interrupt signal to gracefully shutdown the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalf("Server failed to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
