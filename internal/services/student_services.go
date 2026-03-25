// Service → required validation + business rule
package services

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/student-management/internal/models"
	"github.com/student-management/internal/predicate"
	"github.com/student-management/internal/repositories"
	"github.com/student-management/pkg/utils"
)

type StudentService struct {
	repo repositories.StudentRepository
}

// CONSTRUCTOR
func NewStudentService(repo repositories.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

// THROW ERROR
var (
	ErrStudentInfo   = errors.New("invalid student data")
	ErrIDRequired    = errors.New("student ID is required")
	ErrNameFormat    = errors.New("invalid student name format")
	ErrEmailFormat   = errors.New("invalid student email format")
	ErrClassFormat   = errors.New("invalid student class format")
	ErrSubjectFormat = errors.New("invalid subject format")
	ErrValidDOB      = errors.New("date of birth cannot be in the future")

	ErrEmailRequired = errors.New("student email is required")
	ErrNameRequired  = errors.New("student name is required")
	ErrEmailExisted  = errors.New("student email already existed")
	ErrStudentClass  = errors.New("student must belong to a class")
	ErrSubjectEmpty  = errors.New("subject is require to check score")

	ErrScore                 = errors.New("score must between 0-10")
	ErrMaxScore              = errors.New("Maximum 10 scores")
	ErrDublicatedSubject     = errors.New("student cannot have duplicate subject score")
	ErrStudentNotFound       = errors.New("student not found")
	ErrSubjectAlreadyExisted = errors.New("subject already existed")
	ErrStudentID             = errors.New("student ID not found")
	ErrStudentEmail          = errors.New("student email does not exist")

	ErrInvalidYear     = errors.New("invalid year of birth")
	ErrInvalidGender   = errors.New("gender must be 'male' or 'female'")
	ErrAddressTooShort = errors.New("address too short")
	ErrInvalidMinMax   = errors.New("min score cannot be greater than max score")
	ErrStudentRank     = errors.New("student rank must be 'excellent', 'good', 'average' or 'weak'")
)

// ADD STUDENT
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

	if !utils.IsValidStudentName(student.FullName) {
		return ErrNameFormat
	}

	if !utils.IsValidStudentEmail(student.Email) {
		return ErrEmailFormat
	}

	// VALIDATE DOB NOT IN THE FUTURE
	if student.DateOfBirth.After(time.Now()) {
		return ErrValidDOB
	}

	student.Class = strings.TrimSpace(student.Class)
	if student.Class == "" {
		return ErrStudentClass
	}
	if !utils.IsValidClass(student.Class) {
		return ErrClassFormat
	}

	if !utils.IsValidScores(student.Scores) {
		return ErrScore
	}

	return s.repo.AddStudent(student)
}

// UPDATE STUDENT
func (s *StudentService) UpdateStudent(student *models.Student) error {
	if !utils.IsValidStudentName(student.FullName) {
		return ErrNameFormat
	}

	if !utils.IsValidStudentEmail(student.Email) {
		return ErrEmailFormat
	}

	if student.DateOfBirth.After(time.Now()) {
		return ErrValidDOB
	}

	if !utils.IsValidClass(student.Class) {
		return ErrClassFormat
	}

	if !utils.IsValidScores(student.Scores) {
		return ErrScore
	}

	return s.repo.UpdateStudent(student)
}

// DELETE STUDENT
func (s *StudentService) DeleteStudent(studentID string) error {
	studentID = strings.TrimSpace(studentID)

	if studentID == "" {
		return ErrStudentID
	}

	student, err := s.repo.GetStudentByID(studentID)

	if err != nil {
		return err
	}

	if student == nil {
		return ErrStudentNotFound
	}

	return s.repo.DeleteStudent(studentID)
}

// GET ALL STUDENT
func (s *StudentService) GetAllStudents() ([]*models.Student, error) {
	return s.repo.GetAllStudents()
}

// GET STUDENT BY ID
func (s *StudentService) GetStudentByID(studentID string) (*models.Student, error) {
	studentID = strings.TrimSpace(studentID)
	if studentID == "" {
		return nil, ErrStudentID 
	}
	return s.repo.GetStudentByID(studentID)
}

// GET STUDENT BY EMAIL
func (s *StudentService) GetStudentByEmail(studentEmail string) (*models.Student, error) {
	studentEmail = strings.TrimSpace(studentEmail)
	if studentEmail == "" {
		return nil, ErrStudentEmail
	}
	if !utils.IsValidStudentEmail(studentEmail) {
		return nil, ErrEmailFormat
	}
	return s.repo.GetStudentByEmail(studentEmail)
}

/* ------------------------ */

// ADD SUBJECT SCORE
func (s *StudentService) AddScore(studentID string, score *models.SubjectScore) error {

	student, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return err
	}

	score.Subject = strings.TrimSpace(score.Subject)
	if !utils.IsValidSubject(score.Subject) {
		return ErrSubjectFormat
	}

	if !utils.IsValidSubjectScore(score.Score) {
		return ErrScore
	}

	if len(student.Scores) >= 10 {
		return ErrMaxScore
	}

	// CHECK DIBLICATE SUBJECT
	for _, s := range student.Scores {
		if strings.EqualFold(s.Subject, score.Subject) {
			return ErrSubjectAlreadyExisted
		}
	}

	return s.repo.AddScore(studentID, score)
}

