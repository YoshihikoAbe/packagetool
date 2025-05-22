package packagetool

import (
	"bufio"
	"errors"
	"io"
)

var errMagic = errors.New("invalid magic")

type File struct {
	io.Reader
	Filename string
}

func (f *File) Skip() {
	io.Copy(io.Discard, f)
}

type Reader interface {
	io.Reader
	io.ByteReader
}

type PackageReader interface {
	Name() string
	Read(rd Reader, cb func(File) error) error
}

func DetectArchiveType(rd *bufio.Reader) (PackageReader, error) {
	magic, err := rd.Peek(3)
	if err != nil {
		return nil, err
	}

	var pr PackageReader
	switch string(magic) {
	case "QAR":
		pr = QarReader{}
	case "MAS":
		pr = &MarReader{}
	default:
		pr = BarReader{}
	}
	return pr, nil
}
