package main

import "errors"

var (
	ErrFromValueNotSpecified = errors.New("from value not specified")
	ErrToValueNotSpecified   = errors.New("to value not specified")
)
