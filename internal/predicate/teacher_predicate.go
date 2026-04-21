package predicate

import "github.com/student-management/internal/models/teacher"

type PredicateTeacher func(*teacherModels.Teacher) bool

func AndTeacher(predicates ...PredicateTeacher) PredicateTeacher{
	return func(s *teacherModels.Teacher) bool {
		for _, predicate := range predicates{
			if !predicate(s){
				return false
			}
		}
		return true
	}
}