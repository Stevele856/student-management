package teacherRepo

import "github.com/student-management/internal/models/teacher"

type TeacherRepository interface {
	AddTeacher(teacher *teacherModels.Teacher) error
	UpdateTeacher(teacher *teacherModels.Teacher) error
	DeleteTeacher(teacherID string) error
	GetAllTeacher() ([]*teacherModels.Teacher, error)
	GetTeacherByID(teacherID string) (*teacherModels.Teacher, error)
	GetTeacherByEmail(teacherEmail string) (*teacherModels.Teacher, error)

	// Advanced
	GetTeacherAssignedBySubject(subject string) ([]*teacherModels.Teacher, error)
	GetTeacherByClass(class string) ([]*teacherModels.Teacher, error)
	GetTeacherStatus(status string) ([]*teacherModels.Teacher, error)

	// FilterTeachers(p predicate.PredicateTeacher) ([]*teacherModels.Teacher, error)
	// GetTeachersWithPagination(page, pageSize int) ([]*teacherModels.Teacher, int, error)
}

