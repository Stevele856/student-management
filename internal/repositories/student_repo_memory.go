// Step 3:  Implement toàn bộ các method trong interface StudentRepository bằng in-memory (dùng map), cộng thêm phần đọc/ghi JSON.
package repositories

import (
	"encoding/json"
	"os"

	"github.com/student-management/internal/models"
)

type StudentMemoryRepo struct {
	students map[string]*models.Student
	filePath string 					// Read/write JSON - Read file when initialized - Write file after Add/update/delete
}

// Load JSON
func (r *StudentMemoryRepo) loadFile() error {
	// Check file if it existed
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			r.students = make(map[string]*models.Student)
			return nil
		}
		return err
	}

	// Initialized empty map to receive JSON data
	var students map[string]*models.Student
	if err := json.Unmarshal(file, &students); err != nil {
		return err
	}

	r.students = make(map[string]*models.Student)
	for _, value := range students {
		r.students[value.ID] = value
	}

	return nil

}

// Save JSON 
func (r *StudentMemoryRepo) saveFile() error {
	// Re-conver from map to slice 
	var students []*models.Student

	for _, value := range r.students{
		students = append(students, value)
	}

	data, err := json.MarshalIndent(students, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}
// Initilize empty StudentMemoryRepo
func NewStudentMemoryRepo(filePath string) (*StudentMemoryRepo, error) {
	repo := &StudentMemoryRepo{
		students: make(map[string]*models.Student),
		filePath: filePath,
	}

	// Load JSON data
	if err := repo.loadFile(); err != nil {
		return nil, err
	}
	return repo, nil
}
