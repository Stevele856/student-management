// validate, business logic
package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/student-management/internal/models"
	"github.com/student-management/internal/repositories"
	"github.com/student-management/pkg/util"
)

type StudentService struct {
	repo repositories.StudentRepository
}

// Constructor
func NewStudentService(repo repositories.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

// throw Err
var (
	ErrStudentID    = errors.New("Student ID is required")
	ErrStudentInfo  = errors.New("Invalid student data")
	ErrStudentName  = errors.New("Invalid student name format")
	ErrStudentEmail = errors.New("Invalid student email format")
)

func (s *StudentService) AddStudent(student *models.Student) error {
	if student.ID == "" {
		student.ID = uuid.New().String()
	}

	if student.FullName == "" || student.Email == "" {
		return ErrStudentInfo
	}

	if !util.IsValidStudentName(student.FullName){
		return ErrStudentName
	}

	if util.IsValidStudentEmail(student.Email){
		return ErrStudentEmail
	}

	return nil

}
