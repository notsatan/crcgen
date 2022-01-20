package main

import (
	"os"

	"github.com/notsatan/crcgen/src/logger"
)

const (
	pkgName   = "main"
	debugMode = "debug_mode"
)

var (
	// initLogger maps calls to logger.Log
	initLogger = logger.Log

	// exit maps calls to os.Exit
	exit = os.Exit

	// closeLogger maps calls to logger.Stop
	closeLogger = logger.Stop
)

func main() {
	defer closeRes()

	if _, ok := os.LookupEnv(debugMode); ok {
		// In `debug` mode, switch logger to write all logs to `stderr`
		err := initLogger(false)
		crashOnErr(err, "Failed to start logger")

		logger.Debugf("(%s/main): Detected `debug` mode", pkgName)
	}
}

/*
crashOnErr is a simple helper function to ensure a "graceful" crash if the main thread
encounters an error.

The function prints an error message to the user, closes resources and force closes the
application
*/
func crashOnErr(err error, cause string) {
	if err == nil {
		return
	}

	logger.Errorf("(%s/main): %s: %s", pkgName, cause, err)
	closeRes() // close resources

	exit(-10)
}

/*
closeRes attempts to close resources, designed to be run before the app closes
*/
func closeRes() {
	if err := closeLogger(); err != nil {
		logger.Errorf("(%s/main): Failed to close logger: %s", pkgName, err)
	}
}
