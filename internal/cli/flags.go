package cli

import (
	"errors"
	"flag"
)

/*
Struct that contains the arguments parsed from the command line (CLI)
Dir: Directory to scan for files
Word: Word to search for in the files
*/
type Arguments struct {
	Dir        string
	Word       string
	MaxWorkers int
}

/*
Parse the flags from the CLI.
If there are no flags given for the CLI, return an error otherwise return the Arguments struct.
*/
func ParseFlags() (*Arguments, error) {
	dir := flag.String("dir", "", "Directory to scan for files")
	word := flag.String("word", "", "Word to search for")
	maxWorkers := flag.Int("max-workers", 10, "Maximum number of workers")

	flag.Parse()

	if *dir == "" || *word == "" {
		return nil, errors.New("missing required fields: --dir and --word must be provided")
	}

	return &Arguments{
		Dir:        *dir,
		Word:       *word,
		MaxWorkers: *maxWorkers,
	}, nil
}
