package lib

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func reset() {
	filepathWalk = filepath.Walk
}

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

func TestWalkPath_PathError(t *testing.T) {
	reset()

	filepathWalk = func(_ string, _ filepath.WalkFunc) error {
		return nil
	}

	walkFunc := filepath.WalkFunc(nil)
	for _, path := range []string{
		".",
		"../..",
		"<invalid path>",
		"/root",
	} {
		if PathExists(path) {
			assert.NoError(t, WalkPath(path, walkFunc))
		} else {
			err := WalkPath(path, walkFunc)

			assert.Error(t, err)
			assert.True(t, IsInvalidPathErr(err))
		}
	}
}

func TestWalkPath_WalkError(t *testing.T) {
	reset()

	// Failure when walking the path should return an error
	filepathWalk = func(_ string, _ filepath.WalkFunc) error {
		return fmt.Errorf("(%s/TestWalk_Error): test error", pkgName)
	}

	assert.Error(t, WalkPath(".", nil))
}

func TestWalkPath(t *testing.T) {
	reset()

	walkFunc := func(path string, info fs.FileInfo, err error) error {
		// `walkFunc` should be run on non-empty path, without any errors, and
		// specifically on files - i.e. no directories
		assert.NotEmptyf(t, path, `walkFunc run on empty path: "%s"`, path)
		assert.NoErrorf(t, err, `walkFunc run on path with error: "%s"`, path)
		assert.Falsef(t, info.IsDir(), `walkFunc run on directory: "%s"`, path)
		return nil
	}

	assert.NoError(t, WalkPath("../..", walkFunc))
}
