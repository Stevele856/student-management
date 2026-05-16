package studentRepo

import "errors"

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrInvalidPage = errors.New("invalid page")
	ErrInvalidPageSize = errors.New("invalid page size")
)