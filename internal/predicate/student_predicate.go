package predicate

import "github.com/student-management/internal/models/student"


type PredicateStudent func(*studentModels.Student) bool

func And(predicates ...PredicateStudent) PredicateStudent {
	return func(s *studentModels.Student) bool {
		for _, v := range predicates {
			if !v(s) {
				return false
			}
		}
		return true
	}
}
