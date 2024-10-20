package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amit8889/go-database/internal/config"
)

func main() {
	fmt.Println("main function is called")
	cfg := config.MustLoad()
	fmt.Println("config loaded successfully")

	// server setup
	server := http.Server{
		Addr: cfg.Server.Host + ":" + cfg.Server.PORT,
	}
	fmt.Println("server is running on port", cfg.Server.PORT)

	// make a go channel for signal handling
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-done
	log.Println("server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Error shutting down server:", err)
	}
	log.Println("server is shut down")
}
