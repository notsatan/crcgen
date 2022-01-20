package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/notsatan/crcgen/src/logger"
)

const (
	envProd  = "production_mode"
	envDebug = "debug_mode"
	envBuild = "__BUILD_MODE__"
)

func resetEnv() {
	// Unset all environment variables
	_ = os.Unsetenv(envProd)
	_ = os.Unsetenv(envDebug)
	_ = os.Unsetenv(envBuild)
}

func reset() {
	exit = os.Exit
	initLogger = logger.Log
	closeLogger = logger.Stop
}

func TestMain(m *testing.M) {
	// Run all tests, unset env variables, and exit
	resetEnv()

	// TODO: Figure out a better way than running tests in debug mode
	_ = os.Setenv(debugMode, "debug")

	val := m.Run()
	resetEnv()
	os.Exit(val)
}

func TestMainMethod(_ *testing.T) {
	main()
}

func TestCrashOnErr(t *testing.T) {
	defer reset()

	calls := 0
	exit = func(int) { calls++ } // increment `calls` when `exit` is run by `CrashOnErr`

	crashOnErr(nil, "") // `exit` should not be called without an error
	assert.Equal(t, calls, 0)

	crashOnErr(errors.New("test"), "test message")
	assert.Equal(t, calls, 1)
}

func TestCloseRes(t *testing.T) {
	defer reset()

	// `closeRes` is designed to run as an independent function - it should consume
	// errors internally

	closeLogger = func() error { return nil }
	assert.NotPanics(t, func() { closeRes() })

	closeLogger = func() error { return fmt.Errorf("(%s/main): test error", pkgName) }
	assert.NotPanics(t, func() { closeRes() })
}
