package teacherRepo

import (
	"github.com/student-management/internal/models/teacher"
	"github.com/student-management/internal/predicate"
)

type TeacherRepository interface {
	AddTeacher(teacher *teacherModels.Teacher) error
	UpdateTeacher(teacher *teacherModels.Teacher) error
	DeleteTeacher(teacherID string) error
	GetAllTeacher() ([]*teacherModels.Teacher, error)
	GetTeacherByID(teacherID string) (*teacherModels.Teacher, error)
	GetTeacherByEmail(teacherEmail string) (*teacherModels.Teacher, error)

	// Advanced
	GetTeacherAssignedBySubject(subject string) ([]*teacherModels.Teacher, error)
	GetTeacherByAssignedClass(classAssigned string) ([]*teacherModels.Teacher, error)
	GetTeacherByStatus(status string) ([]*teacherModels.Teacher, error)

	FilterTeachers(p predicate.PredicateTeacher) ([]*teacherModels.Teacher, error)
	// GetTeachersWithPagination(page, pageSize int) ([]*teacherModels.Teacher, int, error)
}

