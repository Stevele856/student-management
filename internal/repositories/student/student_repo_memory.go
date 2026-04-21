// Step 3:  implement interface
package studentRepo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/student-management/internal/models/student"
	"github.com/student-management/internal/predicate"
)

// CHECK IMPLEMENTATION FUNCTION WHETHER IT MATCH WITH INTERFACE
// var _ StudentRepository = &InMemoStudentRepo{}

type InMemoStudentRepo struct {
	students map[string]*studentModels.Student
	filePath string
	// Read/write JSON - Read file when initialized - Write file after Add/update/delete
}

// LOAD JSON
func (r *InMemoStudentRepo) loadFile() error {

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			r.students = make(map[string]*studentModels.Student)
			return nil
		}
		return err
	}

	var data []*studentModels.Student
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	r.students = make(map[string]*studentModels.Student)
	for _, value := range data {
		r.students[value.ID] = value
	}

	return nil
}

// SAVE JSON
func (r *InMemoStudentRepo) saveFile() error {

	var studentData []*studentModels.Student
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
		students: make(map[string]*studentModels.Student),
		filePath: filePath,
	}

	// Load JSON data
	if err := repo.loadFile(); err != nil {
		return nil, err
	}
	return repo, nil
}

// CRUD STUDENT
func (r *InMemoStudentRepo) AddStudent(student *studentModels.Student) error {

	if _, existed := r.students[student.ID]; existed {
		return fmt.Errorf("student with ID %s existed", student.ID)
	}

	r.students[student.ID] = student
	return r.saveFile()
}

func (r *InMemoStudentRepo) UpdateStudent(student *studentModels.Student) error {
	if _, existed := r.students[student.ID]; !existed {
		return fmt.Errorf("student with ID %s not existed", student.ID)
	}

	r.students[student.ID] = student
	return r.saveFile()
}

func (r *InMemoStudentRepo) DeleteStudent(studentID string) error {
	if _, existed := r.students[studentID]; !existed {
		return fmt.Errorf("student with ID %s does not exist", studentID)
	}
	delete(r.students, studentID)

	return r.saveFile()
}

func (r *InMemoStudentRepo) GetAllStudents() ([]*studentModels.Student, error) {
	students := make([]*studentModels.Student, 0, len(r.students))

	for _, s := range r.students {
		students = append(students, s)
	}
	return students, nil
}

func (r *InMemoStudentRepo) GetStudentByID(studentID string) (*studentModels.Student, error) {
	student, existed := r.students[studentID]
	if !existed {
		return nil, fmt.Errorf("student with ID %s does not exist", studentID)
	}

	return student, nil
}

func (r *InMemoStudentRepo) GetStudentByEmail(studentEmail string) (*studentModels.Student, error) {
	for _, student := range r.students {
		if strings.EqualFold(student.Email, studentEmail) {
			return student, nil
		}
	}
	return nil, fmt.Errorf("student with Email %s does not exist", studentEmail)
}

/* --------------------------- */

// IMPLEMENT CRUD SCORE
func (r *InMemoStudentRepo) AddScore(studentID string, score *studentModels.SubjectScore) error {

	student, existed := r.students[studentID]

	if !existed {
		return fmt.Errorf("student with ID %s does not exist", studentID)
	}

	student.Scores = append(student.Scores, score)

	return r.saveFile()
}

func (r *InMemoStudentRepo) UpdateScore(studentID string, score *studentModels.SubjectScore) error {
	student, existed := r.students[studentID]

	if !existed {
		return fmt.Errorf("student with ID %s does not exist", studentID)
	}

	for i, s := range student.Scores {
		if strings.EqualFold(s.Subject, score.Subject) {
			// s.Score = score.Score -> s is copy, not pointer
			student.Scores[i].Score = score.Score
			return r.saveFile()
		}
	}
	return fmt.Errorf("subject %s does not exist", score.Subject)
}

func (r *InMemoStudentRepo) DeleteScore(studentID, subject string) error {
	student, existed := r.students[studentID]

	if !existed {
		return fmt.Errorf("student with ID %s does not exist", studentID)
	}

	for i, s := range student.Scores {
		if strings.EqualFold(s.Subject, subject) {
			student.Scores = append(student.Scores[:i], student.Scores[i+1:]...)
			return r.saveFile()
		}

	}
	return fmt.Errorf("subject %s does not exist", subject)
}

/*
student.Scores = []*student.SubjectScore{
	{Subject: "Toan", Score: 7.5}		i
	{Subject: "Tieng Anh", Score: 8}	i+1
	{Subject: "Tieng Viet", Score: 6}	i+2

	- [:i] - Lấy từ đầu đến trước index i
	- [i+1:]... - Lấy từ index i+1 đến hết
}
*/

func (r *InMemoStudentRepo) GetScoresByStudentID(studentID string) ([]*studentModels.SubjectScore, error) {

	student, existed := r.students[studentID]

	if !existed {
		return nil, fmt.Errorf("student with ID %s does not exist", studentID)
	}

	return student.Scores, nil
}

func (r *InMemoStudentRepo) GetScoresBySubject(studentID, subject string) (*studentModels.SubjectScore, error) {
	student, existed := r.students[studentID]

	if !existed {
		return nil, fmt.Errorf("student with ID %s does not exist", studentID)
	}

	for _, s := range student.Scores {
		if strings.EqualFold(s.Subject, subject) {
			return s, nil
		}
	}
	return nil, fmt.Errorf("subject %s does not exist", subject)
}

/* ------------------------- */

// PREDICATE FOR FILTER
func (r *InMemoStudentRepo) FilterStudents(p predicate.PredicateStudent) ([]*studentModels.Student, error) {
	result := []*studentModels.Student{}

	for _, s := range r.students {
		if !p(s) {
			continue
		}
		result = append(result, s)
	}
	return result, nil
}

/* ----------------------------- */

// BULK ADD STUDENTS (CSV)
func (r *InMemoStudentRepo) BulkAddStudents(students []*studentModels.Student) error {
	if len(students) == 0 {
		return nil
	}

	// Check for duplicate IDs in the input and existing students
	for _, student := range students {
		if _, existed := r.students[student.ID]; existed {
			return fmt.Errorf("student with ID %s already exists", student.ID)
		}
	}

	// Add all students to the map
	for _, student := range students {
		r.students[student.ID] = student
	}

	// Save all students to file
	return r.saveFile()
}
