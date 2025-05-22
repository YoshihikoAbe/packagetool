package packagetool

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Dumper struct {
	Dest   string
	Silent bool

	buffer [1024]byte
}

func (dumper *Dumper) DumpArchive(f File) error {
	if runtime.GOOS != "windows" {
		f.Filename = strings.ReplaceAll(f.Filename, "\\", "/")
	}

	outPath := filepath.Join(dumper.Dest, f.Filename)
	outDir, _ := filepath.Split(outPath)
	if !dumper.Silent {
		fmt.Println(f.Filename, "->", outPath)
	}
	if outDir != "" {
		if err := os.MkdirAll(outDir, 0700); err != nil {
			return err
		}
	}

	wr, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer wr.Close()

	_, err = io.CopyBuffer(wr, f, dumper.buffer[:])
	return err
}
