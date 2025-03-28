package counter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

/*
The counter service is a simple service that counts the number of times a word is found in a file.
The count is stored as an unsigned 64-bit integer since the count could be very large, and won't go below 0.
We use atomic operations to ensure thread safety.

We prioritize the use of atomic operations over mutexes for the count since
it performs increments and read operations without locking.

The worker pool is used to limit the number of concurrent workers.
It uses a buffered channel to limit concurrency.
*/
type counter struct {
	wg         sync.WaitGroup
	count      atomic.Uint64
	word       string
	root       string
	workerPool chan struct{}
}

// Create a new counter instance
func New(word string, rootDir string, maxWorkers int) ICounter {
	return &counter{
		count:      atomic.Uint64{},
		word:       word,
		root:       rootDir,
		workerPool: make(chan struct{}, maxWorkers),
	}
}

// Helper function for running tasks with worker pool
func (c *counter) runTask(task func()) {
	// Increment WaitGroup counter
	c.wg.Add(1)

	/*
		The defers are executed in a specific order (LIFO - Last In First Out)
		1. c.wg.Done() - Decrement WaitGroup counter
		2. <-c.workerPool - Release the worker slot from the pool
		to ensure that the worker slot is released before the WaitGroup counter is decremented.
	*/

	go func() {
		// Decrement WaitGroup counter
		defer c.wg.Done()

		// Acquire a worker slot from the pool
		c.workerPool <- struct{}{}

		// Release the worker slot from the pool
		defer func() { <-c.workerPool }()

		// Execute the task
		task()
	}()
}

// Reset the counter
func (c *counter) Reset() {
	c.count.Store(0)
}

// Return the current count
func (c *counter) GetCount() uint64 {
	return c.count.Load()
}

// Increment count by a given amount
func (c *counter) Increment(amount int) error {
	/*
		Validation check for ensuring that the amount is non-negative.
		In a scenario, where a number would be negative and converted to a uint64,
		it would result in a very large number, which is not the intended behavior.
	*/
	if amount < 0 {
		return errors.New("amount must be non-negative value")
	}

	c.count.Add(uint64(amount))

	return nil
}

// Update the root path
func (c *counter) UpdateRoot(path string) {
	c.root = path
}

/*
Look for a directory with the given name under the root directory.
The function returns the path to the directory if found, otherwise it throws an error.
*/
func (c *counter) LookForDirectory(directoryName string) (string, error) {
	var foundPath string

	// Walk the directory tree and look for the directory
	err := filepath.WalkDir(c.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check if the current entry is a directory and matches the name
		if d.IsDir() && d.Name() == directoryName {
			// Store the path to the directory with the joined path
			foundPath = filepath.Join(c.root, directoryName)
			return nil
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("directory '%s' not found under root '%s'", directoryName, c.root)
	}

	return foundPath, nil
}

// Count the number of times a word is found in a file
func (c *counter) countWord(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	word := strings.ToLower(c.word)
	wordBytes := []byte(word)

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// Break out of the loop if the error is EOF (end of file)
			if err == io.EOF {
				break
			}
			return err
		}

		// Convert line to bytes and count occurrences of the word
		// Ensure that the line is converted to lowercase to match case-insensitive words
		// Ensure we trim any spaces, in an example where the user may have added a space to a word by accident

		lineBytes := []byte(strings.TrimSpace(strings.ToLower(line)))
		lineCount := bytes.Count(lineBytes, wordBytes)

		// Increment the global counter atomically for each occurrence to ensure thread safety
		c.count.Add(uint64(lineCount))
	}

	return nil
}

func (c *counter) processFile(path string, ch chan<- error) {
	// Skip non-text files
	if !strings.HasSuffix(path, ".txt") {
		ch <- nil
		return
	}

	// Count the word in the file and send error to channel if any
	if err := c.countWord(path); err != nil {
		select {
		case ch <- err:
			log.Printf("Sent error successfully")
		default:
			log.Fatalf("Failed to send error: %v", err)
		}
	}
}

func (c *counter) processDirectory(dirPath string, errChan chan<- error) {
	err := filepath.WalkDir(dirPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip processing the directory itself to avoid infinite recursion
			if dirPath != path {
				c.runTask(func() {
					c.processDirectory(path, errChan)
				})
			}
			return nil
		}

		// Dispatch file to a worker for processing
		c.runTask(func() {
			c.processFile(path, errChan)
		})

		return nil
	})

	if err != nil {
		// Send error to channel
		errChan <- err
	}
}

/*
Count the number of times a word is found starting from the root directory recursively.
The use of a worker pool allows for concurrent processing of files, directories and
an effective way to manage synchronization.
*/

func (c *counter) Count() <-chan error {
	// Use a buffered channel to prevent blocking
	errChan := make(chan error, 100)

	// Start processing the root directory
	c.runTask(func() {
		c.processDirectory(c.root, errChan)
	})

	// Wait for all goroutines to complete and close the error channel
	go func() {
		c.wg.Wait()
		close(errChan)
		log.Printf("Finished counting: %d", c.count.Load())
	}()

	return errChan
}
