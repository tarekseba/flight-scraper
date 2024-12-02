package utils

import (
	"fmt"
	"runtime"
)

type Once[T any] struct {
	action func() T
	val    T
	done   bool
}

func NewOnce[T any](action func() T) Once[T] {
	return Once[T]{
		action: action,
		done:   false,
	}
}

func (o *Once[T]) Compute() T {
	if o.done {
		return o.val
	}
	o.done = true
	o.val = o.action()
	return o.val
}

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
