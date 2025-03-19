package main

import (
	"log"
	"time"

	"github.com/LegationPro/corpus-reader/internal/cli"
	"github.com/LegationPro/corpus-reader/internal/service/counter"
)

func main() {
	parsedArgs, err := cli.ParseFlags()

	if err != nil {
		log.Fatal(err)
	}

	// Benchmark timer for the counterService
	startTime := time.Now()

	counterService := counter.New(parsedArgs.Word, parsedArgs.Dir)
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

	// Calculate elapsed time
	duration := time.Since(startTime)

	log.Printf("Word '%s' found %d times\n", parsedArgs.Word, counterService.GetCount())
	log.Printf("Duration: %s", duration)
}
