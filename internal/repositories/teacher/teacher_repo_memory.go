package teacherRepo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	teacherModels "github.com/student-management/internal/models/teacher"
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
	var teacherData []*teacherModels.Teacher

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
	return r.saveFile()
}

func (r *InMemoTeacherRepo) UpdateTeacher(teacher *teacherModels.Teacher) error {
	if _, existed := r.teachers[teacher.ID]; !existed {
		return fmt.Errorf("teacher with ID %s does not existed", teacher.ID)
	}

	r.teachers[teacher.ID] = teacher

	return r.saveFile()
}

func (r *InMemoTeacherRepo) DeleteTeacher(teacherID string) error {
	if _, existed := r.teachers[teacherID]; !existed {
		return fmt.Errorf("teacher with ID %s does not existed", teacherID)
	}

	delete(r.teachers, teacherID)
	return r.saveFile()
}

func (r *InMemoTeacherRepo) GetAllTeacher() ([]*teacherModels.Teacher, error) {
	teachers := make([]*teacherModels.Teacher, 0, len(r.teachers))

	for _, teacher := range r.teachers {
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func (r *InMemoTeacherRepo) GetTeacherByID(teacherID string) (*teacherModels.Teacher, error) {
	teacher, existed := r.teachers[teacherID]

	if !existed {
		return nil, fmt.Errorf("teacher with ID %s does not existed", teacherID)
	}

	return teacher, nil
}

func (r *InMemoTeacherRepo) GetTeacherByEmail(teacherEmail string) (*teacherModels.Teacher, error) {
	for _, teacher := range r.teachers {
		if strings.EqualFold(teacher.Email, teacherEmail) {
			return teacher, nil
		}
	}

	return nil, fmt.Errorf("teacher with email %s does not existed", teacherEmail)
}

func (r *InMemoTeacherRepo) GetTeacherAssignedBySubject(subject string) ([]*teacherModels.Teacher, error) {
	result := []*teacherModels.Teacher{}

	for _, teacher := range r.teachers {
		for _, subj := range teacher.SubjectTaught {
			if strings.EqualFold(subj, subject) {
				result = append(result, teacher)
				break
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no teachers found teaching subject %s", subject)
	}

	return result, nil
}
