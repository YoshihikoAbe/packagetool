package packagetool

import (
	"bytes"
	"encoding/binary"
	"io"
)

type QarReader struct{}

func (QarReader) Name() string {
	return "QAR"
}

func (QarReader) Read(rd Reader, cb func(File) error) error {
	b := make([]byte, 144)

	// read header
	if _, err := io.ReadFull(rd, b[:8]); err != nil {
		return nil
	}
	if string(b[:4]) != "QAR\x00" {
		return errMagic
	}
	entires := binary.LittleEndian.Uint32(b[4:])

	for i := uint32(0); i < entires; i++ {
		// read file metadata
		if _, err := io.ReadFull(rd, b); err != nil {
			return nil
		}
		size := int64(binary.LittleEndian.Uint64(b[136:]))

		if err := cb(File{
			Reader:   io.LimitReader(rd, size),
			Filename: string(bytes.TrimRight(b[:128], "\x00")),
		}); err != nil {
			return err
		}
	}

	return nil
}
