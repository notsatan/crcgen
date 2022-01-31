package main

import (
	"os"

	"github.com/notsatan/crcgen/src/cmd"
	"github.com/notsatan/crcgen/src/logger"
)

const pkgName = "main"

var (
	run  = cmd.Run
	exit = os.Exit
)

func main() {
	if err := run(); err != nil {
		logger.Errorf("(%s/main): %s", pkgName, err)
		exit(-10)
	}
}
