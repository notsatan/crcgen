package main

import (
	"os"

	"github.com/notsatan/crcgen/src"
	"github.com/notsatan/crcgen/src/logger"
)

const pkgName = "main"

var (
	run  = src.Run
	exit = os.Exit
)

func main() {
	if err := run(); err != nil {
		logger.Errorf("(%s/main): %s: %s", pkgName, err)
		exit(-10)
	}
}
