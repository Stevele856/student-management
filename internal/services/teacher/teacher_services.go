package teacher

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/student-management/internal/models"
	teacherModels "github.com/student-management/internal/models/teacher"
	teacherRepo "github.com/student-management/internal/repositories/teacher"
	"github.com/student-management/pkg/utils"
)

type TeacherService struct {
	repo teacherRepo.TeacherRepository
}

// CONSTRUCTOR

func NewTeacherService(repo teacherRepo.TeacherRepository) *TeacherService {
	return &TeacherService{repo: repo}
}

var (
	ErrTeacherData              = errors.New("invalid teacher data")
	ErrNameRequired             = errors.New("teacher name is required")
	ErrEmailRequired            = errors.New("teacher email is required")
	ErrNameFormat               = errors.New("invalid teacher name format")
	ErrEmailFormat              = errors.New("error teacher email format")
	ErrValidDOB                 = errors.New("date of birth cannot be in the future")
	ErrTeacherGender            = errors.New("teacher gender should be male or female")
	ErrTeacherStatus            = errors.New("missing teacher status")
	ErrEmployeeID               = errors.New("invalid teacher identity number")
	ErrFormatPhoneNumber        = errors.New("invalid phone number format")
	ErrSubjectFormat            = errors.New("invalid subject format")
	ErrClassFormat              = errors.New("invalid teacher class assigned")
	ErrSubjectRequired          = errors.New("subject is require for teacher")
	ErrTeacherHireDate          = errors.New("invalid teacher hire date")
	ErrClassRequired            = errors.New("teacher's class required")
	ErrSubjectDuplicated        = errors.New("subject already existed")
	ErrClassDuplicate           = errors.New("class being dublicated for this teacher")
	ErrTeacherEmailExisted      = errors.New("teacher email already existed")
	ErrTeacherEmployeeIDExisted = errors.New("teacher employee ID already existed")
	ErrTeacherIDRequired = errors.New("teacher ID required")
)

// ADD TEACHER
func (t *TeacherService) AddTeacher(teacher *teacherModels.Teacher) error {
	if teacher == nil {
		return ErrTeacherData
	}

	teacher = normalizeTeacher(teacher)
	if err := validateTeacher(teacher); err != nil {
		return err
	}

	existed, err := t.repo.GetTeacherByEmail(teacher.Email)
	if err == nil && existed != nil {
		return ErrTeacherEmailExisted
	}
	// DB/Network
	if err != nil && !errors.Is(err, teacherRepo.ErrTeacherNotFound) {
		return err
	}

	// Unique Employee ID
	existed, err = t.repo.GetTeacherByEmployeeID(teacher.EmployeeID)
	if err != nil && !errors.Is(err, teacherRepo.ErrTeacherNotFound) {
		return err
	}
	if err == nil && existed != nil {
		return ErrTeacherEmployeeIDExisted
	}

	if teacher.ID == "" {
		teacher.ID = uuid.New().String()
	}

	now := time.Now().UTC()

	if teacher.CreatedAt.IsZero() {
		teacher.CreatedAt = now
	}

	teacher.UpdatedAt = now

	return t.repo.AddTeacher(teacher)
}

// UPDATE TEACHER
func (t *TeacherService) UpdateTeacher(teacher *teacherModels.Teacher) error {
	if teacher == nil {
		return ErrTeacherData
	}

	teacher.ID = strings.TrimSpace(teacher.ID)
	if teacher.ID == ""{
		return ErrTeacherIDRequired
	}
	
	existing, err := t.repo.GetTeacherByID(teacher.ID)
	if err != nil{
		return err
	}

	teacher = normalizeTeacher(teacher)

	// Preserve system fields from old record
	teacher.ID = existing.ID
	teacher.CreatedAt = existing.CreatedAt
	if teacher.PasswordHash == ""{
		teacher.PasswordHash = existing.PasswordHash
	}

	if err := validateTeacher(teacher); err != nil {
		return err
	}

	existedEmail, err := t.repo.GetTeacherByEmail(teacher.Email)
	if err != nil && !errors.Is(err, teacherRepo.ErrTeacherNotFound){
		return err
	}

	if err == nil && existedEmail != nil && existedEmail.ID != teacher.ID{
		return ErrTeacherEmailExisted
	}

	existedEmployeeID, err := t.repo.GetTeacherByEmployeeID(teacher.EmployeeID)
	if err != nil && !errors.Is(err, teacherRepo.ErrTeacherNotFound){
		return err
	}

	if err == nil && existedEmployeeID != nil && existedEmployeeID.ID != teacher.EmployeeID{
		return ErrTeacherEmployeeIDExisted
	}

	teacher.UpdatedAt = time.Now().UTC()

	return t.repo.UpdateTeacher(teacher)

}

