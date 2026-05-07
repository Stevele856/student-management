package teacherModels

import (
	"time"

	"github.com/student-management/internal/models"
)

type TeacherStatus string

const (
	Active   TeacherStatus = "active"
	Inactive TeacherStatus = "inactive"
	OnLeave  TeacherStatus = "on-leave"
)

type Teacher struct {
	// Basic info
	ID          string        `json:"id"`
	FullName    string        `json:"full_name"`
	Email       string        `json:"email"`
	DateOfBirth time.Time     `json:"date_of_birth"`
	Gender      models.Gender `json:"gender"`
	Address     string        `json:"address"`
	Phone       string        `json:"phone"`

	// Professional info
	EmployeeID    string   `json:"employee_id"`
	SubjectTaught []string `json:"subject_taught"`
	ClassAssigned []string `json:"class_assigned"`

	// metadata field
	HireDate time.Time `json:"hire_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Authentication & account control
	PasswordHash string        `json:"-"` // Never exposed
	Status       TeacherStatus `json:"status"`
}
