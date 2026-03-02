// Step 2: interface

package repositories

import "github.com/student-management/internal/models"

type StudentRepository interface {

	// CRUD student
	AddStudent(student *models.Student) error
	UpdateStudent(student *models.Student) error
	DeleteStudent(id string) error
	GetAllStudents() ([]*models.Student, error)
	GetStudentByID(id string) (*models.Student, error)

	// CRUD grades
	AddGrade(studentID string, grade *models.Grade) error
	UpdateGrade(studentID string, grade *models.Grade) error
	DeleteGrade(studentID, subject string) error
	GetGradeByStudentID(studentID string) ([]*models.Grade, error)
	GetGradeBySubject(studentID, subject string) (*models.Grade, error)

	// Search/filter
	SearchStudentByName(studentName string) ([]*models.Student, error)
	GetStudentByClass(class string) ([]*models.Student, error)
}
