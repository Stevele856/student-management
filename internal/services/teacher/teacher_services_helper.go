package teacher

import (
	"errors"
	"strings"
	"time"

	"github.com/student-management/internal/models"
	teacherModels "github.com/student-management/internal/models/teacher"
	teacherRepo "github.com/student-management/internal/repositories/teacher"
	"github.com/student-management/pkg/utils"
)

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

func (t *TeacherService) ensureTeacherUnique(email, employeeID, currentTeacherID string) error {
	// Unique Teacher email
	existed, err := t.repo.GetTeacherByEmail(email)

	// DB/Network
	if err != nil && !errors.Is(err, teacherRepo.ErrTeacherNotFound) {
		return err
	}

	if err == nil && existed != nil && existed.ID != currentTeacherID {
		return ErrTeacherEmailExisted
	}

	// Unique Employee ID
	existed, err = t.repo.GetTeacherByEmployeeID(employeeID)
	
	// DB/Network
	if err != nil && !errors.Is(err, teacherRepo.ErrTeacherNotFound) {
		return err
	}
	if err == nil && existed != nil && existed.ID != currentTeacherID  {
		return ErrTeacherEmployeeIDExisted
	}

	return nil
}

func (t *TeacherService) requireTeacherID(teacherID string) (string, error){
	teacherID = strings.TrimSpace(teacherID)
	if teacherID == ""{
		return "", ErrTeacherIDRequired
	}
	return teacherID, nil
}

func (t *TeacherService) ensureTeacherExistsByID(id string) (*teacherModels.Teacher, error){
	existing, err := t.repo.GetTeacherByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil,teacherRepo.ErrTeacherNotFound
	}
	return existing, nil
}
