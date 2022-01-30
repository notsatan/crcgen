package lib

import (
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIsInvalidPathErr(t *testing.T) {
	for _, err := range []error{
		fmt.Errorf("(%s): invalid path", pkgName),
		errors.New("(TestIsInvalidPathErr): test error"),
		nil,
		fmt.Errorf("(%s/TestIsInvalidPathErr): test error", pkgName),
		errInvalidPath,
	} {
		switch errors.Is(err, errInvalidPath) {
		case true:
			assert.True(t, IsInvalidPathErr(err))

		case false:
			assert.False(t, IsInvalidPathErr(err))
		}
	}
}

func TestPathExists(t *testing.T) {
	for _, path := range []string{
		".",
		"/root",
		"/user",
		"invalid path",
	} {
		_, err := os.Stat(path)

		assert.Equal(t, err == nil, PathExists(path))
	}
}
