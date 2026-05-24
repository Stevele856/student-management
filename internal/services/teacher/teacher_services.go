package teacher

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	teacherModels "github.com/student-management/internal/models/teacher"
	"github.com/student-management/internal/predicates"
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
	ErrTeacherIDRequired        = errors.New("teacher ID required")
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

	if err := t.ensureTeacherUnique(teacher.Email, teacher.EmployeeID, ""); err != nil {
		return err
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

	teacherID, err := t.requireTeacherID(teacher.ID)
	if err != nil {
		return err
	}

	existingTeacher, err := t.ensureTeacherExistsByID(teacherID)
	if err != nil {
		return err
	}

	teacher = normalizeTeacher(teacher)

	// Preserve system fields from old record
	teacher.ID = existingTeacher.ID
	teacher.CreatedAt = existingTeacher.CreatedAt
	if teacher.PasswordHash == "" {
		teacher.PasswordHash = existingTeacher.PasswordHash
	}

	if err := validateTeacher(teacher); err != nil {
		return err
	}

	if err := t.ensureTeacherUnique(teacher.Email, teacher.EmployeeID, teacher.ID); err != nil {
		return err
	}

	teacher.UpdatedAt = time.Now().UTC()

	return t.repo.UpdateTeacher(teacher)

}

// DELELTE TEACHER
func (t *TeacherService) DeleteTeacher(teacherID string) error {
	teacherID, err := t.requireTeacherID(teacherID)
	if err != nil {
		return err
	}
	if _, err := t.ensureTeacherExistsByID(teacherID); err != nil {
		return err
	}

	return t.repo.DeleteTeacher(teacherID)
}

// GET ALL TEACHER
func (t *TeacherService) GetAllTeachers() ([]*teacherModels.Teacher, error) {
	return t.repo.GetAllTeachers()
}

// GET TEACHER BY ID
func (t *TeacherService) GetTeacherByID(teacherID string) (*teacherModels.Teacher, error) {
	teacherID, err := t.requireTeacherID(teacherID)
	if err != nil {
		return nil, err
	}
	existingTeacher, err := t.ensureTeacherExistsByID(teacherID)
	if err != nil {
		return nil, err
	}
	return existingTeacher, nil
}

func (t *TeacherService) GetTeacherByEmail(teacherEmail string) (*teacherModels.Teacher, error) {
	teacherEmail = strings.ToLower(strings.TrimSpace(teacherEmail))

	if teacherEmail == "" {
		return nil, ErrEmailRequired
	}

	if !utils.IsValidEmail(teacherEmail) {
		return nil, ErrEmailFormat
	}

	existing, err := t.repo.GetTeacherByEmail(teacherEmail)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, teacherRepo.ErrTeacherNotFound
	}

	return existing, nil

}

// List query should return [] + nil - repo did it
func (t *TeacherService) GetTeacherAssignedBySubject(subjectAssigned string) ([]*teacherModels.Teacher, error) {
	subjectAssigned = strings.ToLower(strings.TrimSpace(subjectAssigned))
	if subjectAssigned == "" {
		return nil, ErrSubjectRequired
	}
	if !utils.IsValidSubject(subjectAssigned) {
		return nil, ErrSubjectFormat
	}

	return t.repo.GetTeacherAssignedBySubject(subjectAssigned)
}

func (t *TeacherService) GetTeacherByAssignedClass(classAssigned string) ([]*teacherModels.Teacher, error) {
	classAssigned = strings.ToUpper(strings.TrimSpace(classAssigned))
	if classAssigned == "" {
		return nil, ErrClassRequired
	}
	if !utils.IsValidClass(classAssigned) {
		return nil, ErrClassFormat
	}

	return t.repo.GetTeacherByAssignedClass(classAssigned)
}

func (t *TeacherService) GetTeacherByStatus(status teacherModels.TeacherStatus) ([]*teacherModels.Teacher, error) {
	status = teacherModels.TeacherStatus(strings.ToLower(strings.TrimSpace(string(status))))
	if status != teacherModels.Active && status != teacherModels.Inactive && status != teacherModels.OnLeave {
		return nil, ErrTeacherStatus
	}
	return t.repo.GetTeacherByStatus(status)
}

func (t *TeacherService) GetTeacherByEmployeeID(employeeID string) (*teacherModels.Teacher, error) {
	employeeID = strings.ToUpper(strings.TrimSpace(employeeID))

	if !utils.IsValidEmployeeID(employeeID) {
		return nil, ErrEmployeeID
	}

	existing, err := t.repo.GetTeacherByEmployeeID(employeeID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, teacherRepo.ErrTeacherNotFound
	}

	return existing, nil
}


func (t *TeacherService) FilterTeachers(filter *teacherModels.FilterTeachers) ([]*teacherModels.Teacher, error){
	if filter == nil {
		return t.repo.GetAllTeachers()
	}
	filter = normalizeFilterTeacher(filter)

	if err := validateFilterTeachers(filter); err != nil {
		return nil, err
	}

	predicate := filterPredicateTeachers(filter)
	return t.repo.FilterTeachers(predicates.AndTeacher(predicate...))
	
}

func (t *TeacherService) GetTeachersPaginated(page, pageSize int) ([]*teacherModels.Teacher, int, error) {
	if page < 1 {
		return nil, 0, teacherRepo.ErrInvalidPage
	}
	if pageSize < 1 {
		return nil, 0, teacherRepo.ErrInvalidPageSize
	}

	teachers, total, err := t.repo.GetTeachersPaginated(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}
