/*
Package logger exposes a simple logging interface for the application. The package is
responsible for writing all application logs, which, by default, are written at
to `stderr` for events at warn level or above.

The behavior of the logger can be modified using the `logger.Log()` function to write
all logs at debug level or above to `stderr`, and a text file named `logs.txt`

To safely close the logger, use `logger.Stop()`
*/
package logger

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	pkgName = "logger"

	osLinux = "linux"
)

var (
	// logWriter is responsible for writing logs to one or more output streams
	logWriter = zapcore.WriteSyncer(os.Stderr) // no logs written by default

	// logLevel indicates the minimum log level required for an event to be logged
	logLevel = zap.NewAtomicLevelAt(zap.WarnLevel)

	// streamCloser is used to close the log stream (when writing logs to a file), gets
	// overwritten by an actual function when a file is opened in `logger.Log`
	streamCloser = func() {}
)

var (
	// openLogFile maps the zap.Open function
	openLogFile = zap.Open

	// sync attempts to sync an instance of zap.SugaredLogger
	sync = (*zap.SugaredLogger).Sync
)

/*
config simply defines the logger configuration
*/
var config = zapcore.EncoderConfig{
	// Empty value removes the field from logs
	TimeKey:       "set",
	NameKey:       "set",
	LevelKey:      "set",
	MessageKey:    "set",
	CallerKey:     "set",
	StacktraceKey: "set",

	// Defaults to a new-line, good enough
	LineEnding: zapcore.DefaultLineEnding,

	// Modify log level tag to be `[INFO]` instead of `INFO`
	EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	},

	// Custom date-time format for log messages
	EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("02:01:06::3:04:05 PM")) // dd/mm/yy::<time>
	},
}

// zapLogger is the central logger used by the package to write logs, logLevel decides
// the minimum log level required for an event to be logged, while logWriter decides
// the streams where the logs are written (file, stderr, etc.)
var zapLogger = zap.New(
	zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		logWriter,
		logLevel,
	),
).Sugar()

func init() {
	streamCloser() // might as well just hit 100% coverage	¯\_(ツ)_/¯
}

/*
Log modifies the logger to log events at or above zap.DebugLevel, and write logs to
a file
*/
func Log(writeToFile bool) (err error) {
	logLevel.SetLevel(zap.DebugLevel)

	if writeToFile {
		logWriter, streamCloser, err = openLogFile([]string{
			// Write logs to `stderr` and a file named `logs.txt`
			"stderr", "logs.txt",
		}...)
	}

	return errors.Wrapf(err, "(%s/Log)", pkgName)
}

/*
Stop closes the logger gracefully
*/
func Stop() error {
	defer streamCloser()

	// detectLinuxErr is a simple lambda to detect errors arising when the `sync` method
	// is called on a logger writing to stderr/stdout in Linux
	// Related to: https://github.com/uber-go/zap/issues/880
	detectLinuxErr := func(err error) bool {
		return runtime.GOOS == osLinux && strings.Contains(err.Error(), "/dev/stderr")
	}

	// Attempt to sync - if an error arises, check if OS is linux. If yes, ignore!
	err := sync(zapLogger)
	if err != nil && detectLinuxErr(err) {
		err = nil
	}

	return errors.Wrapf(err, "(%s/Stop)", pkgName)
}

/*********************************** Boilerplate **************************************/

/*
Debug writes log messages to the `debug` level
*/
func Debug(args ...interface{}) {
	zapLogger.Debug(args...)
}

/*
Debugf uses `fmt.Sprintf` style formatting to log messages at the `debug` level
*/
func Debugf(message string, args ...interface{}) {
	zapLogger.Debugf(message, args...)
}

/*
Info writes log messages at the `info` level
*/
func Info(args ...interface{}) {
	zapLogger.Info(args...)
}

/*
Infof uses `fmt.Sprintf` style formatting to log messages at the `info` level
*/
func Infof(template string, args ...interface{}) {
	zapLogger.Infof(template, args...)
}

/*
Warn writes log messages at the `warn` level
*/
func Warn(args ...interface{}) {
	zapLogger.Warn(args...)
}

/*
Warnf uses `fmt.Sprintf` style formatting to log messages at the `warn` level
*/
func Warnf(msg string, args ...interface{}) {
	zapLogger.Warnf(msg, args...)
}

/*
Error writes log messages at the `error` level
*/
func Error(args ...interface{}) {
	zapLogger.Error(args...)
}

/*
Errorf uses `fmt.Sprintf` style formatting to log messages at the `error` level
*/
func Errorf(template string, args ...interface{}) {
	zapLogger.Errorf(template, args...)
}
