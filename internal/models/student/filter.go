package studentModels

import "github.com/student-management/internal/models"

type FilterStudents struct {
	Name        string  `json:"full_name"`
	Class       string  `json:"class"`
	YearOfBirth int     `json:"year_of_birth"`
	Gender      models.Gender  `json:"gender"`
	Address     string  `json:"address"`
	MinAvgScore float64 `json:"min_avg_score"`
	MaxAvgScore float64 `json:"max_avg_score"`
	StudentRank Rank    `json:"student_rank"`
}
