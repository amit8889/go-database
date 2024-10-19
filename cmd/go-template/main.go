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

	"github.com/amit8889/go-project-template/internal/config"
)

func main() {
	fmt.Println("main function is called")
	cfg := config.MustLoad()
	fmt.Println("config loadded succesfully")
	//db setup

	//server setup

	server := http.Server{
		Addr: cfg.Server.Host + ":" + cfg.Server.PORT,
	}
	fmt.Println("server is running on port", cfg.Server.PORT)
	// make a go channel
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGALRM)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGALRM)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("==>Failed to start server", err)
		}
	}()
	<-done
	log.Println("server is shutting down")
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Error shutting down server:", err)
	}
	log.Println("server is shut down")

}
