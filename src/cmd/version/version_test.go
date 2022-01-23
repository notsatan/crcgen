package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionString(t *testing.T) {
	assert.Condition(
		t, func() bool { return strings.Trim(version, "\n ") == Get() },
		"`version` string is not trimmed",
	)
}
