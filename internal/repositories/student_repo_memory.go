// Step 3:  implement interface
package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/student-management/internal/models"
)

type InMemoStudentRepo struct {
	students map[string]*models.Student
	filePath string // Read/write JSON - Read file when initialized - Write file after Add/update/delete
}

// LOAD JSON
func (r *InMemoStudentRepo) loadFile() error {

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			r.students = make(map[string]*models.Student)
			return nil
		}
		return err
	}

	var data []*models.Student
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	r.students = make(map[string]*models.Student)
	for _, value := range data {
		r.students[value.ID] = value
	}

	return nil
}

// SAVE JSON
func (r *InMemoStudentRepo) saveFile() error {

	var studentData []*models.Student
	for _, student := range r.students {
		studentData = append(studentData, student)
	}

	students, err := json.MarshalIndent(studentData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, students, 0644)
}

// INITIALIZE EMPTY CONSTRUCTOR
func NewStudentMemoryRepo(filePath string) (*InMemoStudentRepo, error) {
	repo := &InMemoStudentRepo{
		students: make(map[string]*models.Student),
		filePath: filePath,
	}

	// Load JSON data
	if err := repo.loadFile(); err != nil {
		return nil, err
	}
	return repo, nil
}

// ADD STUDENT
func (r *InMemoStudentRepo) AddStudent(student *models.Student) error {

	if _, existed := r.students[student.ID]; existed {
		return fmt.Errorf("student with ID %s existed", student.ID)
	}

	r.students[student.ID] = student
	return r.saveFile()
}

// UPDATE STUDENT
func (r *InMemoStudentRepo) UpdateStudent(student *models.Student) error {
	if _, existed := r.students[student.ID]; !existed {
		return fmt.Errorf("student with ID %s not existed", student.ID)
	}

	r.students[student.ID] = student
	return r.saveFile()
}

// DELETE STUDENT
func (r *InMemoStudentRepo) DeleteStudent(studentID string) error {
	if _, existed := r.students[studentID]; !existed {
		return fmt.Errorf("student with ID %s does not existed", studentID)
	}
	delete(r.students, studentID)

	return r.saveFile()
}

// GET ALL STUDENT
func (r *InMemoStudentRepo) GetAllStudent() ([]*models.Student, error) {
	students := make([]*models.Student, 0, len(r.students))

	for _, s := range r.students {
		students = append(students, s)
	}
	return students, nil
}

// GET STUDENT BY ID
func (r *InMemoStudentRepo) GetStudentByID(studentID string) (*models.Student, error) {
	student, existed := r.students[studentID]
	if !existed {
		return nil, fmt.Errorf("student with ID %s does not existed", studentID)
	}

	return student, nil
}

// GET STUDENT BY EMAIL
func (r *InMemoStudentRepo) GetStudentByEmail(studentEmail string) (*models.Student, error) {
	for _, student := range r.students {
		if strings.EqualFold(student.Email, studentEmail) {
			return student, nil
		}
	}
	return nil, fmt.Errorf("student with Email %s does not existed", studentEmail)
}

// ADD SCORE
func (r *InMemoStudentRepo) AddScore(studentID string, score *models.SubjectScore) error {

	student, existed := r.students[studentID]

	if !existed {
		return fmt.Errorf("student with ID %s does not existed", studentID)
	}

	student.Scores = append(student.Scores, score)

	return r.saveFile()
}

// UPDATE SCORE
func (r *InMemoStudentRepo) UpdateScore(studentID string, score *models.SubjectScore) error {
	student, existed := r.students[studentID]

	if !existed {
		return fmt.Errorf("student with ID %s does not existed", studentID)
	}

	for _, s := range student.Scores{
		if strings.EqualFold(s.Subject, score.Subject){
			s.Score = score.Score
			return r.saveFile()
		}
	}
	return fmt.Errorf("subject %s does not existed", score.Subject)
}
