package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LegationPro/corpus-reader/internal/server"
)

/*
The entry point for the server.
This function sets up the server configuration and starts the server.
It handles shutdown process by listening for shutdown signals and shutting down the server gracefully.
*/
func main() {
	parsedArgs, err := server.ParseFlags()

	if err != nil {
		log.Fatal(err)
	}

	// Create a context that will be used for the shutdown process
	shutdownContext, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Configure the server
	s := server.New(server.Config{
		Addr:         parsedArgs.Addr,
		ReadTimeout:  parsedArgs.ReadTimeout,
		WriteTimeout: parsedArgs.WriteTimeout,
		IdleTimeout:  parsedArgs.IdleTimeout,
	})

	// Start the server in a seperate goroutine
	go s.Start()

	// Wait for shutdown signal
	<-shutdownContext.Done()

	// Stop the server gracefully
	s.Stop()
}
