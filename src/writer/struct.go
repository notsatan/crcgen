package writer

import "path/filepath"

/*
Checksums contains the various checksums generated for files
*/
type Checksums struct {
	CRC32 string
}

/*
FileInfo defines the structure of the output for an individual file - this will be
written to the output file
*/
type FileInfo struct {
	// Name is the file name
	Name string

	// Checksums contains checksums generated for the file
	Checksums Checksums

	// Size contains the file size, in bytes - not intended to be human-readable
	Size int64

	// LastMod indicates time when the file was last modified. Represents epoch time,
	// not intended to be human-readable
	LastMod int64
}

/*
DirInfo defines contents of a directory. Each directory can contain multiple
directories, and files
*/
type DirInfo struct {
	// Path contains the full path to the directory
	Path string

	// Dirs maps all the directories present in this directory as DirInfo objects
	Dirs []DirInfo

	// Files maps all the files present in the directory as FileInfo objects
	Files []FileInfo

	// LastMod indicates the time when the directory was last modified. Represents epoch
	// time, not intended to be human-readable
	LastMod int64
}

/*
Name returns the name of the directory
*/
func (dir *DirInfo) Name() string {
	return filepath.Base(dir.Path)
}

/*
CalcModTime calculates the LastMod time for a directory, in case this has already been
calculated, the previous value is directly returned

For directories without a set value of LastMod time, this method will iterate over each
file and directory, fetching the last mod time for each, and setting the greatest value
as the last modification time for this directory

Note: For worst-case scenario, this method ends up being recursive -- a call is made
to CalcModTime for each directory in Dirs
*/
func (dir *DirInfo) CalcModTime() int64 {
	if dir.LastMod > 0 {
		return dir.LastMod
	}

	// To calculate the last mod time, iterate over each file and directory within the
	// current directory. Since last mod is the epoch time when the item was last
	// modified, the greatest value of last mod time will be the modification time
	// for this directory
	var modTime int64

	for i := range dir.Files {
		if dir.Files[i].LastMod > modTime {
			modTime = dir.Files[i].LastMod
		}
	}

	for i := range dir.Dirs {
		if time := dir.Dirs[i].CalcModTime(); time > modTime {
			modTime = time
		}
	}

	dir.LastMod = modTime // save this value for future use
	return modTime
}
