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

// Load/Save JSON
func (r *StudentMemoryRepo) loadFile() error {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			r.students = make(map[string]*models.Student)
			return nil
		}
		return err
	}

	var students []*models.Student
	if err := json.Unmarshal(file, &students); err != nil {
		return err
	}

	r.students = make(map[string]*models.Student)
	for _, s := range students {
		r.students[s.ID] = s
	}

	return nil

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
