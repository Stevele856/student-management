package predicate

import (
	"strings"

	"github.com/student-management/internal/models/student"
	"github.com/student-management/pkg/utils"
)

func ByName(name string) PredicateStudent {
	return func(s *models.Student) bool {
		return strings.Contains(
			strings.ToLower(s.FullName),
			strings.ToLower(name),
		)
	}
}

func ByClass(class string) PredicateStudent {
	return func(s *models.Student) bool {
		return strings.EqualFold(s.Class, class)
	}
}

func ByYear(year int) PredicateStudent {
	return func(s *models.Student) bool {
		return s.DateOfBirth.Year() == year
	}
}

func ByGender(gender string) PredicateStudent {
	return func(s *models.Student) bool {
		return strings.EqualFold(string(s.Gender), gender)
	}
}

func ByAddress(address string) PredicateStudent {
	return func(s *models.Student) bool {
		return strings.EqualFold(s.Address, address)
	}
}

func ByAvgScore(min, max float64) PredicateStudent {
	return func(s *models.Student) bool {
		avgScore := utils.CalAvgScore(s.Scores)

		if min > 0 && avgScore < min {
			return false
		}
		if max > 0 && avgScore > max {
			return false
		}
		return true
	}
}

func ByRank(rank models.Rank) PredicateStudent {
	return func(s *models.Student) bool {
		avgScore := utils.CalAvgScore(s.Scores)

		return utils.CalcStudentRankBaseOnAvgScore(avgScore) == rank
	}
}

/*
	-	Tìm học sinh có điểm TB > 8 => minScore = 8 => Điểm TB dưới 8 => loại
	-> f.minScore > 0 && avg < f.minScore{
			return false
		}
			Tìm học sinh có điểm TB < 5 => maxScore = 5 => Điểm TB lớn hơn 5 => loại
	-> f.maxScore > 0 && avg > f.maxScore{
			return false
		}
*/
