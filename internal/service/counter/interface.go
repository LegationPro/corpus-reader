package counter

// Interface for the counter service
type ICounter interface {
	// Returns the current count
	GetCount() uint64
	// Increments the count by a given amount
	Increment(incrementBy int) error
	// Count the number of times a word is found in the file
	countWord(filePath string) error
	// Count the number of times a word is found starting from the root directory recursively

	processFile(path string, ch chan<- error)
	processDirectory(path string, ch chan<- error)

	Count() <-chan error
}
