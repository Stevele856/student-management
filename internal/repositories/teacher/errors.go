package teacherRepo

import "errors"

var (
	ErrTeacherNotFound = errors.New("teacher not found")
	ErrTeacherAlreadyExists = errors.New("teacher already exists")
	ErrInvalidPage = errors.New("invalid page")
	ErrInvalidPageSize = errors.New("invalid page size")
)

/*
sentinel error là lỗi được khai báo sẵn thành biến package-level để mọi nơi có thể so sánh bằng errors.Is(...).
*/