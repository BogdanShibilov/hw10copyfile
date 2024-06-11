package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
)

var (
	from   string
	to     string
	offset int64
	limit  int64
)

func init() {
	flag.StringVar(&from, "from", "", "path to source file")
	flag.StringVar(&to, "to", "", "path to destination file")
	flag.Int64Var(&offset, "offset", 0, "amount of bytes which will be skipped")
	flag.Int64Var(&limit, "limit", 4096, "amount of bytes which will be copied")
}

func main() {
	flag.Parse()
	err := validateFlags()
	if err != nil {
		log.Panic(err)
	}

	sourceFile := readFile()
	defer func(sourceFile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile)

	fi, _ := sourceFile.Stat()
	if offset > fi.Size() {
		return
	}
	if offset+limit > fi.Size() {
		limit = fi.Size() - offset
	}

	destFile := createDestFile()
	defer func(destFile *os.File) {
		_ = destFile.Close()
	}(destFile)

	offset, err = sourceFile.Seek(offset, 0)
	if err != nil {
		log.Panic(err)
	}

	reader := io.LimitReader(sourceFile, limit)

	var bufSize int
	if limit > 256 {
		bufSize = 256
	} else {
		bufSize = int(limit)
	}
	buf := make([]byte, bufSize)
	var readProgress int
	for {
		r, err := reader.Read(buf)
		readProgress += r
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Panic(err)
		}
		_, err = destFile.Write(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Read %d bytes out of %d\n", readProgress, limit)
	}
}

func validateFlags() error {
	if from == "" {
		return ErrFromValueNotSpecified
	}
	if to == "" {
		return ErrToValueNotSpecified
	}
	return nil
}

func readFile() *os.File {
	sourceFile, err := os.Open(from)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Panicf("Source file %s does not exist\n", from)
		}
		if errors.Is(err, fs.ErrPermission) {
			log.Panicf("Permission denied to read file %s\n", from)
		}
		log.Panic(err)
	}
	return sourceFile
}

func createDestFile() *os.File {
	destFile, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			log.Panicf("Destination file %s already exists\n", to)
		}
		if errors.Is(err, fs.ErrPermission) {
			log.Panicf("Permission denied to write file %s\n", to)
		}
		log.Panic(err)
	}
	return destFile
}
