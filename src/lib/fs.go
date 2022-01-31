package lib

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const pkgName = "lib"

// errInvalidPath indicates that the root path was invalid
var errInvalidPath = fmt.Errorf("(%s): invalid path", pkgName)

var filepathWalk = filepath.Walk

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

/*
WalkPath walks through the root directory, running the input `walkFunc` function on all
files. If the root path is invalid, a custom error is returned, use IsInvalidPathErr
to check for this.

Note: The parameter `walkFunc` will be selectively run on files
*/
func WalkPath(path string, walkFunc filepath.WalkFunc) error {
	if !PathExists(path) {
		return errInvalidPath
	}

	err := filepathWalk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		return walkFunc(path, info, err)
	})

	return errors.Wrapf(err, "(%s/WalkPath)", pkgName)
}
