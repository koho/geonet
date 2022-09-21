package lib

import "errors"

var (
	ErrDuplicatedFormatter = errors.New("duplicated formatter")
	ErrNotImplemented      = errors.New("method not implemented")
	ErrNotFound            = errors.New("formatter not found")
)
