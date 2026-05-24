package teacherModels

import (
	"fmt"
	"strings"
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

// Import and Export CSV
const DateLayout = "2006-01-02"

// CSVHeader returns the first row of teacher CSV file.
func CSVHeader() []string {
	return []string{
		"id",
		"name",
		"email",
		"dob",
		"gender",
		"address",
		"phone",
		"employee_id",
		"subject",
		"class",
		"hire_date",
		"status",
	}
}

// ToCSVRows converts one teacher into CSV rows.
// A teacher with multiple subjects/classes becomes multiple rows.
func (t *Teacher) ToCSVRows() [][]string {
	subjects := t.SubjectTaught
	classes := t.ClassAssigned

	if len(subjects) == 0 {
		subjects = []string{""}
	}
	if len(classes) == 0 {
		classes = []string{""}
	}

	base := []string{
		t.ID,
		t.FullName,
		t.Email,
		t.DateOfBirth.Format(DateLayout),
		string(t.Gender),
		t.Address,
		t.Phone,
		t.EmployeeID,
	}

	limit := len(subjects)
	if len(classes) > limit {
		limit = len(classes)
	}
	rows := make([][]string, 0, limit)
	for i := 0; i < limit; i++ {
		subject := ""
		class := ""
		if i < len(subjects) {
			subject = subjects[i]
		}
		if i < len(classes) {
			class = classes[i]
		}

		row := make([]string, len(base))
		copy(row, base)
		row = append(row, subject, class, t.HireDate.Format(DateLayout), string(t.Status))
		rows = append(rows, row)
	}

	return rows
}

// TeacherFromCSVRows builds a Teacher from a group of rows sharing the same ID.
func TeacherFromCSVRows(rows [][]string) (*Teacher, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows provided")
	}

	first := rows[0]
	if len(first) < 12 {
		return nil, fmt.Errorf("invalid row length: expected 12 columns, got %d", len(first))
	}

	dob, err := time.Parse(DateLayout, first[3])
	if err != nil {
		return nil, fmt.Errorf("invalid date_of_birth %q: %w", first[3], err)
	}

	gender, err := models.ParseGender(first[4])
	if err != nil {
		return nil, fmt.Errorf("invalid gender: %w", err)
	}

	hireDate, err := time.Parse(DateLayout, first[10])
	if err != nil {
		return nil, fmt.Errorf("invalid hire_date %q: %w", first[10], err)
	}

	status := TeacherStatus(strings.ToLower(strings.TrimSpace(first[11])))

	t := &Teacher{
		ID:            first[0],
		FullName:      first[1],
		Email:         first[2],
		DateOfBirth:   dob,
		Gender:        gender,
		Address:       first[5],
		Phone:         first[6],
		EmployeeID:    first[7],
		SubjectTaught: make([]string, 0),
		ClassAssigned: make([]string, 0),
		HireDate:      hireDate,
		Status:        status,
	}

	subjectSet := make(map[string]struct{})
	classSet := make(map[string]struct{})

	for _, row := range rows {
		if len(row) < 12 {
			return nil, fmt.Errorf("invalid row length: expected 12 columns, got %d", len(row))
		}

		subject := strings.TrimSpace(row[8])
		if subject != "" {
			key := strings.ToLower(subject)
			if _, existed := subjectSet[key]; !existed {
				subjectSet[key] = struct{}{}
				t.SubjectTaught = append(t.SubjectTaught, subject)
			}
		}

		class := strings.TrimSpace(row[9])
		if class != "" {
			key := strings.ToUpper(class)
			if _, existed := classSet[key]; !existed {
				classSet[key] = struct{}{}
				t.ClassAssigned = append(t.ClassAssigned, class)
			}
		}
	}

	return t, nil
}
