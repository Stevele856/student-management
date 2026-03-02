// validate, business logic
package services

import (
	"github.com/student-management/internal/models"
	"github.com/student-management/internal/repositories"
)

type StudentService struct {
	repo repositories.StudentRepository
}


// Constructor
func NewStudentService(repo repositories.StudentRepository) *StudentService{
	return &StudentService{repo: repo}
}

func (s *StudentService) AddStudent(student *models.Student) error {
	
}