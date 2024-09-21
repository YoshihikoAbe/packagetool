package packagetool

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type MarReader struct{}

func (MarReader) Name() string {
	return "MAR"
}

func (MarReader) Read(rd io.Reader, cb func(File) error) error {
	b := make([]byte, 128)

	// read magic
	magic := b[:8]
	io.ReadFull(rd, magic)
	if string(magic) != "MASMAR0\x00" {
		return errMagic
	}

	for {
		t, err := rd.(io.ByteReader).ReadByte()
		if err != nil {
			return err
		}

		// EOF
		if t == 255 {
			return nil
		}

		// read filename
		size, err := readNullTerminated(b, rd.(io.ByteReader))
		if err != nil {
			return err
		}
		filename := string(b[:size])

		if t == 1 {
			// read file size
			if _, err := io.ReadFull(rd, b[:4]); err != nil {
				return err
			}
			size := binary.LittleEndian.Uint32(b)

			if err := cb(File{
				Reader:   io.LimitReader(rd, int64(size)),
				Filename: filename,
			}); err != nil {
				return err
			}
		} else if t != 2 {
			return fmt.Errorf("invalid file type: %d", t)
		}
	}
}

func readNullTerminated(out []byte, rd io.ByteReader) (int, error) {
	for i := range out {
		b, err := rd.ReadByte()
		if err != nil {
			return 0, err
		}

		if b == 0 {
			return i, nil
		}
		out[i] = b
	}
	return 0, errors.New("filename too long")
}
