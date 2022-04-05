package multiwriter

import "io"

type MultiWriter struct {
	Writers []io.Writer
}

func (l MultiWriter) Write(p []byte) (int, error) {
	n := len(p)
	for i := range l.Writers {
		buf, err := l.Writers[i].Write(p)
		if err != nil || buf != n {
			return buf, err
		}
	}
	return n, nil
}

func New(writers ...io.Writer) MultiWriter {
	return MultiWriter{
		Writers: writers,
	}
}
