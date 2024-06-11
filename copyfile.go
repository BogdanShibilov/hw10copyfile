package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
)

func CopyFile(source, destination string, offset, limit int64) error {
	sourceFile := readFile(source)
	defer sourceFile.Close()

	fi, _ := sourceFile.Stat()
	if offset > fi.Size() {
		return nil
	}
	if offset+limit > fi.Size() {
		limit = fi.Size() - offset
	}

	destFile := createDestFile(destination)
	defer destFile.Close()

	offset, err := sourceFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	reader := io.LimitReader(sourceFile, limit)

	return cp(reader, destFile, limit)
}

func cp(reader io.Reader, destFile *os.File, limit int64) error {
	bufSize := getBufSize(limit)
	buf := make([]byte, bufSize)
	var readProgress int

	for {
		r, err := reader.Read(buf)
		readProgress += r
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		_, err = destFile.Write(buf)
		if err != nil {
			return err
		}
		fmt.Printf("Read %d bytes out of %d\n", readProgress, limit)
	}

	return nil
}

func getBufSize(limit int64) int {
	var bufSize int
	if limit > 256 {
		bufSize = 256
	} else {
		bufSize = int(limit)
	}
	return bufSize
}

func readFile(source string) *os.File {
	sourceFile, err := os.Open(source)
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

func createDestFile(dest string) *os.File {
	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			log.Panicf("Permission denied to write file %s\n", to)
		}
		log.Panic(err)
	}
	return destFile
}
