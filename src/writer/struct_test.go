package writer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirInfo_Name(t *testing.T) {
	for input, expected := range map[string]string{
		"/tmp/path/dest/": "dest",
		"/destination":    "destination",
		"/test/output/":   "output",
		"/tmp/input.jpg/": "input.jpg",
	} {
		obj := DirInfo{Path: input}

		name := obj.Name()
		assert.Equalf(t, expected, name, `(input, output): ("%s", "%s")`, input, name)
	}
}

func TestDirInfo_CalcModTime_DirectReturn(t *testing.T) {
	// Ensure the function performs a direct return if last mod time is specified, and
	// the result should be an int64
	obj := DirInfo{
		LastMod: 10,
		Files: []FileInfo{
			{LastMod: 25},
		},
	}

	// The result should be 25, but direct return means value returned would be 10
	time := obj.CalcModTime()
	assert.Equalf(t, int64(10), time, "result returned: %v", time)
}

func TestDirInfo_CalcModTime_FileModTime(t *testing.T) {
	// Ensure mod time is being picked from files
	obj := DirInfo{
		LastMod: 0, // ensures no direct return
		Files: []FileInfo{
			{LastMod: 2},
			{LastMod: 13},
			{LastMod: 30},
			{LastMod: 32},
		},
	}
	time := obj.CalcModTime()
	assert.Equalf(t, int64(32), time, "result returned: %v", time)
}

func TestDirInfo_CalcModTime_DirModTime(t *testing.T) {
	// Ensure mod time is being picked from directories
	obj := DirInfo{
		LastMod: 0, // no direct return
		Dirs: []DirInfo{
			{LastMod: 15},
			{LastMod: 4},
			{LastMod: 3},
			{LastMod: 210},
		},
	}
	time := obj.CalcModTime()
	assert.Equalf(t, int64(210), time, "result returned: %v", time)
}

func TestDirInfo_CalcModTime(t *testing.T) {
	// Ensure mod time is being picked from files
	obj := DirInfo{
		LastMod: 0, // ensures no direct return
		Files: []FileInfo{
			{LastMod: 13},
			{LastMod: 2},
			{LastMod: 11},
			{LastMod: 29},
		},

		// Directories that contain a direct
		Dirs: []DirInfo{
			{Files: []FileInfo{{LastMod: 10}, {LastMod: 14}}},
			{LastMod: 0o3},
			{Files: []FileInfo{{LastMod: 10}, {LastMod: 14}, {LastMod: 28}}},
			{LastMod: 19, Files: []FileInfo{{LastMod: 1000}}}, // 19 - direct return
		},
	}

	// Even though a file contains mod time of `1000`, it won't be the answer because
	// the parent contains a value for `LastMod` which will be returned by `CalcModTime`
	time := obj.CalcModTime()
	assert.Equalf(t, int64(29), time, "result returned: %v", time)
}

func TestNewDir(t *testing.T) {
	// Ensure path is combined when both `dirName` and `parentPath` are supplied, and
	// vice-versa. The rest of the method is linear logic!

	for expected, input := range map[string]struct {
		dirName    string
		parentPath string
	}{
		// 		in/put 		: 			output
		"/path/to/directory": {parentPath: "/path/to", dirName: "directory"},
		"/obj/directory":     {parentPath: "/obj", dirName: "directory"},
		"/dirName":           {parentPath: "/dirName"},
	} {
		obj := NewDir(input.dirName, input.parentPath, nil, nil, 0)

		assert.Equalf(
			t, expected, obj.Path, `(input, output): ("%s", "%s")`, input, obj.Path,
		)
	}
}
