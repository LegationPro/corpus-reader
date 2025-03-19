package counter

// Interface for the counter service
type ICounter interface {
	// Returns the current count
	GetCount() uint64
	// Increments the count by a given amount
	Increment(incrementBy int) error
	// Count the number of times a word is found in the file
	countWord(filePath string) error

	// Reset the counter
	Reset()

	// Process a file and send error to channel if any
	processFile(path string, ch chan<- error)

	// Process a directory and send error to channel if any
	processDirectory(path string, ch chan<- error)

	// Count the number of times a word is found starting from the root directory recursively
	Count() <-chan error
}
