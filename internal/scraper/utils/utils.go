package utils

import (
	"fmt"
	"runtime"
)

func AnnotateError(err error) error {
	if err == nil {
		return err
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}
	err = fmt.Errorf("[%s:%d]: %w", file, line, err)
	return err
}
