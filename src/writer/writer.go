/*
Package writer handles the part of writing output to a file

Use the function Start to initialize the package
*/
package writer

import (
	"fmt"
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
	absPath          = filepath.Abs                // maps to filepath.Abs
)

var (
	errInvalidFile = fmt.Errorf("(%s): could not detect output file in path", pkgName)
	errInvalidExt  = fmt.Errorf("(%s): output file has invalid extension", pkgName)
	errAbsPath     = fmt.Errorf("(%s): couldn't convert path to absolute", pkgName)
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

func IsAbsPathErr(err error) bool {
	return errors.Is(err, errAbsPath)
}

/*
Start initializes package `writer` - should be executed before calls to any function
from this package

Returns error if the output file could not be parsed from the path, if the output file
contains an invalid extension, or if the path could not be converted to absolute path.
Use the functions IsInvalidFileErr, IsInvalidExtErr and IsAbsPathErr to explicitly
check for these errors

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
