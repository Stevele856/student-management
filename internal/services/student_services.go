// Service → required validation + business rule
package services

import (
	"errors"
	"strings"
	"time"

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
	ErrStudentInfo   = errors.New("invalid student data")
	ErrIDRequired    = errors.New("student ID is required")
	ErrNameFormat    = errors.New("invalid student name format")
	ErrEmailFormat   = errors.New("invalid student email format")
	ErrClassFormat   = errors.New("invalid student class format")
	ErrSubjectFormat = errors.New("invalid subject format")
	ErrValidDOB      = errors.New("date of birth cannot be in the future")

	ErrEmailRequired = errors.New("student name is required")
	ErrNameRequired  = errors.New("student email is required")
	ErrEmailExisted  = errors.New("student email already existed")
	ErrStudentClass  = errors.New("student must belong to a class")

	ErrScore                 = errors.New("score must between 0-10")
	ErrMaxScore              = errors.New("Maximum 10 scores")
	ErrDublicatedSubject     = errors.New("student cannot have duplicate subject score")
	ErrStudentNotFound       = errors.New("student not found")
	ErrSubjectAlreadyExisted = errors.New("subject already existed")
)

func (s *StudentService) AddStudent(student *models.Student) error {
	if student == nil {
		return ErrStudentInfo
	}

	if student.ID == "" {
		student.ID = uuid.New().String()
	}

	student.FullName = strings.TrimSpace(student.FullName)
	student.Email = strings.ToLower(strings.TrimSpace(student.Email))
	if student.FullName == "" || student.Email == "" {
		return ErrStudentInfo
	}

	existed, err := s.repo.GetStudentByEmail(student.Email)
	if err != nil {
		return err
	}
	if existed != nil {
		return ErrEmailExisted
	}

	if !util.IsValidStudentName(student.FullName) {
		return ErrNameFormat
	}

	if !util.IsValidStudentEmail(student.Email) {
		return ErrEmailFormat
	}

	// Validate DOB (not in the future)
	if student.DateOfBirth.After(time.Now()) {
		return ErrValidDOB
	}

	student.Class = strings.TrimSpace(student.Class)
	if student.Class == "" {
		return ErrStudentClass
	}
	if !util.IsValidClass(student.Class) {
		return ErrClassFormat
	}

	if !util.IsValidScores(student.Scores) {
		return ErrScore
	}

	return s.repo.AddStudent(student)
}

func (s *StudentService) AddSubjectScore(studentID string, score *models.SubjectScore) error {
	student, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return err
	}

	if student == nil {
		return ErrStudentNotFound
	}

	score.Subject = strings.TrimSpace(score.Subject)
	if !util.IsValidSubject(score.Subject) {
		return ErrSubjectFormat
	}

	if !util.IsValidSubjectScore(score.Score) {
		return ErrScore
	}

	if len(student.Scores) > 10 {
		return ErrMaxScore
	}

	// Dublicate subject
	for _, s := range student.Scores {
		if strings.EqualFold(s.Subject, score.Subject) {
			return ErrSubjectAlreadyExisted
		}
	}

	return s.repo.AddScores(studentID, score)
}

/*
student == nil => student not found
student != nil =>  student found
*/
