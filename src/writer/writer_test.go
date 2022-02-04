package writer

import (
	"fmt"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

//nolint:deadcode,unused // keeping this for future use
func reset() {
	readViperConfigs = (*viper.Viper).ReadInConfig
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

func TestStart(t *testing.T) {
	// Checks to ensure `Start` can be run exactly once
	readViperConfigs = func(*viper.Viper) error { return nil }

	once = sync.Once{} // reset to isolate this test

	invalidInput := "/root/" // expect an error since no file is specified
	assert.Errorf(
		t, Start(invalidInput),
		`expected failure for invalid input: "%s"`, invalidInput,
	)

	// On the second run, the function `Start` should return directly, and no error
	// should be possible even for invalid input
	assert.NoErrorf(t, Start(invalidInput), `function "Start" running multiple times`)
}

func TestInternalStart(t *testing.T) {
	readViperConfigs = func(*viper.Viper) error { return nil }

	for input, err := range map[string]error{
		"":             errInvalidFile, // no file specified
		"file.txt":     errInvalidExt,  // invalid extension
		"/tmp/":        errInvalidFile, // no file specified
		"/dest/file":   errInvalidExt,  // no extension specified
		"config.YAML":  nil,            // default to working directory
		"/file.yAML":   nil,            // case-insensitivity ensured
		"/config.YmL":  nil,
		"/config.JSon": nil,
	} {
		result := start(input)

		if err == nil {
			assert.NoErrorf(t, result, `unexpected error for input: "%s"`, input)
		} else {
			assert.Errorf(t, result, `no error for invalid input: "%s"`, input)
			assert.Truef(
				t, errors.Is(result, err), `failed at: ("%s", %v)`, input, result,
			)
		}
	}
}

func TestValidateExt(t *testing.T) {
	for input, expected := range map[string]bool{
		"jSon":             true, // case-insensitive expected
		"YaMl":             true,
		"yML":              true,
		"invalid input":    false,
		"unexpected-input": false,
	} {
		result := validateExt(input)
		assert.Equalf(t, expected, result, `check passed for invalid input "%s"`, input)
	}
}
