package packagetool

import (
	"bytes"
	"encoding/binary"
	"io"
)

type BarReader struct{}

func (BarReader) Name() string {
	return "BAR"
}

func (BarReader) Read(rd io.Reader, cb func(File) error) error {
	b := make([]byte, 256)

	// read header
	if _, err := io.ReadFull(rd, b[:12]); err != nil {
		return err
	}

	entries := binary.LittleEndian.Uint16(b[10:])
	for i := uint16(0); i < entries; i++ {
		// read filename
		if _, err := io.ReadFull(rd, b); err != nil {
			return err
		}
		nameData, _, _ := bytes.Cut(b, []byte{0})
		name := string(nameData)

		// for weird files
		if binary.LittleEndian.Uint32(b[252:]) != 3 {
			// skip
			io.ReadFull(rd, b[:4])
		}

		// read file metadata
		if _, err := io.ReadFull(rd, b[:12]); err != nil {
			return err
		}
		size := binary.LittleEndian.Uint64(b[4:])

		if err := cb(File{
			Reader:   io.LimitReader(rd, int64(size)),
			Filename: name,
		}); err != nil {
			return err
		}
	}

	return nil
}
