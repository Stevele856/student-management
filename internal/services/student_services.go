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
	ErrStudentData   = errors.New("invalid student data")
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
	ErrSubjectNotFound = errors.New("subject not found")
)

// ADD STUDENT
func (s *StudentService) AddStudent(student *models.Student) error {

	if student == nil {
		return ErrStudentData
	}

	student = normalizeStudent(student)

	if err := validateStudent(student); err != nil {
		return err
	}

	existed, err := s.repo.GetStudentByEmail(student.Email)
	if err == nil && existed != nil {
		return ErrEmailExisted
	}

	// GEN UUID WHEN MAKE SURE THAT STUDENT WILL BE ADDED IN DB
	if student.ID == "" {
		student.ID = uuid.New().String()
	}

	return s.repo.AddStudent(student)
}

// UPDATE STUDENT
func (s *StudentService) UpdateStudent(student *models.Student) error {

	if student == nil {
		return ErrStudentData
	}

	if student.ID == "" {
		return ErrIDRequired
	}
	existing, err := s.repo.GetStudentByID(student.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return ErrStudentNotFound
	}

	student = normalizeStudent(student)

	if err := validateStudent(student); err != nil {
		return err
	}

	existedEmail, err := s.repo.GetStudentByEmail(student.Email)
	if err != nil {
		return err
	}
	if existedEmail != nil && existedEmail.ID != student.ID {
		return ErrEmailExisted
	}

	return s.repo.UpdateStudent(student)
}

