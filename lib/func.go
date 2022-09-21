package lib

import (
	"strings"
)

var (
	formatterMap = make(map[string]Formatter)
)

func RegisterFormatter(name string, f Formatter) error {
	name = strings.TrimSpace(name)
	if _, ok := formatterMap[name]; ok {
		return ErrDuplicatedFormatter
	}
	formatterMap[name] = f
	return nil
}

func GetFormatter(name string) (Formatter, error) {
	if f, ok := formatterMap[name]; ok {
		return f, nil
	}
	return nil, ErrNotFound
}
