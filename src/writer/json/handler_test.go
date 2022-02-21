package json

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/notsatan/crcgen/src/writer"
)

func handler() *jsonHandler {
	return &jsonHandler{}
}

func TestJsonHandler_Marshal(t *testing.T) {
	input := &writer.DirInfo{
		Path:    "/test/path",
		LastMod: 0,
		Files: []writer.FileInfo{
			{Path: "/path/to/test/file.json", Size: 300},
		},
	}

	// Supposed to be the JSON equivalent of the `input` object
	rawJSON := `{
	"Path": "/test/path",
	"Dirs": null,
	"Files": [
		{
			"Path": "/path/to/test/file.json",
			"Checksums": {
				"CRC32": ""
			},
			"Size": 300,
			"LastMod": 0
		}
	],
	"LastMod": 0
}`

	var output bytes.Buffer // strip extra spaces
	require.NoError(t, json.Compact(&output, []byte(rawJSON)))

	// Without indentation
	res, err := handler().Marshal(input)
	require.NoError(t, err)
	assert.Equal(t, output.Bytes(), res)

	res, err = handler().Marshal(input, true)
	require.NoError(t, err)
	assert.Equal(t, []byte(rawJSON), res)
}

func TestJsonHandler_Unmarshal(t *testing.T) {
	in := `
{
  "Path": "/test/path",
  "Dirs": null,
  "Files": [
    {
      "Path": "/path/to/test/file.json",
      "Checksums": {
        "CRC32": ""
      },
      "Size": 300,
      "LastMod": 0
    }
  ],
  "LastMod": 0
}
`

	// Should be the struct-equivalent of the JSON string above
	output := &writer.DirInfo{
		Path:    "/test/path",
		LastMod: 0,
		Files: []writer.FileInfo{
			{Path: "/path/to/test/file.json", Size: 300},
		},
	}

	var input bytes.Buffer // strip extra spacing from input JSON
	require.NoError(t, json.Compact(&input, []byte(in)))

	var info writer.DirInfo
	require.NoError(t, handler().Unmarshal(input.Bytes(), &info))

	assert.Equal(t, output, &info) // compare
}

func TestJsonHandler_FileTypes(t *testing.T) {
	assert.Equal(t, []string{"json"}, handler().FileTypes())
}
