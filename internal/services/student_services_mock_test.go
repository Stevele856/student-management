package services_test

import (
	"github.com/student-management/internal/models"
	"github.com/student-management/internal/predicate"
)

// Mock Student Repository, implement entire interface StudentRepository
// each field is each function, which test case need behavior => assign on it

type MockStudentRepository struct {
	AddStudentFn        func(student *models.Student) error
	UpdateStudentFn     func(student *models.Student) error
	DeleteStudentFn     func(studentID string) error
	GetAllStudentsFn    func() ([]*models.Student, error)
	GetStudentByIDFn    func(studentID string) (*models.Student, error)
	GetStudentByEmailFn func(studentEmail string) (*models.Student, error)

	AddScoreFn            func(studentID string, score *models.SubjectScore) error
	UpdateScoreFn         func(studentID string, score *models.SubjectScore) error
	DeleteScoreFn         func(studentID, subject string) error
	GetScoresByStudentIDFn func(studentID string) ([]*models.SubjectScore, error)
	GetScoresBySubjectFn   func(studentID, subject string) (*models.SubjectScore, error)

	FilterStudentsFn func(p predicate.PredicateStudent) ([]*models.Student, error)
}

/*-------IMPLEMENT INTERFACE ------- */
func (m *MockStudentRepository) AddStudent(student *models.Student) error{
	return m.AddStudentFn(student)
}

func (m *MockStudentRepository) UpdateStudent(student *models.Student) error{
	return m.UpdateStudentFn(student)
}

func (m *MockStudentRepository) DeleteStudent(studentID string) error{
	return m.DeleteStudentFn(studentID)
}

func (m *MockStudentRepository) GetAllStudents() ([]*models.Student, error) {
	return m.GetAllStudentsFn()
}

func (m *MockStudentRepository) GetStudentByID(studentID string) (*models.Student, error){
	return m.GetStudentByIDFn(studentID)
}

func (m *MockStudentRepository) GetStudentByEmail(studentEmail string) (*models.Student, error){
	return m.GetStudentByEmailFn(studentEmail)
}

func (m *MockStudentRepository) AddScore(studentID string, score *models.SubjectScore) error{
	return m.AddScoreFn(studentID, score)
}

func (m *MockStudentRepository) UpdateScore(studentID string, score *models.SubjectScore) error{
	return m.UpdateScoreFn(studentID, score)
}

func (m *MockStudentRepository) DeleteScore(studentID, subject string) error {
	return m.DeleteScoreFn(studentID,subject)
}

func (m *MockStudentRepository) GetScoresByStudentID(studentID string) ([]*models.SubjectScore, error){
	return m.GetScoresByStudentIDFn(studentID)
}

func (m *MockStudentRepository) GetScoresBySubject(studentID, subject string) (*models.SubjectScore, error){
	return m.GetScoresBySubjectFn(studentID, subject)
}

func (m *MockStudentRepository) FilterStudents(p predicate.PredicateStudent) ([]*models.Student, error){
	return m.FilterStudentsFn(p)
}