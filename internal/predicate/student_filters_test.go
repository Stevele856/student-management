package predicate_test

import (
	"testing"
	"time"

	"github.com/student-management/internal/models"
	"github.com/student-management/internal/predicate"
)

func mockStudent(name, class, gender string, year int) *models.Student {
	return &models.Student{
		FullName:    name,
		Class:       class,
		Gender:      gender,
		DateOfBirth: time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func mockStudentWithScore(score []*models.SubjectScore) *models.Student {
	return &models.Student{
		Scores: score,
	}
}

func mockStudentWithMultipleCategories(name, class, gender string, year int, score []*models.SubjectScore) *models.Student {
	return &models.Student{
		FullName:    name,
		Class:       class,
		Gender:      gender,
		DateOfBirth: time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Scores:      score,
	}
}

func TestByName(t *testing.T) {
	p := predicate.ByName("Nguyen")

	// Match
	student1 := mockStudent("Nguyen Van A", "10A", "male", 2000)
	if !p(student1) {
		t.Error("expected true, got false")
	}

	// Not match
	student2 := mockStudent("Tran Van B", "11B", "male", 2000)
	if p(student2) {
		t.Error("expected false, got true")
	}
}

func TestByYear(t *testing.T) {
	p := predicate.ByYear(2000)

	// Match
	student1 := mockStudent("Nguyen Van A", "10A", "male", 2000)
	if !p(student1) {
		t.Error("expedted true, got false")
	}

	// Not match
	student2 := mockStudent("Nguyen Van B", "10A", "male", 1999)
	if p(student2) {
		t.Error("expected false, got true")
	}
}

func TestByAvgScore(t *testing.T) {

	// avg = (8.5 + 8 + 9) / 3 = 8.5
	scores := []*models.SubjectScore{
		{Subject: "Toan", Score: 8.5},
		{Subject: "Anh Van", Score: 8.0},
		{Subject: "Tin Hoc", Score: 9.0},
	}

	student := mockStudentWithScore(scores)

	p := predicate.ByAvgScore(8.0, 9.0) // => true

	if !p(student) {
		t.Error("expected true, got false")
	}

	p1 := predicate.ByAvgScore(6.0, 8.0)

	if p1(student) {
		t.Error("expected false, got true")
	}
}

func TestByRank(t *testing.T) {
	// avg = (8.5 + 8 + 9) / 3 = 8.5
	scores := []*models.SubjectScore{
		{Subject: "Toan", Score: 8.5},
		{Subject: "Anh Van", Score: 8.0},
		{Subject: "Tin Hoc", Score: 9.0},
	}
	student := mockStudentWithScore(scores)

	p := predicate.ByRank(models.Good)

	if !p(student) {
		t.Error("expected true, got false")
	}
}

func TestAnd(t *testing.T) {
	p := predicate.And(
		predicate.ByName("Nguyen"),
		predicate.ByClass("10A"),
		predicate.ByAvgScore(7.0, 8.5),
		predicate.ByRank(models.Good),
	)

	// match
	studentA := mockStudentWithMultipleCategories(
		"Nguyen Van A", "10A", "male", 2000,
		[]*models.SubjectScore{
			{Subject: "Toan", Score: 7.0},
			{Subject: "Anh Van", Score: 8.0},
			{Subject: "Toan", Score: 8.0},
		},
	)

	if !p(studentA) {
		t.Error("expected true, got false")
	}

	// not match
	studentB := mockStudentWithMultipleCategories(
		"Nguyen Van A", "10B", "male", 2000,
		[]*models.SubjectScore{
			{Subject: "Toan", Score: 8.0},
			{Subject: "Anh Van", Score: 8.0},
			{Subject: "Toan", Score: 8.0},
		},
	)

	if p(studentB) {
		t.Error("expected false, got true")
	}
}
