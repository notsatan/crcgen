package main

import (
	"os"
	"testing"
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

func TestMain(m *testing.M) {
	// Run all tests, unset env variables, and exit
	resetEnv()
	val := m.Run()
	resetEnv()
	os.Exit(val)
}

func TestMainMethod(_ *testing.T) {
	main()
}
