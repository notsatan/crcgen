package writer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFileTypes []string

// mockHandler is a test utility to create a mock Handler object
type mockHandler struct{}

func (*mockHandler) FileTypes() []string { return testFileTypes }

func (*mockHandler) Marshal(*DirInfo, ...bool) ([]byte, error) { return nil, nil }

func (*mockHandler) Unmarshal([]byte, *DirInfo) error { return nil }

var _ = Handler(&mockHandler{}) // verify mockHandler implements Handler

func TestAddHandler(t *testing.T) {
	reset := func() {
		outHandlers = map[string]Handler{} // isolate test case
	}

	reset()
	defer reset()

	testFileTypes = []string{"Mp4", "JsOn", " .mkv ", " ziP", "7Z", ".tXt", ".mp4"}
	AddHandler(&mockHandler{})

	var keys []string
	for key := range outHandlers {
		keys = append(keys, key)
	}

	expected := []string{"mp4", "json", "mkv", "zip", "7z", "txt"}
	assert.Equal(t, len(expected), len(keys))
	for _, ext := range keys {
		flag := false
		for _, check := range expected {
			if check == ext {
				flag = true
				break
			}
		}

		assert.Truef(t, flag, `extension "%s" not found in expected keys`, ext)
	}
}
