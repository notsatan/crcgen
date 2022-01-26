package logger

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// openFixtureCalled indicates if openFixture was called
var openFixtureCalled = false

// openFixture is a test fixture to overwrite openLogFile lambda - returns empty
// response, modifies value of openFixtureCalled to indicate the function was called
var openFixture = func(...string) (zapcore.WriteSyncer, func(), error) {
	openFixtureCalled = true
	return nil, func() {}, nil
}

// reset rolls back values of variables/lambdas to their defaults
func reset() {
	openLogFile = zap.Open
	sync = (*zap.SugaredLogger).Sync

	openFixtureCalled = false
	streamCloser = func() {}
	logLevel.SetLevel(zap.WarnLevel)
}

func TestLog_LogToStderr(t *testing.T) {
	// Ensures a call to `Log` modifies the log level without opening the log file

	defer reset()

	openLogFile = openFixture // replace actual function with the fixture
	assert.NoError(t, Log(false))

	assert.Equal(t, zap.DebugLevel, logLevel.Level()) // log level should be modified
	assert.False(
		t, openFixtureCalled, "unexpected attempt to open the log file",
	)
}

func TestLog_LogToFile(t *testing.T) {
	// Ensure logging to a file opens the log file, and modifies the log level

	defer reset()

	openLogFile = openFixture
	assert.NoError(t, Log(true)) // logs written to file

	assert.Equal(t, zap.DebugLevel, logLevel.Level())
	assert.True(
		t, openFixtureCalled, "no attempt made to open the log file",
	)
}

func TestLog_OpenFileError(t *testing.T) {
	// Tests the scenario when opening a log file results in an error - log level should
	// be modified, and the function should return an error

	defer reset()

	// Modify `openLogFile` lambda to return an error
	openLogFile = func(...string) (zapcore.WriteSyncer, func(), error) {
		return nil, func() {}, fmt.Errorf(
			"(%s/TestLog_OpenFileError): test error", pkgName,
		)
	}

	assert.Error(t, Log(true), "`logger.Log` failed to return error during failure")
	assert.Equal(t, zap.DebugLevel, logLevel.Level())
}

func TestStop_NoError(t *testing.T) {
	// No error should occur when `sync` does not fail

	defer reset()

	sync = func(*zap.SugaredLogger) error { return nil }
	assert.NoError(t, Stop(), "unexpected error when stopping the logger")
}

func TestStop_LinuxError(t *testing.T) {
	// Failure to sync the logger is expected on linux - error should be suppressed

	defer reset()

	sync = func(*zap.SugaredLogger) error { return errors.New("/dev/stderr") }

	// Error would be suppressed only for tests running on Linux
	if runtime.GOOS == "linux" {
		assert.NoError(t, Stop(), "failed to suppress logger sync error on linux")
	} else {
		assert.Error(t, Stop(), "logger didn't return error after failing to sync")
	}
}

func TestStop_SyncError(t *testing.T) {
	// If `Stop` fails, an error should be returned

	defer reset()

	sync = func(*zap.SugaredLogger) error {
		return fmt.Errorf("(%s/TestStop_SyncError): test error", pkgName)
	}

	assert.Error(t, Stop(), "logger did not propagate error")
}

func TestPublicMethods(t *testing.T) {
	//nolint:godox // suppress error for the `to-do` -- below low priority
	// TODO: Write this test better
	normal := []func(...interface{}){
		Debug,
		Info,
		Warn,
		Error,
	}

	formatted := []func(string, ...interface{}){
		Debugf,
		Infof,
		Warnf,
		Errorf,
	}

	for _, function := range normal {
		function("test")
	}

	for _, function := range formatted {
		function("test")
	}
}
