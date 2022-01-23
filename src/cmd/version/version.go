package version

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
// version contains the raw version number in the form `1.2.3`
var version string

func init() {
	// Convert version to be of the form v1.2.3; and trim newlines and spaces
	version = "v" + strings.Trim(version, "\n ")
}

/*
Get returns the current version number as a string
*/
func Get() string {
	return version
}
