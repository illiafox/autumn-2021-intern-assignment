package methods

import (
	"encoding/json"
	"fmt"
)

type errJSON struct {
	Ok  bool   `json:"ok"`
	Err string `json:"err,omitempty"`
}

func jsonError(format string, a ...any) []byte {
	buf, _ := json.Marshal(errJSON{
		Ok:  false,
		Err: fmt.Errorf(format, a...).Error(),
	})
	return buf
}

func jsonErrorString(err string) []byte {
	buf, _ := json.Marshal(errJSON{
		Ok:  false,
		Err: err,
	})
	return buf
}

var successBytes = []byte("{\"ok\":true}")

func jsonSuccess() []byte {
	return successBytes
}