// DELETE STUDENT
func (s *StudentService) DeleteStudent(studentID string) error {
	studentID = strings.TrimSpace(studentID)

	if studentID == "" {
		return ErrStudentID
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
func (s *StudentService) UpdateScore(studentID string, score *models.SubjectScore) error {
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
func (s *StudentService) GetScoresByStudentID(studentID string) ([]*models.SubjectScore, error) {
	studentID = strings.TrimSpace(studentID)

	if studentID == "" {
		return nil, ErrIDRequired
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
func (s *StudentService) GetScoresBySubject(studentID, subject string) (*models.SubjectScore, error) {
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

// FILTER STUDENTS
func (s *StudentService) FilterStudents(filter *models.FilterStudents) ([]*models.Student, error) {
	if filter == nil {
		return s.repo.GetAllStudents()
	}

	filter = normalizeFilterStudent(filter)

	if err := validateFilterStudents(filter); err != nil {
		return nil, err
	}

	// PREDICATE STUDENT
	predicates := filterPredicateStudents(filter)

	return s.repo.FilterStudents(predicate.And(predicates...))

}

func normalizeStudent(s *models.Student) *models.Student {
	cp := *s
	cp.FullName = strings.TrimSpace(cp.FullName)
	cp.Email = strings.ToLower(strings.TrimSpace(cp.Email))
	cp.Class = strings.TrimSpace(cp.Class)
	return &cp
}

func validateStudent(student *models.Student) error {
	if student.FullName == "" {
		return ErrNameRequired
	}

	if student.Email == "" {
		return ErrEmailRequired
	}

	if !utils.IsValidStudentName(student.FullName) {
		return ErrNameFormat
	}

	if !utils.IsValidStudentEmail(student.Email) {
		return ErrEmailFormat
	}

	if student.DateOfBirth.After(time.Now()) {
		return ErrValidDOB
	}

	if student.Class == "" {
		return ErrStudentClass
	}

	if !utils.IsValidClass(student.Class) {
		return ErrClassFormat
	}

	if !utils.IsValidScores(student.Scores) {
		return ErrScore
	}
	return nil
}

func normalizeScore(studentID string, score *models.SubjectScore) (string, *models.SubjectScore) {
	cp := *score
	cp.Subject = strings.TrimSpace(cp.Subject)
	return strings.TrimSpace(studentID), &cp
}

func validateScore(studentID string, score *models.SubjectScore) error {
	if studentID == "" {
		return ErrIDRequired
	}
	if score == nil {
		return ErrStudentData
	}
	if !utils.IsValidSubject(score.Subject) {
		return ErrSubjectFormat
	}

	if !utils.IsValidSubjectScore(score.Score) {
		return ErrScore
	}
	return nil
}

func validateAddScoreRules(student *models.Student, score *models.SubjectScore) error {
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

func validateSubjectExists(student *models.Student, subject string) error {

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
		return ErrIDRequired
	}

	if subject == "" {
		return ErrSubjectEmpty
	}

	if !utils.IsValidSubject(subject) {
		return ErrSubjectFormat
	}

	return nil
}

func normalizeFilterStudent(filter *models.FilterStudents) *models.FilterStudents {
	cp := *filter
	cp.Name = strings.TrimSpace(cp.Name)
	cp.Class = strings.TrimSpace(cp.Class)
	cp.Gender = strings.TrimSpace(cp.Gender)

	return &cp
}

func validateFilterStudents(filter *models.FilterStudents) error {
	// VALIDATE NAME
	if filter.Name != "" && !utils.IsValidStudentName(filter.Name) {
		return ErrNameFormat
	}

	// VALIDATE CLASS
	if filter.Class != "" && !utils.IsValidClass(filter.Class) {
		return ErrClassFormat
	}

	// VALIDATE DATE OF BIRTH
	if filter.YearOfBirth != 0 && filter.YearOfBirth > time.Now().Year() {
		return ErrInvalidYear
	}

	// VALIDATE GENDER
	if filter.Gender != "" && filter.Gender != "male" && filter.Gender != "female" {
		return ErrInvalidGender
	}

	// VALIDATE ADDRESS
	if filter.Address != "" && len([]rune(filter.Address)) < 5 {
		return ErrAddressTooShort
	}
	// VALIDATE SCORE RANGE
	if err := validateFilterScoreRange(filter); err != nil {
		return err
	}

	// VALIDATE STUDENT RANK
	if err := validateStudentRank(filter); err != nil {
		return err
	}
	return nil
}

func validateFilterScoreRange(filter *models.FilterStudents) error {
	// VALIDATE SCORE RANGE
	if filter.MinAvgScore < 0 || filter.MinAvgScore > 10 ||
		filter.MaxAvgScore < 0 || filter.MaxAvgScore > 10 {
		return ErrScore
	}

	if filter.MinAvgScore > filter.MaxAvgScore &&
		filter.MaxAvgScore != 0 {
		return ErrInvalidMinMax
	}

	return nil
}

func validateStudentRank(filter *models.FilterStudents) error {
	// VALIDATE STUDENT RANK
	if filter.StudentRank != "" {
		rank := strings.TrimSpace(string(filter.StudentRank))
		if rank != "excellent" && rank != "good" && rank != "average" && rank != "weak" {
			return ErrStudentRank
		}
	}
	return nil
}

func filterPredicateStudents(filter *models.FilterStudents) []predicate.PredicateStudent {
	predicates := []predicate.PredicateStudent{}

	if filter.Name != "" {
		predicates = append(predicates, predicate.ByName(filter.Name))
	}

	if filter.Class != "" {
		predicates = append(predicates, predicate.ByClass(filter.Class))
	}

	if filter.Gender != "" {
		predicates = append(predicates, predicate.ByGender(filter.Gender))
	}

	if filter.Address != "" {
		predicates = append(predicates, predicate.ByAddress(filter.Address))
	}

	if filter.MinAvgScore > 0 && filter.MaxAvgScore > 0 {
		predicates = append(predicates, predicate.ByAvgScore(filter.MinAvgScore, filter.MaxAvgScore))
	}

	if filter.StudentRank != "" {
		predicates = append(predicates, predicate.ByRank(filter.StudentRank))
	}
	return predicates
}

/*
student == nil => student not found
student != nil =>  student found
*/

/* Update student
- Full Update - PUT
*/
