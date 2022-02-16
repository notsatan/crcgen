/*
Package cmd is the central handler that retains the flow-of-control throughout the
lifetime
*/
package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/notsatan/crcgen/src/cmd/version"
	"github.com/notsatan/crcgen/src/logger"
	_ "github.com/notsatan/crcgen/src/writer/json"
)

const (
	pkgName   = "cmd"
	debugMode = "debug_mode"
)

var (
	exit        = os.Exit
	initLogger  = logger.Log
	execCmd     = Root.Execute
	cmdUsage    = (*cobra.Command).Usage
	closeLogger = logger.Stop
)

// Root is the central command that nests all subcommands under it
var Root = &cobra.Command{
	Use:   "crcgen",
	Short: "Show help for crcgen commands and flags",
	Long: `
crcgen batch generates file checksums for files in a directory

`,
	Version: version.Get(),
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.SetOut(cmd.OutOrStdout())

		err := cmdUsage(cmd) // dump usage message to `stdout`
		if err != nil {
			logger.Warnf("(%s/Root.Run): %v", pkgName, err)
			closeRes()
			exit(-10)
		}
	},
}

func init() {
	// Prints out the version as `crcgen v1.0.1`
	Root.SetVersionTemplate(
		"{{with .Name}}{{printf \"%s \" .}}{{end}}{{printf \"%s\" .Version}}\n\n",
	)

	setupCmdTemplate(Root)
}

/*
Run sets up the resources needed and runs the Root command, passing the flow-of-control
to the main command-line interface

Once the root command completes, Run is also responsible to safely close the resources,
before returning the flow-of-control
*/
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

	err := execCmd() // run the root command
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
