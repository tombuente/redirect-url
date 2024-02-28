package xerrors

import "errors"

var (
	ErrSQLNotFound = errors.New("db: not found")
	ErrSQLInternal = errors.New("db: internal error")
)
