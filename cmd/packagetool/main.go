package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/YoshihikoAbe/packagetool"
)

var (
	out  string
	list bool
)

func main() {
	flag.StringVar(&out, "o", "./", "Path to the ouput directory")
	flag.BoolVar(&list, "l", false, "List archive contents")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] FILENAME\nList of available options:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(filename)
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	pr, err := packagetool.DetectArchiveType(rd)
	if err != nil {
		fatal("failed to determine archive type:", err)
	}
	fmt.Println("archive type:", pr.Name())

	callback := dumpArchive
	if list {
		callback = listArchive
	}

	start := time.Now()
	if err := pr.Read(rd, callback); err != nil {
		fatal(err)
	}
	fmt.Println("time elapsed:", time.Now().Sub(start))
}

func listArchive(f packagetool.File) error {
	fmt.Println(f.Filename)
	io.Copy(io.Discard, f)
	return nil
}

func dumpArchive(f packagetool.File) error {
	if runtime.GOOS != "windows" {
		f.Filename = strings.ReplaceAll(f.Filename, "\\", "/")
	}

	path := filepath.Join(out, f.Filename)
	dir, _ := filepath.Split(path)
	fmt.Println(f.Filename, "->", path)

	if dir != "" {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	wr, err := os.Create(path)
	if err != nil {
		return err
	}
	defer wr.Close()

	_, err = io.Copy(wr, f)
	return err
}

func fatal(v ...any) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}
