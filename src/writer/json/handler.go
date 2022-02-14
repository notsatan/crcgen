/*
Package json handles reading from, and writing to JSON output files
*/
package json

import (
	"encoding/json"

	"github.com/notsatan/crcgen/src/writer"
)

type jsonHandler struct{}

func (*jsonHandler) Marshal(info *writer.DirInfo, indent ...bool) ([]byte, error) {
	if len(indent) == 0 || !indent[0] {
		return json.Marshal(info) // without indents
	}

	// Indent with tabs when indentation is required
	return json.MarshalIndent(info, "", "\t")
}

func (*jsonHandler) Unmarshal(data []byte, info *writer.DirInfo) error {
	return json.Unmarshal(data, info)
}

func (*jsonHandler) FileTypes() []string {
	return []string{"json"}
}
