/*
Package json handles reading from, and writing to JSON output files
*/
package json

import (
	"encoding/json"

	"github.com/notsatan/crcgen/src/writer"
)

type jsonHandler struct{}

func (*jsonHandler) Marshal(info *writer.DirInfo) ([]byte, error) {
	return json.Marshal(info)
}

func (*jsonHandler) Unmarshal(data []byte, info *writer.DirInfo) error {
	return json.Unmarshal(data, info)
}

func (*jsonHandler) FileTypes() []string {
	return []string{"json"}
}
