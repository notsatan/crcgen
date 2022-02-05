package writer

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func reset() {
	readViperConfigs = (*viper.Viper).ReadInConfig
	absPath = filepath.Abs
}

func TestIsInvalidExtErr(t *testing.T) {
	for err, expected := range map[error]bool{
		nil:                      false,
		errInvalidExt:            true,
		errInvalidFile:           false,
		fmt.Errorf("test error"): false,
		fmt.Errorf(""):           false,
		fmt.Errorf("(%s): output file has invalid extension", pkgName): false,
	} {
		assert.Equalf(
			t, expected, IsInvalidExtErr(err),
			`failed to match "%v" -> "%v"`, expected, err,
		)
	}
}

func TestIsInvalidFileErr(t *testing.T) {
	for err, expected := range map[error]bool{
		nil:                      false,
		errInvalidExt:            false,
		errInvalidFile:           true,
		fmt.Errorf("test error"): false,
		fmt.Errorf(""):           false,
		fmt.Errorf("(%s): output file has invalid extension", pkgName): false,
	} {
		assert.Equalf(
			t, expected, IsInvalidFileErr(err),
			`failed to match "%v" -> "%v"`, expected, err,
		)
	}
}

func TestIsAbsPathErr(t *testing.T) {
	for err, expected := range map[error]bool{
		nil:                             false,
		errInvalidExt:                   false,
		errInvalidFile:                  false,
		errors.Wrap(errAbsPath, "test"): true,
		fmt.Errorf("test error"):        false,
		fmt.Errorf(""):                  false,
		fmt.Errorf("(%s): output file has invalid extension", pkgName): false,
	} {
		assert.Equalf(
			t, expected, IsAbsPathErr(err),
			`failed to match "%v" -> "%v"`, expected, err,
		)
	}
}

func TestStart(t *testing.T) {
	// Checks to ensure `Start` can be run exactly once

	once = sync.Once{}       // reset to isolate this test
	invalidInput := "/root/" // should fail - no file is specified

	err := Start(invalidInput) // should fail at `writer.fixPath` in `writer.start`
	assert.Errorf(t, err, `expected failure for invalid input: "%s"`, invalidInput)

	// On the second run, the function `Start` should return directly, and no error
	// should be possible even for invalid input
	assert.NoErrorf(t, Start(invalidInput), `function "Start" running multiple times`)
}

func TestInternalStart(t *testing.T) {
	// Dry run the internal start function

	reset()
	readViperConfigs = func(*viper.Viper) error { return nil }
	validInput := "output.json"

	err := start(validInput)
	assert.NoErrorf(t, err, "unexpected error: %v", err)
}

func TestFixPath(t *testing.T) {
	reset()
	for input, val := range map[string]struct {
		err  error
		path string // contains relative path, needs to be converted
	}{
		"":             {err: errInvalidFile},   // no file specified
		"file.txt":     {err: errInvalidExt},    // invalid extension
		"file.out":     {err: errInvalidExt},    // invalid extension
		"/tmp/":        {err: errInvalidFile},   // no file specified
		"/dest/file":   {err: errInvalidExt},    // no extension specified
		"config.YAML":  {path: "./config.YAML"}, // default to working directory
		"file.mp4":     {err: errInvalidExt},    // invalid extension
		"/file.yAML":   {path: "/file.yAML"},    // case-insensitivity ensured
		"/config.YmL":  {path: "/config.YmL"},
		"/config.JSon": {path: "/config.JSon"},
	} {
		result, err := fixPath(input)

		if val.err == nil {
			assert.NoErrorf(t, err, `unexpected error for input: "%s"`, input)

			path, e := filepath.Abs(val.path)
			assert.NoErrorf(t, e, `unexpected error for test input: "%s"`, input)
			assert.Equalf(t, path, result, `failed for input: "%v"`, input)
		} else {
			assert.Emptyf(t, result, `expected empty path for input: "%v"`, input)
			assert.Errorf(t, err, `no error for invalid input: "%s"`, input)
			assert.Truef(
				t, errors.Is(err, val.err), `(input, error): ("%s", "%v")`, input, err,
			)
		}
	}
}

func TestFixPath_AbsPathFail(t *testing.T) {
	reset()

	absPath = func(string) (string, error) { return "", errAbsPath }
	path, err := fixPath("/test/path.json")

	assert.Error(t, err, "expected an error when absolute path can't be formed")
	assert.Emptyf(t, path, `expected path to be empty for an error: "%v"`, path)
	assert.Truef(t, IsAbsPathErr(err), `expected custom error type: "%v"`, err)
}