// STANDARD DATA
func normalizeTeacher(teacher *teacherModels.Teacher) *teacherModels.Teacher {
	if teacher == nil {
		return nil
	}
	t := *teacher
	t.ID = strings.TrimSpace(t.ID)
	t.FullName = strings.TrimSpace(t.FullName)
	t.Email = strings.ToLower(strings.TrimSpace(t.Email))
	t.Address = strings.TrimSpace(t.Address)
	t.Phone = strings.TrimSpace(t.Phone)

	t.EmployeeID = strings.ToUpper(strings.TrimSpace(t.EmployeeID))
	for i := range t.SubjectTaught {
		t.SubjectTaught[i] = strings.TrimSpace(t.SubjectTaught[i])
	}

	for i := range t.ClassAssigned {
		t.ClassAssigned[i] = strings.TrimSpace(t.ClassAssigned[i])
	}
	t.Status = teacherModels.TeacherStatus(strings.ToLower(strings.TrimSpace(string(t.Status))))

	return &t
}

func validateTeacher(teacher *teacherModels.Teacher) error {
	if teacher == nil {
		return ErrTeacherData
	}

	if teacher.FullName == "" {
		return ErrNameRequired
	}

	if teacher.Email == "" {
		return ErrEmailRequired
	}

	if !utils.IsValidEmail(teacher.Email) {
		return ErrEmailFormat
	}

	if teacher.Gender != models.GenderMale && teacher.Gender != models.GenderFemale {
		return ErrTeacherGender
	}

	if !utils.IsValidName(teacher.FullName) {
		return ErrNameFormat
	}

	if teacher.DateOfBirth.IsZero() || teacher.DateOfBirth.After(time.Now()) {
		return ErrValidDOB
	}

	if teacher.Status != teacherModels.Active && teacher.Status != teacherModels.Inactive && teacher.Status != teacherModels.OnLeave {
		return ErrTeacherStatus
	}

	if teacher.EmployeeID == "" || !utils.IsValidEmployeeID(teacher.EmployeeID) {
		return ErrEmployeeID
	}

	if teacher.Phone == "" || !utils.IsValidPhoneNumber(teacher.Phone) {
		return ErrFormatPhoneNumber
	}

	if len(teacher.SubjectTaught) == 0 {
		return ErrSubjectRequired
	}

	// Check dublicate subject
	uniqueSubject := make(map[string]struct{}, len(teacher.SubjectTaught))

	for _, subj := range teacher.SubjectTaught {
		if subj == "" {
			return ErrSubjectRequired
		}
		if !utils.IsValidSubject(subj) {
			return ErrSubjectFormat
		}

		subjectKey := strings.ToLower(strings.TrimSpace(subj))
		if _, existed := uniqueSubject[subjectKey]; existed {
			return ErrSubjectDuplicated
		}
		uniqueSubject[subjectKey] = struct{}{}
	}

	if len(teacher.ClassAssigned) == 0 {
		return ErrClassRequired
	}

	// Check if dublicate Class
	uniqueClassAssigned := make(map[string]struct{}, len(teacher.ClassAssigned))

	for _, class := range teacher.ClassAssigned {
		if class == "" {
			return ErrClassRequired
		}
		if !utils.IsValidClass(class) {
			return ErrClassFormat
		}
		classKey := strings.ToUpper(strings.TrimSpace(class))
		if _, existed := uniqueClassAssigned[classKey]; existed {
			return ErrClassDuplicate
		}
		uniqueClassAssigned[classKey] = struct{}{}
	}

	if teacher.HireDate.After(time.Now()) {
		return ErrTeacherHireDate
	}

	if teacher.HireDate.Before(teacher.DateOfBirth) {
		return ErrTeacherHireDate
	}

	if !IsValidHireTeacher(25, teacher) {
		return ErrTeacherHireDate
	}

	return nil
}

func IsValidHireTeacher(minAge int, teacher *teacherModels.Teacher) bool {
	if teacher == nil || teacher.DateOfBirth.IsZero() || teacher.HireDate.IsZero() {
		return false
	}

	validAge := teacher.DateOfBirth.AddDate(minAge, 0, 0)
	return !teacher.HireDate.Before(validAge)
}
