package teacherRepo

import "errors"

var (
	ErrTeacherNotFound = errors.New("teacher not found")
)

/*
sentinel error là lỗi được khai báo sẵn thành biến package-level để mọi nơi có thể so sánh bằng errors.Is(...).
*/