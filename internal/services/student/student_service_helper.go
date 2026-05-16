package student

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/student-management/internal/models"
	studentModels "github.com/student-management/internal/models/student"
	"github.com/student-management/pkg/utils"
)

func normalizeStudent(student *studentModels.Student) *studentModels.Student {
	if student == nil {
		return nil
	}
	s := *student
	s.FullName = strings.TrimSpace(s.FullName)
	s.Email = strings.ToLower(strings.TrimSpace(s.Email))
	s.Class = strings.TrimSpace(s.Class)
	s.Address = strings.TrimSpace(s.Address)
	return &s
}

func (s *StudentService) requireStudentID(studentID string) (string, error) {
	studentID = strings.TrimSpace(studentID)
	if studentID == "" {
		return "", ErrStudentIDRequired
	}
	return studentID, nil
}

func (s *StudentService) ensureStudentExistsByID(studentID string) (*studentModels.Student, error) {
	existing, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrStudentNotFound
	}
	return existing, nil
}

func (s *StudentService) ensureStudentUnique(email, studentID string) error {
	// Unique student email
	existedEmail, err := s.repo.GetStudentByEmail(email)

	// DB or network error
	if err != nil && !errors.Is(err, ErrStudentNotFound){
		return err
	}

	if err == nil && existedEmail != nil && existedEmail.ID != studentID {
		return ErrStudentEmailExisted
	}
	return nil
}

func validateStudent(student *studentModels.Student) error {
	if student.FullName == "" {
		return ErrNameRequired
	}

	if student.Email == "" {
		return ErrEmailRequired
	}

	if !utils.IsValidName(student.FullName) {
		return ErrNameFormat
	}

	if !utils.IsValidEmail(student.Email) {
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
		return ErrScoreData
	}
	return nil
}

func normalizeFilterStudent(filter *studentModels.FilterStudents) *studentModels.FilterStudents {
	cp := *filter
	cp.Name = strings.TrimSpace(cp.Name)
	cp.Class = strings.TrimSpace(cp.Class)
	cp.Gender = models.Gender(strings.ToLower(string(cp.Gender)))

	return &cp
}

func validateFilterStudents(filter *studentModels.FilterStudents) error {
	// VALIDATE NAME
	if filter.Name != "" && !utils.IsValidName(filter.Name) {
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
	if filter.Gender != "" && filter.Gender != models.GenderMale && filter.Gender != models.GenderFemale {
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

func validateStudentRank(filter *studentModels.FilterStudents) error {
	// VALIDATE STUDENT RANK
	if filter.StudentRank != "" {
		rank := strings.TrimSpace(string(filter.StudentRank))
		if rank != "excellent" && rank != "good" && rank != "average" && rank != "weak" {
			return ErrStudentRank
		}
	}
	return nil
}

func (s *StudentService) BulkAddStudents(students []*studentModels.Student) error {
	if len(students) == 0 {
		return nil
	}

	// Normalize and validate all students
	validatedStudents := make([]*studentModels.Student, len(students))
	emailSet := make(map[string]bool) // Track emails within input for duplicates

	for i, student := range students {
		if student == nil {
			return ErrStudentData
		}

		student = normalizeStudent(student)

		if err := validateStudent(student); err != nil {
			return err
		}

		// Check for duplicate emails within the input
		email := strings.ToLower(student.Email)
		if emailSet[email] {
			return ErrStudentEmailExisted
		}
		emailSet[email] = true

		validatedStudents[i] = student
	}

	// Check for email conflicts with existing students
	for _, student := range validatedStudents {
		existed, err := s.repo.GetStudentByEmail(student.Email)
		if err == nil && existed != nil {
			return ErrStudentEmailExisted
		}
	}

	// Generate UUIDs for students without IDs
	for _, student := range validatedStudents {
		if student.ID == "" {
			student.ID = uuid.New().String()
		}
	}

	return s.repo.BulkAddStudents(validatedStudents)
}
