/*
Package logger exposes a simple logging interface for the application
*/
package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// logWriter is responsible for writing logs to one or more output streams
	logWriter = zapcore.WriteSyncer(os.Stderr) // no logs written by default

	// logLevel indicates the minimum log level required for an event to be logged
	logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
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
	zapLogger.Info("Logger is running")
}
