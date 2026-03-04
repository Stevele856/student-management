package util

import (
	"net/mail"
	"regexp"
	"strings"
)

func IsValidStudentEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var studentName = regexp.MustCompile(`^[\p{L}]+(?:[\s'-][\p{L}]+)*$`)

func IsValidStudentName(name string) bool{
	name = strings.TrimSpace(name)

	if len(name) < 3 || len(name) > 50 {
		return false
	}

	return studentName.MatchString(name)
}