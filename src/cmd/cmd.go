/*
Package cmd defines the root command
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/notsatan/crcgen/src/cmd/version"
)

// Root is the central command that nests all subcommands under it
var Root = &cobra.Command{
	Use:   "crcgen",
	Short: "Show help for crcgen commands and flags",
	Long: `
crcgen batch generates file checksums for files in a directory

`,
	Version: version.Get(),
}

func init() {
	// Prints out the version as `crcgen v1.0.1`
	Root.SetVersionTemplate(
		"{{with .Name}}{{printf \"%s \" .}}{{end}}{{printf \"%s\" .Version}}\n\n",
	)
}
