package main

import (
	"fmt"
	"go-markov/markov"
	"os"
)

const (
	// Exit codes
	ExitCodeFatal = iota + 1
)

func main() {
	run()
}

// run is the entry point into the markov-image application
func run() {
	if len(os.Args) < 3 {
		usage()
	}
	m := markov.New()
	fatal(m.ReadFile(os.Args[1]), "failed to read input file")
	fatal(m.WriteFile(os.Args[2]), "failed to write output file")
}

// usage prints the applications usage/help and exits
func usage() {
	fmt.Printf(
		"usage: %s <input>.png <output>.png"+
			" -- generate a new image using Markov chaining\n",
		os.Args[0],
	)
	os.Exit(0)
}

// fatal prints the given message and error, and exits with a fatal code, if the
// given error is not nil.
func fatal(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		os.Exit(ExitCodeFatal)
	}
}
