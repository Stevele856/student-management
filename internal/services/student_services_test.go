package services

import (
	"errors"

	"github.com/student-management/internal/models"
)

type MockStudentRepo struct {
	students map[string]*models.Student
	getByEmailFn func(email string) (*models.Student, error)
	addStudentFn func(student *models.Student) error
}

func NewMockStudentRepo() *MockStudentRepo{
	return &MockStudentRepo{
		students: make(map[string]*models.Student),
	}
}

// IMPLEMENT INTERFACES

func (m *MockStudentRepo) GetStudentByEmail(email string) (*models.Student, error){
	if m.getByEmailFn != nil {
		return m.getByEmailFn(email)
	}
	return nil, errors.New("student not found")
}

func (m *MockStudentRepo) AddStudent(student *models.Student) error {
	if m.addStudentFn != nil {
		return m.addStudentFn(student)
	}

	m.students[student.ID] = student
	return nil
}