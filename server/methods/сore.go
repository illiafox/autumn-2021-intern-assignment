package methods

import (
	"encoding/json"
	"io"

	"autumn-2021-intern-assignment/database/model"
)

type Methods struct {
	db model.Repository
}

func New(db model.Repository) *Methods {
	return &Methods{db: db}
}

type errJSON struct {
	Ok  bool   `json:"ok"`
	Err string `json:"err,omitempty"`
}

func EncodeError(writer io.Writer, err error) (int, error) {
	str := err.Error()

	buf, err := json.Marshal(errJSON{
		Ok:  false,
		Err: str,
	})

	if err != nil {
		return writer.Write([]byte(str))
	}

	return writer.Write(buf)
}

func EncodeString(writer io.Writer, str string) (int, error) {
	buf, err := json.Marshal(errJSON{
		Ok:  false,
		Err: str,
	})

	if err != nil {
		return writer.Write([]byte(str))
	}

	return writer.Write(buf)
}

var success = []byte("{\"ok\":true}")
