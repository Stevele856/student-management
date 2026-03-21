// Step 2: interface

package repositories

import (
	"github.com/student-management/internal/models"
	"github.com/student-management/internal/predicate"
)

type StudentRepository interface {

	// CRUD student
	AddStudent(student *models.Student) error
	UpdateStudent(student *models.Student) error
	DeleteStudent(studentID string) error
	GetAllStudents() ([]*models.Student, error)
	GetStudentByID(studentID string) (*models.Student, error)
	GetStudentByEmail(StudentEmail string) (*models.Student, error)

	// CRUD scores
	AddScore(studentID string, score *models.SubjectScore) error
	UpdateScore(studentID string, score *models.SubjectScore) error
	DeleteScore(studentID, subject string) error
	GetScoresByStudentID(studentID string) ([]*models.SubjectScore, error)
	GetScoresBySubject(studentID, subject string) (*models.SubjectScore, error)

	// Filter
	FilterStudents(p predicate.PredicateStudent) ([]*models.Student, error)	
}
