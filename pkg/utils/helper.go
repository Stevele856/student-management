// Utils → format validation
package utils

import (
	"net/mail"
	"regexp"
	"strings"

	studentModels "github.com/student-management/internal/models/student"
)

func IsValidEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}

	_, err := mail.ParseAddress(email)
	return err == nil
}

var validName = regexp.MustCompile(`^[\p{L}]+(?:[\s'-][\p{L}]+)*$`)

func IsValidName(name string) bool {

	if len(name) < 3 || len(name) > 50 {
		return false
	}

	return validName.MatchString(name)
}

var studentClass = regexp.MustCompile(`^\d{1,2}-?[A-Z]$`)
func IsValidClass(class string) bool {
	return studentClass.MatchString(class)
}

func IsValidScores(scores []*studentModels.SubjectScore) bool {
	// MAXIMUM 10 SUBJECTS
	if len(scores) > 10 {
		return false
	}

	for _, s := range scores {
		if !IsValidSubjectScore(s.Score) {
			return false
		}
	}
	return true
}

func IsValidSubjectScore(score float64) bool {
	return score >= 0 && score <= 10
}

var studentSubject = regexp.MustCompile(`^[\p{L}\s]{3,30}$`)
func IsValidSubject(subject string) bool {
	return studentSubject.MatchString(subject)

}

func CalAvgScore(scores []*studentModels.SubjectScore) float64 {
	if len(scores) == 0 {
		return 0
	}

	var total float64
	for _, s := range scores {
		total += s.Score
	}

	return total / float64(len(scores))
}

func CalcStudentRankBaseOnAvgScore(avg float64) studentModels.Rank {
	switch {
	case avg >= 9.0:
		return studentModels.Excellent
	case avg >= 7.0:
		return studentModels.Good
	case avg >= 5.0:
		return studentModels.Average
	default:
		return studentModels.Weak
	}
}

var validEmployeeIDPattern = regexp.MustCompile(`^T\d{3}$`)

func IsValidEmployeeID(employeeID string) bool {
	employeeID = strings.TrimSpace(employeeID)
	return validEmployeeIDPattern.MatchString(employeeID)
}

var validPhoneNumber = regexp.MustCompile(`^(0|\+84)(3|5|7|8|9)[0-9]{8}$`)

func IsValidPhoneNumber(phone string) bool {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, ".", "")
	return validPhoneNumber.MatchString(phone)
}


