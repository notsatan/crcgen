/*
Package writer handles the part of writing output to a file

Use the function writer.Start to initialize the package
*/
package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"

	"github.com/notsatan/crcgen/src/logger"
)

const pkgName = "writer"

/*
RootDir contains contents from the existing output file. Populated within the Start
function. Modifications can be made as needed, use the Write function to write the
updated contents back to the output file
*/
var RootDir DirInfo

var (
	fileIsDir = os.FileInfo.IsDir // maps to method IsDir in os.FileInfo
	closeFile = (*os.File).Close  // maps to method Close in os.File

	osReadFile  = os.ReadFile  // maps to os.ReadFile
	osWriteFile = os.WriteFile // maps to os.WriteFile
	pathStats   = os.Stat      // maps to os.Stat
	createFile  = os.Create    // maps to os.Create
	absPath     = filepath.Abs // maps to filepath.Abs
)

/*
Custom error
*/
var (
	errPathIsDir   = fmt.Errorf("(%s): path points to an existing directory", pkgName)
	errInvalidFile = fmt.Errorf("(%s): could not detect output file in path", pkgName)
	errInvalidExt  = fmt.Errorf("(%s): output file has invalid extension", pkgName)
	errAbsPath     = fmt.Errorf("(%s): couldn't convert path to absolute", pkgName)
	errNotWritable = fmt.Errorf("(%s): path is not writeable", pkgName)
	errReadFile    = fmt.Errorf("(%s): output file cannot be read", pkgName)
)

// once ensures Start can call inner start function exactly one time
var once sync.Once

// filePath contains full path to the output file, resolved and initialized in start
var filePath string

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
IsReadFileErr checks if an error was being caused because the output cannot be read
*/
func IsReadFileErr(err error) bool {
	return errors.Is(err, errReadFile)
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

	if path, err := fixPath(confFile); err != nil {
		return errors.Wrap(err, logTag)
	} else {
		filePath = path
	}

	// Create output file if needed, ignored if file exists
	if err := createOutFile(filePath); err != nil {
		logger.Errorf("%s: failed to create output file: %v", logTag, err)
		return errors.Wrap(err, logTag)
	}

	// Read content from output file, and unmarshal into `RootDir`
	if err := readFile(&RootDir); err != nil {
		logger.Errorf("%s: failed to read output file: %v", logTag, err)
		return errors.Wrap(errReadFile, logTag) // return custom error
	}

	return nil
}

/*
fixPath validates, and fixes the path to the output file. This includes validating the
path, ensuring the file extension is valid, converting relative paths into absolute
paths, and more
*/
func fixPath(path string) (string, error) {
	const logTag = "(" + pkgName + "/fixPath)"

	if _, file := filepath.Split(path); file == "" {
		logger.Errorf(`%s: output path is empty`, logTag)
		return "", errors.Wrap(errInvalidFile, logTag)
	}

	// Extract file extension, and validate the same
	if ext := filepath.Ext(path); !validateExt(ext) {
		logger.Errorf(`%s: output file invalid ext detected: "%s"`, logTag, path)
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

/*
readFile reads contents of the output file, using these to unmarshal all the data into
a DirInfo object. When the function completes its execution, the DirInfo object will
contain the contents of the output file
*/
func readFile(info *DirInfo) error {
	data, err := osReadFile(filePath)
	if len(data) == 0 || err != nil {
		// Direct return if output file has nothing to read, or an error occurred
		return errors.Wrapf(err, "(%s/readFile)", pkgName)
	}

	handler := getHandler(filepath.Ext(filePath)) // fetch handler based on file name
	err = handler.Unmarshal(data, info)
	return errors.Wrapf(err, "(%s/readFile)", pkgName)
}

/*
Write writes a DirInfo object to the output file while replacing existing contents in
the file
*/
func Write(info *DirInfo) error {
	const writePerm = 600 // assigns read, write

	handler := getHandler(filepath.Ext(filePath)) // fetch handler based on file name

	data, err := handler.Marshal(info, true)
	if err != nil {
		return errors.Wrapf(err, "(%s/Write)", pkgName)
	}

	err = osWriteFile(filePath, data, writePerm)
	return errors.Wrapf(err, "(%s/Write)", pkgName)
}
