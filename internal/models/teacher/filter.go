package teacherModels

import (
	"time"

	"github.com/student-management/internal/models"
)

type FilterTeachers struct {
	Name          string        `json:"full_name"`
	Gender        models.Gender `json:"gender"`
	Email         string        `json:"email"`
	EmployeeID    string        `json:"employee_id"`
	Status        TeacherStatus `json:"status"`
	ClassAssigned []string      `json:"class_assigned"`
	SubjectTaught []string      `json:"subject_taught"`
	HireDateFrom  *time.Time    `json:"hire_date_from"`
	HireDateTo    *time.Time    `json:"hire_date_to"`
}
