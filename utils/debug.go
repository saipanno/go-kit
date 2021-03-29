package utils

import (
	"bytes"
	"encoding/json"
)

// PrettyPrint ...
func PrettyPrint(data interface{}) string {

	var out bytes.Buffer
	bs, _ := json.Marshal(data)

	json.Indent(&out, bs, "", "  ") // nolint: errcheck
	return out.String()
}
