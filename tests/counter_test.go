package tests

import (
	"testing"

	"github.com/LegationPro/corpus-reader/internal/service/counter"
)

func TestCounterService(t *testing.T) {
	// CLI flag arguments
	rootDir := "../corpus"
	word := "john"
	maxWorkers := 10

	counterService := counter.New(word, rootDir, maxWorkers)
	errChan := counterService.Count()

	// Read errors from the channel, if any
	for err := range errChan {
		if err != nil {
			t.Fatalf("error when counting word: %v", err)
		}
	}

	expectedCount := 13
	actualCount := counterService.GetCount()

	if actualCount != uint64(expectedCount) {
		t.Errorf("expected count %d, got %d", expectedCount, actualCount)
	}
}
