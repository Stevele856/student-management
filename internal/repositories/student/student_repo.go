// Step 2: interface

package studentRepo

import (
	"github.com/student-management/internal/models/student"
	"github.com/student-management/internal/predicate"
)

type StudentRepository interface {

	// CRUD student
	AddStudent(student *studentModels.Student) error
	UpdateStudent(student *studentModels.Student) error
	DeleteStudent(studentID string) error
	GetAllStudents() ([]*studentModels.Student, error)
	GetStudentByID(studentID string) (*studentModels.Student, error)
	GetStudentByEmail(StudentEmail string) (*studentModels.Student, error)

	// CRUD scores
	AddScore(studentID string, score *studentModels.SubjectScore) error
	UpdateScore(studentID string, score *studentModels.SubjectScore) error
	DeleteScore(studentID, subject string) error
	GetScoresByStudentID(studentID string) ([]*studentModels.SubjectScore, error)
	GetScoresBySubject(studentID, subject string) (*studentModels.SubjectScore, error)

	// Filter
	FilterStudents(p predicate.PredicateStudent) ([]*studentModels.Student, error)	

	// CSV
	BulkAddStudents(students []*studentModels.Student) error
}
