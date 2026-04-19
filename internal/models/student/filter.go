package models


type FilterStudents struct {
	Name        string
	Class       string
	YearOfBirth int
	Gender      string
	Address     string
	MinAvgScore float64
	MaxAvgScore float64
	StudentRank Rank
}
