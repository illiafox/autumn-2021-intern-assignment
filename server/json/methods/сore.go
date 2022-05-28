package methods

import (
	"encoding/json"
	"io"

	rep "autumn-2021-intern-assignment/database/model"
)

type Methods struct {
	db rep.Repository
}

func New(db rep.Repository) *Methods {
	return &Methods{db: db}
}

// // //

// Error
// @Description Api error
type Error struct {
	Ok  bool   `json:"ok" default:"false"`
	Err string `json:"err"`
}

func WriteError(writer io.Writer, err error) (int, error) {
	str := err.Error()

	buf, err := json.Marshal(Error{
		Ok:  false,
		Err: str,
	})

	if err != nil {
		return writer.Write([]byte(str))
	}

	return writer.Write(buf)
}

func WriteString(writer io.Writer, str string) (int, error) {
	buf, err := json.Marshal(Error{
		Ok:  false,
		Err: str,
	})

	if err != nil {
		return writer.Write([]byte(str))
	}

	return writer.Write(buf)
}

var success = []byte("{\"ok\":true}")
