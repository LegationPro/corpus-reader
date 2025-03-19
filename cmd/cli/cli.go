package main

import (
	"log"

	"github.com/LegationPro/corpus-reader/internal/cli"
	"github.com/LegationPro/corpus-reader/internal/service/counter"
)

func main() {
	// Parse the flags from the CLI
	parsedArgs, err := cli.ParseFlags()

	if err != nil {
		log.Fatal(err)
	}

	// Create a new counter service
	counterService := counter.New(parsedArgs.Word, parsedArgs.Dir, parsedArgs.MaxWorkers)
	errChan := counterService.Count()

	// Read errors from the channel
	for err := range errChan {
		if err != nil {
			log.Printf("Error occurred: %v", err)
		}
	}

	if err != nil {
		log.Fatalf("error when counting word: %v", err)
	}

	log.Printf("Word '%s' found %d times\n", parsedArgs.Word, counterService.GetCount())
}
