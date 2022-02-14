package writer

import (
	"strings"

	"github.com/notsatan/crcgen/src/logger"
)

// outHandlers maps the available Handlers to the extension they handle. Extensions
// do not contain period(s), and stored as the key in lower-case
var outHandlers map[string]Handler

/*
AddHandler registers a Handler that will be used to handle a particular output file,
these handlers are detected based on file extensions
*/
func AddHandler(handler Handler) {
	for _, ext := range handler.FileTypes() {
		ext = strings.ToLower(strings.Trim(ext, " ."))
		if _, ok := outHandlers[ext]; ok {
			logger.Debugf("(%s/AddHandler): duplicate extension `%s`", pkgName, ext)
		}

		outHandlers[ext] = handler
	}
}

/*
Handler defines a simple interface to interact with data from the output file - this
includes reading data from the file, and writing to the same file when done

Each Handler defines an object that can interact with certain file types (based on file
extensions)
*/
type Handler interface {
	// Marshal takes an object of DirInfo as input, and converts it into a byte array
	// that can be easily written to a file.
	//
	// The second argument indicates if the marshall-ed output needs to be indented
	Marshal(info *DirInfo, indent ...bool) ([]byte, error)

	// Unmarshal parses encoded data and stores the result in the DirInfo object
	Unmarshal([]byte, *DirInfo) error

	// FileTypes returns an array of strings - indicating the supported output file
	// extensions. The extensions are case-insensitive
	FileTypes() []string
}
