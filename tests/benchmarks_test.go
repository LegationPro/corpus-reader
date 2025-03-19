package tests

import (
	"testing"

	"github.com/LegationPro/corpus-reader/internal/service/counter"
)

func BenchmarkCounterService(b *testing.B) {
	// CLI flag arguments
	rootDir := "../corpus"
	word := "john"

	counterService := counter.New(word, rootDir)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Track memory usage
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
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
