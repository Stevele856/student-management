// Service → required validation + business rule
package student

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	studentModels "github.com/student-management/internal/models/student"
	predicates "github.com/student-management/internal/predicates"

	studentRepo "github.com/student-management/internal/repositories/student"
	"github.com/student-management/pkg/utils"
)

type StudentService struct {
	repo studentRepo.StudentRepository
}

// CONSTRUCTOR
func NewStudentService(repo studentRepo.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

// THROW ERROR
var (
	ErrStudentData       = errors.New("invalid student data")
	ErrStudentIDRequired = errors.New("student ID is required")
	ErrNameFormat        = errors.New("invalid student name format")
	ErrEmailFormat       = errors.New("invalid student email format")
	ErrClassFormat       = errors.New("invalid student class format")
	ErrSubjectFormat     = errors.New("invalid subject format")
	ErrValidDOB          = errors.New("date of birth cannot be in the future")

	ErrEmailRequired = errors.New("student email is required")
	ErrNameRequired  = errors.New("student name is required")
	ErrStudentEmailExisted  = errors.New("student email already existed")
	ErrStudentClass  = errors.New("student must belong to a class")
	ErrSubjectEmpty  = errors.New("subject is require to check score")

	ErrScoreData             = errors.New("invalid score data")
	ErrScoreRange            = errors.New("score must between 0-10")
	ErrMaxScore              = errors.New("Maximum 10 scores")
	ErrDublicatedSubject     = errors.New("student cannot have duplicate subject score")
	ErrStudentNotFound       = errors.New("student not found")
	ErrSubjectAlreadyExisted = errors.New("subject already existed")
	ErrStudentIDNotFound     = errors.New("student ID not found")
	ErrStudentEmail          = errors.New("student email does not exist")

	ErrInvalidYear     = errors.New("invalid year of birth")
	ErrInvalidGender   = errors.New("gender must be 'male' or 'female'")
	ErrAddressTooShort = errors.New("address too short")
	ErrInvalidMinMax   = errors.New("min score cannot be greater than max score")
	ErrStudentRank     = errors.New("student rank must be 'excellent', 'good', 'average' or 'weak'")
	ErrSubjectNotFound = errors.New("subject not found")
)

// ADD STUDENT
func (s *StudentService) AddStudent(student *studentModels.Student) error {

	if student == nil {
		return ErrStudentData
	}

	student = normalizeStudent(student)

	if err := validateStudent(student); err != nil {
		return err
	}

	existed, err := s.repo.GetStudentByEmail(student.Email)
	if err == nil && existed != nil {
		return ErrStudentEmailExisted
	}

	// GEN UUID WHEN MAKE SURE THAT STUDENT WILL BE ADDED IN DB
	if student.ID == "" {
		student.ID = uuid.New().String()
	}

	return s.repo.AddStudent(student)
}

// UPDATE STUDENT
func (s *StudentService) UpdateStudent(student *studentModels.Student) error {

	if student == nil {
		return ErrStudentData
	}

	studentID, err := s.requireStudentID(student.ID)
	if err != nil {
		return err
	}

	existingStudent, err := s.ensureStudentExistsByID(studentID)
	if err != nil {
		return err
	}

	student = normalizeStudent(student)

	// Preserve system fields from old record
	student.ID = existingStudent.ID

	if err := validateStudent(student); err != nil {
		return err
	}

	if err := s.ensureStudentUnique(student.Email, student.ID); err != nil {
		return err
	}

	student.UpdatedAt = time.Now().UTC()

	return s.repo.UpdateStudent(student)
}

// DELETE STUDENT
func (s *StudentService) DeleteStudent(studentID string) error {
	studentID, err := s.requireStudentID(studentID)
	if err != nil {
		return err
	}

	existing, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return err
	}

	if existing == nil {
		return ErrStudentNotFound
	}

	return s.repo.DeleteStudent(studentID)
}

// GET ALL STUDENT
func (s *StudentService) GetAllStudents() ([]*studentModels.Student, error) {
	return s.repo.GetAllStudents()
}

// GET STUDENT BY ID
func (s *StudentService) GetStudentByID(studentID string) (*studentModels.Student, error) {
	studentID, err := s.requireStudentID(studentID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetStudentByID(studentID)
}

// GET STUDENT BY EMAIL
func (s *StudentService) GetStudentByEmail(studentEmail string) (*studentModels.Student, error) {
	studentEmail = strings.TrimSpace(studentEmail)
	if studentEmail == "" {
		return nil, ErrEmailRequired
	}
	if !utils.IsValidEmail(studentEmail) {
		return nil, ErrEmailFormat
	}
	return s.repo.GetStudentByEmail(studentEmail)
}

// FILTER STUDENTS
func (s *StudentService) FilterStudents(filter *studentModels.FilterStudents) ([]*studentModels.Student, error) {
	if filter == nil {
		return s.repo.GetAllStudents()
	}

	filter = normalizeFilterStudent(filter)

	if err := validateFilterStudents(filter); err != nil {
		return nil, err
	}

	// PREDICATE STUDENT
	predicate := filterPredicateStudents(filter)

	return s.repo.FilterStudents(predicates.AndStudent(predicate...))

}

/* ------------------------ */

// ADD SUBJECT SCORE
func (s *StudentService) AddScore(studentID string, score *studentModels.SubjectScore) error {

	if score == nil {
		return ErrScoreData
	}
	studentID, score = normalizeScore(studentID, score)

	if err := validateScore(studentID, score); err != nil {
		return err
	}

	student, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return err
	}

	if student == nil {
		return ErrStudentNotFound
	}

	if err := validateAddScoreRules(student, score); err != nil {
		return err
	}

	return s.repo.AddScore(studentID, score)
}

