// Step 1

package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/student-management/internal/models"
)

type Rank string

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
	Gender      models.Gender   `json:"gender"`
	Address     string          `json:"address"`
	Class       string          `json:"class"`
	Email       string          `json:"email"`
	Scores      []*SubjectScore `json:"scores"`
}

// Import and Export CSV
const DateLayout = "2006-01-02"

// Return the first row of CSV file
func CSVHeader() []string {
	return []string{"id", "name", "dob", "gender", "address", "class", "email", "subject", "score"}
}

// return [][]string - slice of row - one student with 3 scores produces 3 csv rows
func (s *Student) ToCSVRows() [][]string {
	base := []string{
		s.ID,
		s.FullName,
		s.DateOfBirth.Format(DateLayout),
		string(s.Gender),
		s.Address,
		s.Class,
		s.Email,
	}

	// edge case, if the student has no score, still write one row with empty subject and score column
	if len(s.Scores) == 0 {
		return [][]string{append(base, "", "")}
	}

	rows := make([][]string, 0, len(s.Scores))
	for _, sc := range s.Scores {
		row := make([]string, len(base))
		copy(row, base)
		row = append(row, sc.Subject, strconv.FormatFloat(sc.Score, 'f', 2, 64))
		rows = append(rows, row)
	}

	return rows
}

// StudentFromCSVRows build a Student from a group of rows sharing the same ID
func StudentFromCSVRows(rows [][]string) (*Student, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows provided")
	}

	first := rows[0]
	if len(first) < 9 {
		return nil, fmt.Errorf("invalid row length: expected 9 columns, got %d", len(first))
	}

	dob, err := time.Parse(DateLayout, first[2])
	if err != nil {
		return nil, fmt.Errorf("invalid date_of_birth %q: %w", first[2], err)
	}

	gender, err := models.ParseGender(first[3])
	if err != nil {
		return nil, fmt.Errorf("invalid gender at row: %w", err)
	}

	s := &Student{
		ID:          first[0],
		FullName:    first[1],
		DateOfBirth: dob,
		Gender:      gender,
		Address:     first[4],
		Class:       first[5],
		Email:       first[6],
		Scores:      make([]*SubjectScore, 0),
	}

	for _, row := range rows {
		subject := row[7]
		scoreStr := row[8]
		if subject == "" {
			continue // no score for this row
		}
		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid score %q for subject %q: %w", scoreStr, subject, err)
		}
		s.Scores = append(s.Scores, &SubjectScore{Subject: subject, Score: score})
	}

	return s, nil
}


