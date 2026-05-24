package teacherRepo

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	teacherModels "github.com/student-management/internal/models/teacher"
	"github.com/student-management/internal/predicates"
)

type InMemoTeacherRepo struct {
	teachers map[string]*teacherModels.Teacher
	filePath string
}

// LOAD JSON
func (r *InMemoTeacherRepo) loadFile() error {
	file, err := os.ReadFile(r.filePath)

	if err != nil {
		if os.IsNotExist(err) {
			r.teachers = make(map[string]*teacherModels.Teacher)
			return nil
		}
		return err
	}

	var teacherData []*teacherModels.Teacher
	if err := json.Unmarshal(file, &teacherData); err != nil {
		return err
	}

	r.teachers = make(map[string]*teacherModels.Teacher)
	for _, teacher := range teacherData {
		r.teachers[teacher.ID] = teacher
	}
	return nil
}

// SAVE JSON
func (r *InMemoTeacherRepo) saveFile() error {
	teacherData := []*teacherModels.Teacher{}

	for _, teacher := range r.teachers {
		teacherData = append(teacherData, teacher)
	}

	teachers, err := json.MarshalIndent(teacherData, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, teachers, 0644)
}

// EMPTY CONSTRUCTOR
func NewTeacherMemoryRepo(filePath string) (*InMemoTeacherRepo, error) {
	repo := &InMemoTeacherRepo{
		teachers: make(map[string]*teacherModels.Teacher),
		filePath: filePath,
	}

	if err := repo.loadFile(); err != nil {
		return nil, err
	}

	return repo, nil

}

// CRUD
func (r *InMemoTeacherRepo) AddTeacher(teacher *teacherModels.Teacher) error {
	if _, existed := r.teachers[teacher.ID]; existed {
		return fmt.Errorf("teacher with ID %s already exsited", teacher.ID)
	}

	r.teachers[teacher.ID] = teacher
	if err := r.saveFile(); err != nil {
		return fmt.Errorf("save teacher add: %w", err)
	}
	return nil
}

func (r *InMemoTeacherRepo) UpdateTeacher(teacher *teacherModels.Teacher) error {
	if _, existed := r.teachers[teacher.ID]; !existed {
		return ErrTeacherNotFound
	}

	r.teachers[teacher.ID] = teacher
	if err := r.saveFile(); err != nil {
		return fmt.Errorf("save teacher update: %w", err)
	}
	return nil
}

func (r *InMemoTeacherRepo) DeleteTeacher(teacherID string) error {
	if _, existed := r.teachers[teacherID]; !existed {
		return ErrTeacherNotFound
	}

	delete(r.teachers, teacherID)
	if err := r.saveFile(); err != nil {
		return fmt.Errorf("save teacher delete: %w", err)
	}
	return nil
}

func (r *InMemoTeacherRepo) GetAllTeachers() ([]*teacherModels.Teacher, error) {
	teachers := make([]*teacherModels.Teacher, 0, len(r.teachers))

	for _, teacher := range r.teachers {
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func (r *InMemoTeacherRepo) GetTeacherByID(teacherID string) (*teacherModels.Teacher, error) {
	teacher, existed := r.teachers[teacherID]

	if !existed {
		return nil, ErrTeacherNotFound
	}

	return teacher, nil
}

func (r *InMemoTeacherRepo) GetTeacherByEmail(teacherEmail string) (*teacherModels.Teacher, error) {
	for _, teacher := range r.teachers {
		if strings.EqualFold(teacher.Email, teacherEmail) {
			return teacher, nil
		}
	}

	return nil, ErrTeacherNotFound
}

func (r *InMemoTeacherRepo) GetTeacherAssignedBySubject(subjectAssigned string) ([]*teacherModels.Teacher, error) {
	result := []*teacherModels.Teacher{}

	for _, teacher := range r.teachers {
		for _, subj := range teacher.SubjectTaught {
			if strings.EqualFold(subj, subjectAssigned) {
				result = append(result, teacher)
				break
			}
		}
	}

	return result, nil
}

func (r *InMemoTeacherRepo) GetTeacherByAssignedClass(classAssigned string) ([]*teacherModels.Teacher, error) {
	result := []*teacherModels.Teacher{}

	for _, teacher := range r.teachers {
		for _, class := range teacher.ClassAssigned {
			if strings.EqualFold(class, classAssigned) {
				result = append(result, teacher)
				break
			}
		}
	}

	return result, nil
}

func (r *InMemoTeacherRepo) GetTeacherByStatus(status teacherModels.TeacherStatus) ([]*teacherModels.Teacher, error) {
	result := []*teacherModels.Teacher{}

	for _, teacher := range r.teachers {
		if teacher.Status == status {
			result = append(result, teacher)
		}
	}

	return result, nil
}

func (r *InMemoTeacherRepo) GetTeacherByEmployeeID(employeeID string) (*teacherModels.Teacher, error) {
	for _, teacher := range r.teachers {
		if strings.EqualFold(teacher.EmployeeID, employeeID) {
			return teacher, nil
		}
	}

	return nil, ErrTeacherNotFound
}

func (r *InMemoTeacherRepo) FilterTeachers(p predicates.PredicateTeacher) ([]*teacherModels.Teacher, error) {
	if p == nil {
		return r.GetAllTeachers()
	}

	result := []*teacherModels.Teacher{}

	for _, teacher := range r.teachers {
		if !p(teacher) {
			continue
		}
		result = append(result, teacher)
	}

	return result, nil
}

func (r *InMemoTeacherRepo) GetTeachersPaginated(page, pageSize int) ([]*teacherModels.Teacher, int, error) {
	if page < 1 {
		return nil, 0, ErrInvalidPage
	}
	if pageSize < 1 {
		return nil, 0, ErrInvalidPageSize
	}

	teachers := make([]*teacherModels.Teacher, 0, len(r.teachers))
	for _, teacher := range r.teachers {
		teachers = append(teachers, teacher)
	}

	// map iteration trong Go không ổn định, nên sort để pagination nhất quán
	sort.Slice(teachers, func(i, j int) bool {
		return teachers[i].ID < teachers[j].ID
	})

	total := len(teachers)
	start := (page - 1) * pageSize

	if start > total {
		return []*teacherModels.Teacher{}, total, nil
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return teachers[start:end], total, nil
}

func (r *InMemoTeacherRepo) BulkAddTeachers(teachers []*teacherModels.Teacher) error {
	if len(teachers) == 0 {
		return nil
	}
	// Check if the incoming teacher data(bulk request) has the same ID, not overwrite the old data with same ID. "payload-duplicate detection"
	seenID := make(map[string]struct{}, len(teachers))
	for _, teacher := range teachers {
		if _, ok := seenID[teacher.ID]; ok {
			return fmt.Errorf("%w: duplicate ID in payload: %s", ErrTeacherAlreadyExists, teacher.ID)
		}
		seenID[teacher.ID] = struct{}{}

		if _, existed := r.teachers[teacher.ID]; existed {
			return fmt.Errorf("%w: %s", ErrTeacherAlreadyExists, teacher.ID)
		}
	}

	for _, teacher := range teachers {
		r.teachers[teacher.ID] = teacher
	}

	if err := r.saveFile(); err != nil {
		return fmt.Errorf("save bulk add teachers: %w", err)
	}
	return nil

}
