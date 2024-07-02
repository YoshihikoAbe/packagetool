package packagetool

import (
	"errors"
	"io"
)

var errMagic = errors.New("invalid magic")

type Reader func(io.Reader, func(File) error) error

type File struct {
	io.Reader
	Filename string
}
