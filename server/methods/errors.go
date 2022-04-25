package methods

import (
	"encoding/json"
	"fmt"
	"io"
)

type errJSON struct {
	Ok  bool   `json:"ok"`
	Err string `json:"err,omitempty"`
}

func jsonError(writer io.Writer, format string, a ...interface{}) (int, error) {
	str := fmt.Errorf(format, a...).Error()

	buf, err := json.Marshal(errJSON{
		Ok:  false,
		Err: str,
	})

	if err != nil {
		return writer.Write([]byte(str))
	}

	return writer.Write(buf)
}

func jsonErrorString(writer io.Writer, str string) (int, error) {
	buf, err := json.Marshal(errJSON{
		Ok:  false,
		Err: str,
	})

	if err != nil {
		return writer.Write([]byte(str))
	}

	return writer.Write(buf)
}

var successBytes = []byte("{\"ok\":true}")

func jsonSuccess(writer io.Writer) (int, error) {
	return writer.Write(successBytes)
}
