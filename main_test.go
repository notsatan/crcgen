package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/notsatan/crcgen/src/cmd"
)

func reset() {
	run = cmd.Run

	exit = os.Exit
}

func TestMain_Fail(t *testing.T) {
	// Should call `exit` if an error is encountered

	reset()

	calls := 0
	exit = func(int) { calls++ }
	run = func() error { return fmt.Errorf("(%s/TestMain_Fail): test error", pkgName) }

	main()
	assert.Equal(t, 1, calls, "failed to force close on error")
}