// UPDATE SCORE
func (s *StudentService) UpdateScore(studentID string, score *studentModels.SubjectScore) error {
	if score == nil {
		return ErrScoreData
	}

	studentID, score = normalizeScore(studentID, score)

	if err := validateScore(studentID, score); err != nil {
		return err
	}

	student, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return err
	}

	if student == nil {
		return ErrStudentNotFound
	}

	if err := validateSubjectExists(student, score.Subject); err != nil {
		return err
	}

	return s.repo.UpdateScore(studentID, score)
}

// DELETE SCORE
func (s *StudentService) DeleteScore(studentID, subject string) error {
	studentID, subject = normalizeSubject(studentID, subject)

	if err := validateSubject(studentID, subject); err != nil {
		return err
	}

	student, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return err
	}

	if student == nil {
		return ErrStudentNotFound
	}

	if err := validateSubjectExists(student, subject); err != nil {
		return err
	}

	return s.repo.DeleteScore(studentID, subject)
}

// GET SCORE BY STUDENT ID
func (s *StudentService) GetScoresByStudentID(studentID string) ([]*studentModels.SubjectScore, error) {
	studentID = strings.TrimSpace(studentID)

	if studentID == "" {
		return nil, ErrStudentIDRequired
	}

	existing, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, ErrStudentNotFound
	}

	return s.repo.GetScoresByStudentID(studentID)
}

// GET SCORE BY SUBJECT
func (s *StudentService) GetScoresBySubject(studentID, subject string) (*studentModels.SubjectScore, error) {
	studentID, subject = normalizeSubject(studentID, subject)

	if err := validateSubject(studentID, subject); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, ErrStudentNotFound
	}

	return s.repo.GetScoresBySubject(studentID, subject)

}

func normalizeScore(studentID string, score *studentModels.SubjectScore) (string, *studentModels.SubjectScore) {
	cp := *score
	cp.Subject = strings.TrimSpace(cp.Subject)
	return strings.TrimSpace(studentID), &cp
}

func validateScore(studentID string, score *studentModels.SubjectScore) error {
	if studentID == "" {
		return ErrStudentIDRequired
	}
	if score == nil {
		return ErrStudentData
	}
	if !utils.IsValidSubject(score.Subject) {
		return ErrSubjectFormat
	}

	if !utils.IsValidSubjectScore(score.Score) {
		return ErrScoreData
	}
	return nil
}

func validateAddScoreRules(student *studentModels.Student, score *studentModels.SubjectScore) error {
	if len(student.Scores) >= 10 {
		return ErrMaxScore
	}

	// CHECK DIBLICATE SUBJECT
	for _, existing := range student.Scores {
		if strings.EqualFold(existing.Subject, score.Subject) {
			return ErrSubjectAlreadyExisted
		}
	}
	return nil
}

func validateSubjectExists(student *studentModels.Student, subject string) error {

	// CHECK DIBLICATE SUBJECT
	for _, existing := range student.Scores {
		if strings.EqualFold(existing.Subject, subject) {
			return nil
		}
	}
	return ErrSubjectNotFound
}

func normalizeSubject(studentID, subject string) (string, string) {
	return strings.TrimSpace(studentID), strings.TrimSpace(subject)
}

func validateSubject(studentID, subject string) error {
	if studentID == "" {
		return ErrStudentIDRequired
	}

	if subject == "" {
		return ErrSubjectEmpty
	}

	if !utils.IsValidSubject(subject) {
		return ErrSubjectFormat
	}

	return nil
}

func validateFilterScoreRange(filter *studentModels.FilterStudents) error {
	// VALIDATE SCORE RANGE
	if filter.MinAvgScore < 0 || filter.MinAvgScore > 10 ||
		filter.MaxAvgScore < 0 || filter.MaxAvgScore > 10 {
		return ErrScoreRange
	}

	if filter.MinAvgScore > filter.MaxAvgScore &&
		filter.MaxAvgScore != 0 {
		return ErrInvalidMinMax
	}

	return nil
}

func filterPredicateStudents(filter *studentModels.FilterStudents) []predicates.PredicateStudent {
	predicate := []predicates.PredicateStudent{}

	if filter.Name != "" {
		predicate = append(predicate, predicates.ByName(filter.Name))
	}

	if filter.Class != "" {
		predicate = append(predicate, predicates.ByClass(filter.Class))
	}

	if filter.Gender != "" {
		predicate = append(predicate, predicates.ByGender(string(filter.Gender)))
	}

	if filter.Address != "" {
		predicate = append(predicate, predicates.ByAddress(filter.Address))
	}

	if filter.MinAvgScore > 0 && filter.MaxAvgScore > 0 {
		predicate = append(predicate, predicates.ByAvgScore(filter.MinAvgScore, filter.MaxAvgScore))
	}

	if filter.StudentRank != "" {
		predicate = append(predicate, predicates.ByRank(filter.StudentRank))
	}
	return predicate
}

/*
student == nil => student not found
student != nil =>  student found
*/

/* Update student
- Full Update - PUT
*/
