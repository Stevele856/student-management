package services

import (
	"errors"
	"testing"
	"time"

	"github.com/student-management/internal/models/student"
)

func baseStudent() *models.Student {
	return &models.Student{
		FullName:    "Nguyen Van A",
		DateOfBirth: time.Now().AddDate(-26, 2, 3),
		Gender:      "male",
		Address:     "Ho Chi Minh City",
		Class:       "4A",
		Email:       "vult@gmail.com",
		Scores:      []*models.SubjectScore{{Subject: "Toan", Score: 6.5}},
	}
}

func makeValidStudent(id string) *models.Student{
	s := baseStudent()
	s.ID = id
	return s
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if want != nil {
		if got == nil {
			t.Errorf("expected error %q, got nil", want)
		} else if got.Error() != want.Error() {
			t.Errorf("expected error %q, got %q", want, got)
		}
	} else {
		if got != nil {
			t.Errorf("expected no error, got %q", got)
		}
	}
}

func assertField(t *testing.T, field, got, want string){
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", field, got, want)
	}
}

func mutateField(t *testing.T, field, got, original string){
	t.Helper()
	if got != original{
		t.Errorf("%s, got %q, want %q", field, got, original)
	}
}

// mock func, closure func
func returnGetStudentByEmail(s *models.Student) func(string) (*models.Student, error){
	return func(email string) (*models.Student, error) {
		return s, nil
	}
}

func returnGetStudentByID(s *models.Student) func(string) (*models.Student, error){
	return func(id string) (*models.Student, error) {
		return s, nil
	}
}

func returnDBError(msg string) func(string) (*models.Student, error) {
	return func(_ string) (*models.Student, error) {
		return nil, errors.New(msg)
	}
}