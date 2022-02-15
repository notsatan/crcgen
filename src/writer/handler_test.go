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

func TestValidateExt(t *testing.T) {
	reset()

	// register a bunch of mock handlers
	outHandlers = map[string]Handler{
		"mp4": &mockHandler{},
		"mkv": &mockHandler{},
		"zip": &mockHandler{},
		"mp3": &mockHandler{},
		"7z":  &mockHandler{},
	}

	for input, expected := range map[string]bool{
		"mp4":   true,
		"MKV":   true,
		"mP3":   true,
		"png":   false,
		".zip":  true,
		"  .7z": true,
		"jpeg":  false,
	} {
		res := validateExt(input)
		assert.Equalf(t, expected, res, `failed validate extension: "%s"`, input)
	}
}

func TestAddHandler(t *testing.T) {
	reset()

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
