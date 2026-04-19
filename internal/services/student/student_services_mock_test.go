package student

import (
	studentModels "github.com/student-management/internal/models/student"
	"github.com/student-management/internal/predicate"
)

// Mock Student Repository, implement entire interface StudentRepository
// each field is each function, which test case need behavior => assign on it

type MockStudentRepository struct {
	AddStudentFn        func(student *studentModels.Student) error
	UpdateStudentFn     func(student *studentModels.Student) error
	DeleteStudentFn     func(studentID string) error
	GetAllStudentsFn    func() ([]*studentModels.Student, error)
	GetStudentByIDFn    func(studentID string) (*studentModels.Student, error)
	GetStudentByEmailFn func(studentEmail string) (*studentModels.Student, error)

	AddScoreFn            func(studentID string, score *studentModels.SubjectScore) error
	UpdateScoreFn         func(studentID string, score *studentModels.SubjectScore) error
	DeleteScoreFn         func(studentID, subject string) error
	GetScoresByStudentIDFn func(studentID string) ([]*studentModels.SubjectScore, error)
	GetScoresBySubjectFn   func(studentID, subject string) (*studentModels.SubjectScore, error)

	FilterStudentsFn func(p predicate.PredicateStudent) ([]*studentModels.Student, error)
	// CSV
	BulkAddStudentsFn func(students []*studentModels.Student) error
}

/*-------IMPLEMENT INTERFACE ------- */
func (m *MockStudentRepository) AddStudent(student *studentModels.Student) error{
	return m.AddStudentFn(student)
}

func (m *MockStudentRepository) UpdateStudent(student *studentModels.Student) error{
	return m.UpdateStudentFn(student)
}

func (m *MockStudentRepository) DeleteStudent(studentID string) error{
	return m.DeleteStudentFn(studentID)
}

func (m *MockStudentRepository) GetAllStudents() ([]*studentModels.Student, error) {
	return m.GetAllStudentsFn()
}

func (m *MockStudentRepository) GetStudentByID(studentID string) (*studentModels.Student, error){
	return m.GetStudentByIDFn(studentID)
}

func (m *MockStudentRepository) GetStudentByEmail(studentEmail string) (*studentModels.Student, error){
	return m.GetStudentByEmailFn(studentEmail)
}

func (m *MockStudentRepository) AddScore(studentID string, score *studentModels.SubjectScore) error{
	return m.AddScoreFn(studentID, score)
}

func (m *MockStudentRepository) UpdateScore(studentID string, score *studentModels.SubjectScore) error{
	return m.UpdateScoreFn(studentID, score)
}

func (m *MockStudentRepository) DeleteScore(studentID, subject string) error {
	return m.DeleteScoreFn(studentID,subject)
}

func (m *MockStudentRepository) GetScoresByStudentID(studentID string) ([]*studentModels.SubjectScore, error){
	return m.GetScoresByStudentIDFn(studentID)
}

func (m *MockStudentRepository) GetScoresBySubject(studentID, subject string) (*studentModels.SubjectScore, error){
	return m.GetScoresBySubjectFn(studentID, subject)
}

func (m *MockStudentRepository) FilterStudents(p predicate.PredicateStudent) ([]*studentModels.Student, error){
	return m.FilterStudentsFn(p)
}

func (m *MockStudentRepository) BulkAddStudents(student []*studentModels.Student) error{
	return m.BulkAddStudentsFn(student)
}