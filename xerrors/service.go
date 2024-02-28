package xerrors

import "errors"

var (
	ErrNotFound = errors.New("service: not found")
	ErrInternal = errors.New("service: internal error")
)