// UPDATE SCORE
func (s *StudentService) UpdateScore(studentID string, score *models.SubjectScore) error {
	studentID = strings.TrimSpace(studentID)
	_, err := s.repo.GetStudentByID(studentID)

	if err != nil {
		return err
	}

	score.Subject = strings.TrimSpace(score.Subject)
	if !utils.IsValidSubject(score.Subject) {
		return ErrSubjectFormat
	}

	if !utils.IsValidSubjectScore(score.Score) {
		return ErrScore
	}

	return s.repo.UpdateScore(studentID, score)
}

// DELETE SCORE
func (s *StudentService) DeleteScore(studentID, subject string) error {
	studentID = strings.TrimSpace(studentID)
	subject = strings.TrimSpace(subject)

	if studentID == "" {
		return ErrIDRequired
	}

	if subject == "" {
		return ErrSubjectEmpty
	}

	_, err := s.repo.GetStudentByID(studentID)

	if err != nil {
		return err
	}
	if !utils.IsValidSubject(subject) {
		return ErrSubjectFormat
	}

	return s.repo.DeleteScore(studentID, subject)
}

// GET SCORE BY STUDENT ID
func (s *StudentService) GetScoresByStudentID(studentID string) ([]*models.SubjectScore, error) {
	studentID = strings.TrimSpace(studentID)

	if studentID == "" {
		return nil, ErrIDRequired
	}

	return s.repo.GetScoresByStudentID(studentID)
}

// GET SCORE BY SUBJECT
func (s *StudentService) GetScoresBySubject(studentID, subject string) (*models.SubjectScore, error) {
	studentID = strings.TrimSpace(studentID)
	subject = strings.TrimSpace(subject)

	if studentID == "" {
		return nil, ErrIDRequired
	}

	if subject == "" {
		return nil, ErrSubjectEmpty
	}

	return s.repo.GetScoresBySubject(studentID, subject)

}

func (s *StudentService) FilterStudents(filter *models.FilterStudents) ([]*models.Student, error) {
	if filter == nil {
		return s.repo.GetAllStudents()
	}

	// VALIDATE NAME
	if filter.Name != "" {
		filter.Name = strings.TrimSpace(filter.Name)
		if !utils.IsValidStudentName(filter.Name) {
			return nil, ErrNameFormat
		}
	}

	// VALIDATE CLASS
	if filter.Class != "" {
		filter.Class = strings.TrimSpace(filter.Class)
		if !utils.IsValidClass(filter.Class) {
			return nil, ErrClassFormat
		}
	}

	// VALIDATE DATE OF BIRTH
	if filter.YearOfBirth != 0 {
		currentYear := time.Now().Year()

		if filter.YearOfBirth > currentYear {
			return nil, ErrInvalidYear
		}
	}

	// VALIDATE GENDER
	if filter.Gender != "" {
		filter.Gender = strings.TrimSpace(filter.Gender)

		if filter.Gender != "male" && filter.Gender != "female" {
			return nil, ErrInvalidGender
		}
	}

	// VALIDATE ADDRESS
	if filter.Address != "" && len([]rune(filter.Address)) < 5 {
		return nil, ErrAddressTooShort
	}

	// VALIDATE SCORE RANGE
	if filter.MinAvgScore < 0 || filter.MinAvgScore > 10 ||
		filter.MaxAvgScore < 0 || filter.MaxAvgScore > 10 {
		return nil, ErrScore
	}

	if filter.MinAvgScore > 0 && filter.MaxAvgScore > 0 &&
		filter.MinAvgScore > filter.MaxAvgScore {
		return nil, ErrInvalidMinMax
	}

	// VALIDATE STUDENT RANK
	if filter.StudentRank != "" {
		rank := strings.TrimSpace(string(filter.StudentRank))
		if rank != "excellent" && rank != "good" && rank != "average" && rank != "weak" {
			return nil, ErrStudentRank
		}
	}

	// PREDICATE STUDENT
	predicates := []predicate.PredicateStudent{}
	
	if filter.Name != ""{
		predicates = append(predicates, predicate.ByName(filter.Name))
	}

	if filter.Class != ""{
		predicates = append(predicates, predicate.ByClass(filter.Class))
	}

	if filter.Gender != ""{
		predicates = append(predicates, predicate.ByGender(filter.Gender))
	}

	if filter.Address != ""{
		predicates = append(predicates, predicate.ByAddress(filter.Address))
	}

	if filter.MinAvgScore > 0 && filter.MaxAvgScore > 0{
		predicates = append(predicates, predicate.ByAvgScore(filter.MinAvgScore, filter.MaxAvgScore))
	}

	if filter.StudentRank != ""{
		predicates = append(predicates, predicate.ByRank(filter.StudentRank))
	}

	return s.repo.FilterStudents(predicate.And(predicates...))

}

/*
student == nil => student not found
student != nil =>  student found
*/

/* Update student
- Full Update - PUT
*/
