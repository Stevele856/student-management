package models

import (
	"fmt"
	"strings"
)


type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

func ParseGender(value string) (Gender, error) {
	g := Gender(strings.ToLower(strings.TrimSpace(value)))

	switch g {
	case GenderMale, GenderFemale:
		return g, nil
	default:
		return "", fmt.Errorf("invalid gender: %q (must be 'male' or 'female')", value)
	}
}
