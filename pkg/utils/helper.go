// Utils → format validation
package utils

import (
	"net/mail"
	"regexp"
	"github.com/student-management/internal/models"
)

func IsValidStudentEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}

	_, err := mail.ParseAddress(email)
	return err == nil
}



var studentName = regexp.MustCompile(`^[\p{L}]+(?:[\s'-][\p{L}]+)*$`)
func IsValidStudentName(name string) bool {

	if len(name) < 3 || len(name) > 50 {
		return false
	}

	return studentName.MatchString(name)
}



var studentClass = regexp.MustCompile(`^\d{1,2}-?[A-Z]$`)
func IsValidClass(class string) bool {
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
	return studentSubject.MatchString(subject)

}

func CalAvgScore(scores []*models.SubjectScore) float64 {
	if len(scores) == 0 {
		return 0
	}

	var total float64
	for _, s := range scores {
		total += s.Score
	}

	return total / float64(len(scores))
}


func CalcStudentRankBaseOnAvgScore(avg float64) models.Rank {
	switch {
	case avg >= 9.0:
		return models.Excellent
	case avg >= 7.0:
		return models.Good
	case avg >= 5.0:
		return models.Average
	default:
		return models.Weak
	}
}
