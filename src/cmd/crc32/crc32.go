/*
Package crc32 implements the CRC-32 check
*/
package crc32

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/notsatan/crcgen/src/cmd"
)

func init() {
	cmd.Root.AddCommand(command)
}

// command defines the CRC32 algorithm
var command = &cobra.Command{
	Use:   "crc32 path",
	Short: "Deals with generation of CRC-32 checksums",
	Long: `
crc32 generates the CRC-32 checksum for files using the IEEE polynomial
by default
`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PreRun: nil,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.SetOut(cmd.OutOrStdout())
			_ = cmd.Usage() // dump usage message to `stdout`

			return fmt.Errorf("could not detect root path")
		}

		return nil
	},
}

/*
Calculate generates a CRC-32 checksum for a file, and returns the same
*/
func Calculate(path string) (string, error) {
	return "", nil
}
