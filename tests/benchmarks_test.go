package tests

import (
	"testing"

	"github.com/LegationPro/corpus-reader/internal/service/counter"
)

func BenchmarkCounterService(b *testing.B) {
	// CLI flag arguments
	rootDir := "../corpus"
	word := "john"
	maxWorkers := 10

	counterService := counter.New(word, rootDir, maxWorkers)

	// Track memory usage
	b.ReportAllocs()

	for b.Loop() {
		counterService.Reset()
		errChan := counterService.Count()

		// Read errors from the channel, if any
		for err := range errChan {
			if err != nil {
				b.Fatalf("error when counting word: %v", err)
			}
		}

		// Get the count
		_ = counterService.GetCount()
	}
}
