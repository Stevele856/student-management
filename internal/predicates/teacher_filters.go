package predicates

import (
	"strings"
	"time"

	"github.com/student-management/internal/models"
	teacherModels "github.com/student-management/internal/models/teacher"
)

func ByTeacherName(name string) PredicateTeacher {
	return func(t *teacherModels.Teacher) bool {
		return strings.Contains(
			strings.ToLower(t.FullName),
			strings.ToLower(name),
		)
	}
}

func ByTeacherGender(gender models.Gender) PredicateTeacher{
	g := strings.TrimSpace(string(gender))
	return func(t *teacherModels.Teacher) bool {
		return strings.EqualFold(string(t.Gender), g)
	}
}

func ByTeacherEmail(email string) PredicateTeacher{
	return func(t *teacherModels.Teacher) bool {
		return strings.EqualFold(t.Email, email)
	}
}

func ByEmployeeID(employeeID string) PredicateTeacher{
	return func(t *teacherModels.Teacher) bool {
		return strings.EqualFold(t.EmployeeID, employeeID)
	}
}

func ByStatus(status teacherModels.TeacherStatus) PredicateTeacher{
	s := strings.TrimSpace(string(status))
	return func(t *teacherModels.Teacher) bool {
		return strings.EqualFold(string(t.Status), s)
	}
}

func ByClassAssigned(class string) PredicateTeacher{
	return func(t *teacherModels.Teacher) bool {
		for _, c := range t.ClassAssigned{
			if strings.EqualFold(c, class){
				return true
			}
		}
		return false
	}
}

func BySubjectTaught(subject string) PredicateTeacher{
	return func(t *teacherModels.Teacher) bool {
		for _, s := range t.SubjectTaught{
			if strings.EqualFold(s, subject){
				return true
			}
		}
		return false
	}
}

func ByHireDateRange(from, to *time.Time) PredicateTeacher{
	return func(t *teacherModels.Teacher) bool {
		if from != nil && t.HireDate.Before(*from){
			return false
		}
		if to != nil && t.HireDate.After(*to){
			return false
		}
		return true
	}
}
