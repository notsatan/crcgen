package writer

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
