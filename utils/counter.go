package utils

import "io"

type Counter struct {
	io.Writer
	Written int
}

func NewCounter(w io.Writer) *Counter {
	return &Counter{Writer: w}
}

func (counter *Counter) Write(p []byte) (int, error) {
	written, err := counter.Writer.Write(p)
	counter.Written += written
	return written, err
}
