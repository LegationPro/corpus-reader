package server

import (
	"flag"
	"time"
)

/*
Struct that contains the arguments parsed from the command line (CLI)
Addr: Address to listen on
ReadTimeout: Read timeout
WriteTimeout: Write timeout
IdleTimeout: Idle timeout
*/
type Arguments struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

/*
Parse the flags from the CLI.
If there are no flags given for the CLI, return an error otherwise return the Arguments struct.
*/
func ParseFlags() (*Arguments, error) {
	addr := flag.String("addr", ":8080", "Address to listen on")
	readTimeout := flag.Duration("read-timeout", 5*time.Second, "Read timeout")
	writeTimeout := flag.Duration("write-timeout", 5*time.Second, "Write timeout")
	idleTimeout := flag.Duration("idle-timeout", 120*time.Second, "Idle timeout")

	flag.Parse()

	return &Arguments{
		Addr:         *addr,
		ReadTimeout:  *readTimeout,
		WriteTimeout: *writeTimeout,
		IdleTimeout:  *idleTimeout,
	}, nil
}
