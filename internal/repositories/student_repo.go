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
	GetStudentByEmail(email string) (*models.Student, error)

	// CRUD scores
	AddScores(studentID string, Scores *models.SubjectScore) error
	UpdateScores(studentID string, Scores *models.SubjectScore) error
	DeleteScores(studentID, subject string) error
	GetScoresByStudentID(studentID string) ([]*models.SubjectScore, error)
	GetScoresBySubject(studentID, subject string) (*models.SubjectScore, error)

	// Search/filter
	SearchStudentByName(studentName string) ([]*models.Student, error)
	GetStudentByClass(class string) ([]*models.Student, error)
}
