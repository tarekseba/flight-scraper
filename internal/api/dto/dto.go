package dto

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Response[T any] struct {
	Data    *T     `json:"data"`
	Err     string `json:"err"`
	Message string `json:"message"`
}

func NewResponse[T any](data *T, err string, message string) Response[T] {
	return Response[T]{
		Data:    data,
		Err:     err,
		Message: message,
	}
}

func JSON[T any](payload *Response[T]) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(payload)
	return buf, err
}

func HandleResponse[T any](response http.ResponseWriter, data T, message string) {
	res := NewResponse[T](&data, "", message)
	payload, _ := JSON(&res)
	io.Copy(response, payload)
	return
}

func HandleError(response http.ResponseWriter, err error) {
	res := NewResponse[interface{}](nil, err.Error(), "")
	payload, _ := JSON(&res)
	io.Copy(response, payload)
	return
}
