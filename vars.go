package gorocks

import "errors"

var (
	// ErrRequestNil returns when http request is nil
	ErrRequestNil = errors.New("http request is nil")
	// ErrRequestBodyNil returns when request body is nil
	ErrRequestBodyNil = errors.New("http request body is nil")
)
