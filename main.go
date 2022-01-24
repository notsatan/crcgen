package main

import (
	"os"

	"github.com/notsatan/crcgen/src"
	"github.com/notsatan/crcgen/src/logger"
)

const pkgName = "main"

var (
	// run maps to src.Run
	run = src.Run

	// exit maps calls to os.Exit
	exit = os.Exit
)

func main() {
	if err := run(); err != nil {
		logger.Errorf("(%s/main): %s: %s", pkgName, err)
		exit(-10)
	}
}
