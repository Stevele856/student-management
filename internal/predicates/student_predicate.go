package predicates

import "github.com/student-management/internal/models/student"


type PredicateStudent func(*studentModels.Student) bool

func AndStudent(predicates ...PredicateStudent) PredicateStudent {
	return func(s *studentModels.Student) bool {
		for _, p := range predicates {
			if !p(s) {
				return false
			}
		}
		return true
	}
}
