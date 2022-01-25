/*
Package src acts as an interface between the main application and package `main`
*/
package src

import (
	"os"

	"github.com/pkg/errors"

	"github.com/notsatan/crcgen/src/cmd"
	"github.com/notsatan/crcgen/src/logger"
)

const (
	pkgName   = "src"
	debugMode = "debug_mode"
)

var (
	initLogger  = logger.Log
	execCmd     = cmd.Root.Execute
	closeLogger = logger.Stop
)

func Run() error {
	defer closeRes()

	if _, ok := os.LookupEnv(debugMode); ok {
		// In `debug` mode, switch logger to write all logs to `stderr`
		err := initLogger(false)
		if err != nil {
			return errors.Wrapf(err, "(%s/Run)", pkgName)
		}

		logger.Debugf("(%s/main): Detected `debug` mode", pkgName)
	}

	// Run the root command and return the result
	err := execCmd()
	return errors.Wrapf(err, "(%s/Run)", pkgName)
}

/*
closeRes attempts to close resources, designed to be run before the app closes
*/
func closeRes() {
	if err := closeLogger(); err != nil {
		logger.Errorf("(%s/main): Failed to close logger: %s", pkgName, err)
	}
}
