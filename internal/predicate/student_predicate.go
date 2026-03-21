package predicate

import "github.com/student-management/internal/models"

type PredicateStudent func(*models.Student) bool

func And(predicates ...PredicateStudent) PredicateStudent {
	return func(s *models.Student) bool {
		for _, v := range predicates {
			if !v(s) {
				return false
			}
		}
		return true
	}
}
