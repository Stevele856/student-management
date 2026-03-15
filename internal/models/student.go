// Step 1

package models

import "time"

// Enum type
type Gender string
type Rank string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

const (
	Excellent Rank = "excellent"
	Good      Rank = "good"
	Average   Rank = "average"
	Weak      Rank = "weak"
)

type SubjectScore struct {
	Subject string  `json:"subject"`
	Score   float64 `json:"score"`
}

type Student struct {
	ID          string          `json:"id"`
	FullName    string          `json:"full_name"`
	DateOfBirth time.Time       `json:"date_of_birth"`
	Gender      string          `json:"gender"`
	Address     string          `json:"address"`
	Class       string          `json:"class"`
	Email       string          `json:"email"`
	Scores      []*SubjectScore `json:"scores"`
}
