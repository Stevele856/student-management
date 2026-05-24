package teacher

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/student-management/internal/models"
	teacherModels "github.com/student-management/internal/models/teacher"
	"github.com/student-management/internal/predicates"

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
	if err == nil && existed != nil && existed.ID != currentTeacherID {
		return ErrTeacherEmployeeIDExisted
	}

	return nil
}

func (t *TeacherService) requireTeacherID(teacherID string) (string, error) {
	teacherID = strings.TrimSpace(teacherID)
	if teacherID == "" {
		return "", ErrTeacherIDRequired
	}
	return teacherID, nil
}

func (t *TeacherService) ensureTeacherExistsByID(id string) (*teacherModels.Teacher, error) {
	existing, err := t.repo.GetTeacherByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, teacherRepo.ErrTeacherNotFound
	}
	return existing, nil
}

func normalizeFilterTeacher(filter *teacherModels.FilterTeachers) *teacherModels.FilterTeachers {
	if filter == nil {
		return nil
	}
	tf := *filter
	tf.Name = strings.TrimSpace(tf.Name)
	tf.Gender = models.Gender(strings.ToLower(string(tf.Gender)))
	tf.Email = strings.ToLower(strings.TrimSpace(tf.Email))
	tf.Status = teacherModels.TeacherStatus(strings.ToLower(strings.TrimSpace(string(tf.Status))))
	tf.EmployeeID = strings.ToUpper(strings.TrimSpace(tf.EmployeeID))
	for i := range tf.SubjectTaught {
		tf.SubjectTaught[i] = strings.TrimSpace(tf.SubjectTaught[i])
	}

	for i := range tf.ClassAssigned {
		tf.ClassAssigned[i] = strings.TrimSpace(tf.ClassAssigned[i])
	}

	return &tf
}

func validateFilterTeachers(filter *teacherModels.FilterTeachers) error {
	if filter.Name != "" && !utils.IsValidName(filter.Name) {
		return ErrNameFormat
	}
	if filter.Gender != "" && filter.Gender != models.GenderMale && filter.Gender != models.GenderFemale {
		return ErrTeacherGender
	}
	if filter.Email != "" && !utils.IsValidEmail(filter.Email) {
		return ErrEmailFormat
	}
	if filter.EmployeeID != "" && !utils.IsValidEmployeeID(filter.EmployeeID) {
		return ErrEmployeeID
	}
	if filter.Status != "" && filter.Status != teacherModels.Active && filter.Status != teacherModels.Inactive && filter.Status != teacherModels.OnLeave {
		return ErrTeacherStatus
	}
	for _, class := range filter.ClassAssigned {
		if !utils.IsValidClass(class) {
			return ErrClassFormat
		}
	}

	for _, subject := range filter.SubjectTaught {
		if !utils.IsValidSubject(subject) {
			return ErrSubjectFormat
		}
	}

	if filter.HireDateFrom != nil && filter.HireDateTo != nil {
		if filter.HireDateFrom.After(*filter.HireDateTo) {
			return errors.New("hire date from cannot be after hire date to")
		}
	}
	return nil
}

func filterPredicateTeachers(filter *teacherModels.FilterTeachers) []predicates.PredicateTeacher {
	predicate := []predicates.PredicateTeacher{}

	if filter.Name != "" {
		predicate = append(predicate, predicates.ByTeacherName(filter.Name))
	}
	if filter.Gender != "" {
		predicate = append(predicate, predicates.ByTeacherGender(filter.Gender))
	}
	if filter.Email != "" {
		predicate = append(predicate, predicates.ByTeacherEmail(filter.Email))
	}
	if filter.EmployeeID != "" {
		predicate = append(predicate, predicates.ByEmployeeID(filter.EmployeeID))
	}
	if filter.Status != "" {
		predicate = append(predicate, predicates.ByStatus(filter.Status))
	}

	for _, subject := range filter.SubjectTaught {
		if subject != "" {
			predicate = append(predicate, predicates.BySubjectTaught(subject))
		}
	}

	for _, class := range filter.ClassAssigned {
		if class != "" {
			predicate = append(predicate, predicates.ByClassAssigned(class))
		}
	}

	if filter.HireDateFrom != nil || filter.HireDateTo != nil {
		predicate = append(predicate, predicates.ByHireDateRange(filter.HireDateFrom, filter.HireDateTo))
	}
	return predicate
}

func (s *TeacherService) BulkAddTeachers(teachers []*teacherModels.Teacher) error {
	if len(teachers) == 0 {
		return nil
	}

	// Normalize and validate all teachers
	validatedTeachers := make([]*teacherModels.Teacher, len(teachers))
	emailSet := make(map[string]struct{}, len(teachers)) // detecting duplicate emails inside the same input payload
	idSet := make(map[string]struct{}, len(teachers))
	employeeIDSet := make(map[string]struct{}, len(teachers))

	for i, teacher := range teachers {
		if teacher == nil {
			return ErrTeacherData
		}

		teacher = normalizeTeacher(teacher)

		if err := validateTeacher(teacher); err != nil {
			return err
		}

		// Duplicate email in payload
		email := strings.ToLower(strings.TrimSpace(teacher.Email))
		if _, exists := emailSet[email]; exists {
			return ErrTeacherEmailExisted
		}
		emailSet[email] = struct{}{}

		// Duplicate non-empty ID in payload
		if teacher.ID != "" {
			id := strings.TrimSpace(teacher.ID)
			if _, exists := idSet[id]; exists {
				return ErrTeacherIDRequired // better: ErrTeacherIDDuplicated if you define one
			}
			idSet[id] = struct{}{}
		}

		validatedTeachers[i] = teacher

		// Duplicate employee ID in payload
		employeeID := strings.ToLower(strings.TrimSpace(teacher.EmployeeID))
		if _, exists := employeeIDSet[employeeID]; exists {
			return ErrTeacherEmployeeIDExisted
		}
		employeeIDSet[employeeID] = struct{}{}
	}

	// Check conflicts with existing records
	for _, teacher := range validatedTeachers {
		existedByEmail, err := s.repo.GetTeacherByEmail(teacher.Email)
		if err == nil && existedByEmail != nil {
			return ErrTeacherEmailExisted
		}
		if err != nil && !errors.Is(err, teacherRepo.ErrTeacherNotFound) {
			return wrapRepoError(err)
		}

		existedByEmpID, err := s.repo.GetTeacherByEmployeeID(teacher.EmployeeID)
		if err == nil && existedByEmpID != nil {
			return ErrTeacherEmployeeIDExisted
		}
		if err != nil && !errors.Is(err, teacherRepo.ErrTeacherNotFound) {
			return wrapRepoError(err)
		}
	}

	// Generate UUID for empty IDs
	for _, teacher := range validatedTeachers {
		if strings.TrimSpace(teacher.ID) == "" {
			teacher.ID = uuid.New().String()
		}
	}

	if err := s.repo.BulkAddTeachers(validatedTeachers); err != nil {
		return wrapRepoError(err)
	}

	return nil
}

func wrapRepoError(err error) error {
	if errors.Is(err, teacherRepo.ErrTeacherNotFound) {
		return teacherRepo.ErrTeacherNotFound
	}
	if errors.Is(err, teacherRepo.ErrTeacherAlreadyExists) {
		return teacherRepo.ErrTeacherAlreadyExists
	}
	if errors.Is(err, teacherRepo.ErrInvalidPage) {
		return teacherRepo.ErrInvalidPage
	}
	if errors.Is(err, teacherRepo.ErrInvalidPageSize) {
		return teacherRepo.ErrInvalidPageSize
	}
	return err
}

