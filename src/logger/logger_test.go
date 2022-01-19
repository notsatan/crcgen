package logger

import "testing"

func TestLogger(t *testing.T) {
	zapLogger.Debug("Logger is running")
}
