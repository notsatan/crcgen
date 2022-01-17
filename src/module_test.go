package src

import (
	"testing"
)

func TestModuleName(t *testing.T) {
	if ProjectName() != "crcgen" {
		t.Errorf("Project name `%s` incorrect", ProjectName())
	}
}
