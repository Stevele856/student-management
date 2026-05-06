package predicate

import "github.com/student-management/internal/models/teacher"

type PredicateTeacher func(*teacherModels.Teacher) bool

func AndTeacher(predicates ...PredicateTeacher) PredicateTeacher{
	return func(s *teacherModels.Teacher) bool {
		for _, p := range predicates{
			if !p(s){
				return false
			}
		}
		return true
	}
}