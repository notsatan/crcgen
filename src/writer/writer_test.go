package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func reset() {
	closeFile = (*os.File).Close
	fileIsDir = os.FileInfo.IsDir
	filePath = ""

	pathStats = os.Stat
	createFile = os.Create
	absPath = filepath.Abs
	osReadFile = os.ReadFile
	osWriteFile = os.WriteFile

	outHandlers = map[string]Handler{}
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

func TestIsPathNotWriteableErr(t *testing.T) {
	for err, expected := range map[error]bool{
		nil:                      false,
		errNotWritable:           true,
		errInvalidFile:           false,
		fmt.Errorf("test error"): false,
		fmt.Errorf(""):           false,
		fmt.Errorf("(%s): output file has invalid extension", pkgName): false,
	} {
		assert.Equalf(
			t, expected, IsPathNotWriteableErr(err),
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

func TestIsReadFileErr(t *testing.T) {
	for err, expected := range map[error]bool{
		nil:                              false,
		errReadFile:                      true,
		errInvalidFile:                   false,
		errors.Wrap(errReadFile, "test"): true,
		fmt.Errorf("test error"):         false,
		fmt.Errorf(""):                   false,
		fmt.Errorf("(%s): output file has invalid extension", pkgName): false,
	} {
		assert.Equalf(
			t, expected, IsReadFileErr(err),
			`failed to match "%v" -> "%v"`, expected, err,
		)
	}
}

func TestIsPathDirErr(t *testing.T) {
	for err, expected := range map[error]bool{
		nil:                      false,
		errNotWritable:           false,
		errPathIsDir:             true,
		fmt.Errorf("test error"): false,
		fmt.Errorf(""):           false,
		fmt.Errorf("(%s): output file has invalid extension", pkgName): false,
	} {
		assert.Equalf(
			t, expected, IsPathDirErr(err),
			`failed to match "%v" -> "%v"`, expected, err,
		)
	}
}

func TestIsHandlerNotFoundErr(t *testing.T) {
	for err, expected := range map[error]bool{
		nil:                      false,
		errNoHandler:             true,
		errPathIsDir:             false,
		fmt.Errorf("test error"): false,
		fmt.Errorf(""):           false,
		fmt.Errorf("(%s): output file has invalid extension", pkgName): false,
	} {
		assert.Equalf(
			t, expected, IsHandlerNotFoundErr(err),
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

	pathStats = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	fileIsDir = func(os.FileInfo) bool { return false }
	createFile = func(string) (*os.File, error) { return nil, nil }
	osReadFile = func(string) ([]byte, error) { return []byte{}, nil }

	outHandlers = map[string]Handler{"json": &mockHandler{}} // mock add `json` handler
	validInput := "output.json"

	// expect a failure if file fails to close
	closeFile = func(*os.File) error { return fmt.Errorf("(%s): test error", pkgName) }
	assert.Errorf(t, start(validInput), `no error returned when file failed to close`)

	closeFile = func(*os.File) error { return nil } // path to undo this
	err := start(validInput)
	assert.NoErrorf(t, err, "unexpected error: %v", err)

	// expect a failure if the output file cannot be read
	err = fmt.Errorf("(%s): test error", pkgName)
	osReadFile = func(string) ([]byte, error) { return nil, err }

	err = start(validInput)
	assert.Errorf(t, err, "unexpected non-failure: %v", err)
}

func TestFixPath(t *testing.T) {
	reset()

	outHandlers = map[string]Handler{
		"json": &mockHandler{},
		"yaml": &mockHandler{},
		"yml":  &mockHandler{},
	}

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
	outHandlers = map[string]Handler{"json": &mockHandler{}}

	path, err := fixPath("/test/path.json")

	assert.Error(t, err, "expected an error when absolute path can't be formed")
	assert.Emptyf(t, path, `expected path to be empty for an error: "%v"`, path)
	assert.Truef(t, IsAbsPathErr(err), `expected custom error type: "%v"`, err)
}

func TestCreateOutFile_PathStats(t *testing.T) {
	reset()

	// Mock functions to isolate scenarios
	createFile = func(string) (*os.File, error) { return nil, nil }
	closeFile = func(*os.File) error { return nil }

	// Wrapper to run uniform tests
	runner := func() error { return createOutFile("") }

	// If file does not exit, the function should pass normally
	pathStats = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	assert.NoError(t, runner(), `unexpected error for mock #1`)

	// Expect an error if file stats can't be read
	pathStats = func(string) (os.FileInfo, error) { return nil, fmt.Errorf("") }
	assert.Error(t, runner(), `no error occurred for mock #3`)

	// If a directory exists at the same path, expect an error
	fileIsDir = func(os.FileInfo) bool { return true }
	pathStats = func(string) (os.FileInfo, error) { return nil, nil }

	result := runner()
	assert.Error(t, result, `no error occurred for mock #2`)
	assert.Truef(t, IsPathDirErr(result), `failed to match error type: "%v"`, result)

	fileIsDir = os.FileInfo.IsDir // reset

	// Mock scenario when output file already exists - expect a direct return
	fileIsDir = func(os.FileInfo) bool { return false }
	assert.NoError(t, runner(), `unexpected error for mock #4`)
}

func TestCreateOutFile_CreateFile(t *testing.T) {
	reset()

	// Mocks to isolate the portion being tested
	pathStats = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	fileIsDir = func(os.FileInfo) bool { return false }

	runner := func() error { return createOutFile("") }

	// Expect an error indicating un-writable path if the file can't be created
	createFile = func(string) (*os.File, error) { return nil, fmt.Errorf("test") }
	result := runner()
	assert.Error(t, result, `no error occurred for mock #1`)
	assert.Truef(
		t, IsPathNotWriteableErr(result), `failed to match error type: "%v"`, result,
	)

	// Mock normal conditions - no error should occur
	closeFile = func(*os.File) error { return nil }
	createFile = func(string) (*os.File, error) { return nil, nil }
	assert.NoError(t, runner())
}

func TestReadFile(t *testing.T) {
	reset()

	// Ensure `readFile` fails in case of an error, and vice-versa
	osReadFile = func(string) ([]byte, error) { return nil, errReadFile }
	assert.Error(t, readFile(&DirInfo{}))

	// Ensure direct return in case the file is empty
	osReadFile = func(string) ([]byte, error) { return []byte{}, nil }
	assert.NoError(t, readFile(&DirInfo{}))

	// Ensure an error is returned if no handler can interact with the output file
	osReadFile = func(string) ([]byte, error) { return []byte{15}, nil }
	filePath = "/path/to/file.mp4"
	assert.Error(t, readFile(&DirInfo{}))

	// Ensure error is returned if unmarshal fails
	outHandlers["yaml"] = &mockHandlerFail{}
	filePath = "output.yaml"
	assert.True(t, IsInvalidFileErr(readFile(&DirInfo{})))

	// No error should be returned for a successful run
	outHandlers = map[string]Handler{"yml": &mockHandler{}}
	filePath = "output.yml"
	assert.NoError(t, readFile(&DirInfo{}))
}

// mockHandlerFail is a wrapper over mockHandler where all methods fail - when possible
type mockHandlerFail struct {
	mockHandler
}

func (*mockHandlerFail) FileTypes() []string { return []string{} }

func (*mockHandlerFail) Marshal(*DirInfo, ...bool) ([]byte, error) {
	return nil, fmt.Errorf("(%s/mockHandlerFail.Marshal): test error", pkgName)
}

func (*mockHandlerFail) Unmarshal([]byte, *DirInfo) error {
	return fmt.Errorf("(%s/mockHandlerFail.Unmarshal): test error", pkgName)
}

func TestWrite(t *testing.T) {
	reset()

	osWriteFile = func(string, []byte, os.FileMode) error {
		return fmt.Errorf("(%s/TestWrite): mock error", pkgName) // mock failure
	}

	// Ensure an error is returned if no handler is available for the filetype
	filePath = "/path/to/incorrect-file.mp4"
	assert.True(t, IsHandlerNotFoundErr(Write(&DirInfo{})))

	// Mock a JSON file, and set a handler that handles JSON file to fail
	filePath = "/path/to/configs.json"
	outHandlers = map[string]Handler{"json": &mockHandlerFail{}}
	assert.Error(t, Write(&DirInfo{})) // expect an error when a handler fails

	// Ensure error is returned when file cannot be written to
	filePath = "/path/to/configs.yml"
	outHandlers["yml"] = &mockHandler{}

	// Error returned since writing will fail
	assert.True(t, IsPathNotWriteableErr(Write(&DirInfo{})))

	osWriteFile = func(string, []byte, os.FileMode) error { return nil } // mock success
	assert.NoError(t, Write(&DirInfo{}))
}
