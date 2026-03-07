// Utils → format validation
package util

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/student-management/internal/models"
)

func IsValidStudentEmail(email string) bool {
	email = strings.TrimSpace(email)

	if len(email) == 0 || len(email) > 254 {
		return false
	}

	_, err := mail.ParseAddress(email)
	return err == nil
}



var studentName = regexp.MustCompile(`^[\p{L}]+(?:[\s'-][\p{L}]+)*$`)
func IsValidStudentName(name string) bool {
	name = strings.TrimSpace(name)

	if len(name) < 3 || len(name) > 50 {
		return false
	}

	return studentName.MatchString(name)
}



var studentClass = regexp.MustCompile(`^\d{1,2}-?[A-Z]$`)
func IsValidClass(class string) bool {
	class = strings.TrimSpace(class)

	return studentClass.MatchString(class)
}



func IsValidScores(scores []*models.SubjectScore) bool {
	// MAXIMUM 10 SUBJECTS
	if len(scores) > 10 {
		return false
	}

	for _, s := range scores {
		if !IsValidSubjectScore(s.Score){
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
	subject = strings.TrimSpace(subject)

	return studentSubject.MatchString(subject)

}
