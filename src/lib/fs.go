package lib

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const pkgName = "lib"

// errInvalidPath indicates that the root path was invalid
var errInvalidPath = fmt.Errorf("(%s): invalid path", pkgName)

/*
IsInvalidPathErr indicates if an error was returned because of invalid input path
*/
func IsInvalidPathErr(err error) bool {
	return errors.Is(err, errInvalidPath)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
