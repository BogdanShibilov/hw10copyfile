package main

import (
	"flag"
	"log"
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

	err = CopyFile(from, to, offset, limit)
	if err != nil {
		log.Panic(err)
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
