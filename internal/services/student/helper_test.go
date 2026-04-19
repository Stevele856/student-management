package student

import (
	"errors"
	"testing"
	"time"

	studentModels "github.com/student-management/internal/models/student"
)

func baseStudent() *studentModels.Student {
	return &studentModels.Student{
		FullName:    "Nguyen Van A",
		DateOfBirth: time.Now().AddDate(-26, 2, 3),
		Gender:      "male",
		Address:     "Ho Chi Minh City",
		Class:       "4A",
		Email:       "vult@gmail.com",
		Scores:      []*studentModels.SubjectScore{{Subject: "Toan", Score: 6.5}},
	}
}

func makeValidStudent(id string) *studentModels.Student{
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
func returnGetStudentByEmail(s *studentModels.Student) func(string) (*studentModels.Student, error){
	return func(email string) (*studentModels.Student, error) {
		return s, nil
	}
}

func returnGetStudentByID(s *studentModels.Student) func(string) (*studentModels.Student, error){
	return func(id string) (*studentModels.Student, error) {
		return s, nil
	}
}

func returnDBError(msg string) func(string) (*studentModels.Student, error) {
	return func(_ string) (*studentModels.Student, error) {
		return nil, errors.New(msg)
	}
}