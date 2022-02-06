/*
Package writer handles the part of writing output to a file

Use the function writer.Start to initialize the package
*/
package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/notsatan/crcgen/src/logger"
)

const pkgName = "writer"

var once sync.Once

// output is the central instance of `viper` being used by package `writer`
var output = viper.New()

var (
	readViperConfigs = (*viper.Viper).ReadInConfig // maps to viper.ReadInConfig
	fileIsDir        = os.FileInfo.IsDir           // maps to method IsDir in os.FileInfo
	closeFile        = (*os.File).Close            // maps to method Close in os.File

	pathStats  = os.Stat      // maps to os.Stat
	createFile = os.Create    // maps to os.Create
	absPath    = filepath.Abs // maps to filepath.Abs
)

var (
	errInvalidFile = fmt.Errorf("(%s): could not detect output file in path", pkgName)
	errInvalidExt  = fmt.Errorf("(%s): output file has invalid extension", pkgName)
	errAbsPath     = fmt.Errorf("(%s): couldn't convert path to absolute", pkgName)
	errNotWritable = fmt.Errorf("(%s): path is not writeable", pkgName)
	errPathIsDir   = fmt.Errorf("(%s): path points to an existing directory", pkgName)
)

/*
IsInvalidFileErr checks if an error returned by package writer was caused because the
output file could not be located from the path
*/
func IsInvalidFileErr(err error) bool {
	return errors.Is(err, errInvalidFile)
}

/*
IsInvalidExtErr checks if an error returned by package `writer` was caused by the file
having an invalid extension
*/
func IsInvalidExtErr(err error) bool {
	return errors.Is(err, errInvalidExt)
}

/*
IsAbsPathErr checks if an error occurred because a relative path could not be converted
to absolute path
*/
func IsAbsPathErr(err error) bool {
	return errors.Is(err, errAbsPath)
}

/*
IsPathNotWriteableErr checks if an error was caused because the path is not writeable
*/
func IsPathNotWriteableErr(err error) bool {
	return errors.Is(err, errNotWritable)
}

/*
IsPathDirErr checks if an error was caused because the path to output file points to an
existing directory
*/
func IsPathDirErr(err error) bool {
	return errors.Is(err, errPathIsDir)
}

/*
Start initializes package `writer` - should be executed before calls to any function
from this package

Returns error if the output file could not be parsed from the path, if the output file
contains an invalid extension, or if the path could not be converted to absolute path,
the path points to an existing directory, or if the path is not writeable. Use the
functions IsInvalidFileErr, IsInvalidExtErr, IsAbsPathErr, IsPathNotWriteableErr and
IsPathDirErr to explicitly check for these errors

Note: This function can be run once (at most)
*/
func Start(confPath string) error {
	var err error
	once.Do(func() {
		err = start(confPath)
	})

	return err
}

/*
start initializes package `writer`, setting up the configurations needed
*/
func start(confFile string) error {
	const logTag = "(" + pkgName + "/Start)"

	path, err := fixPath(confFile)
	if err != nil {
		return errors.Wrap(err, logTag)
	}

	output.SetConfigFile(path)

	err = createOutFile(path) // create output file if needed, ignored if file exists
	if err != nil {
		return errors.Wrap(err, logTag)
	}

	err = readViperConfigs(output)
	return errors.Wrap(err, logTag)
}

/*
fixPath validates, and fixes the path to the output file. This includes validating the
path, ensuring the file extension is valid, converting relative paths into absolute
paths, and more
*/
func fixPath(path string) (string, error) {
	const logTag = "(" + pkgName + "/fixPath)"

	if _, file := filepath.Split(path); file == "" {
		logger.Errorf(`%s: config path is empty`, logTag)
		return "", errors.Wrap(errInvalidFile, logTag)
	}

	// Extract and validate extension from path - remove `dot`, and convert to lowercase
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))
	switch ext {
	case "json", "yaml", "yml":
		// supported extensions, ignore

	default:
		logger.Errorf(`%s: config file invalid ext detected: "%s"`, logTag, path)
		return "", errors.Wrap(errInvalidExt, logTag)
	}

	p, err := absPath(path)
	if err != nil {
		logger.Errorf(`%s: could not resolve path to absolute: "%v"`, logTag, path)
		return "", errors.Wrap(errAbsPath, logTag)
	}

	return p, nil
}

/*
createOutFile checks if a file exists at the path for the output file, if not, it
creates an empty file

Errors returned can be checked using functions IsPathNotWriteableErr and IsPathDirErr
for granular error detection
*/
func createOutFile(path string) error {
	switch file, err := pathStats(path); {
	case os.IsNotExist(err):
		// do nothing, file needs to be created

	case err != nil:
		return errors.Wrapf(err, "(%s/createOutFile)", pkgName)

	case fileIsDir(file):
		return errors.Wrapf(errPathIsDir, "(%s/createOutFile)", pkgName)

	default: // path exists, and points to a file
		return nil
	}

	if file, err := createFile(path); err == nil {
		return errors.Wrapf(closeFile(file), "(%s/createOutFile)", pkgName)
	} else {
		// Assume failure is caused because the path is not writeable
		return errors.Wrapf(errNotWritable, "(%s/createOutFile)", pkgName)
	}
}
