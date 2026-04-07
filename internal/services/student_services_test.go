package services

import (
	"errors"
	"testing"
	"time"

	"github.com/student-management/internal/models"
)

func TestValidateStudent(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		student     *models.Student
		expectedErr error
	}{
		{
			name: "valid student",
			student: &models.Student{
				FullName:    "Nguyen Van A",
				DateOfBirth: now.AddDate(-1, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "4A",
				Email:       "vult@gmail.com",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedErr: nil,
		},

		{
			name: "empty fullname",
			student: &models.Student{
				FullName:    "",
				DateOfBirth: now.AddDate(-1, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "4A",
				Email:       "vult@gmail.com",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedErr: ErrNameRequired,
		},

		{
			name: "invalid name format",
			student: &models.Student{
				FullName:    "Nguy4n V4n A",
				DateOfBirth: now.AddDate(-1, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "4A",
				Email:       "vult@gmail.com",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedErr: ErrNameFormat,
		},

		{
			name: "invalid email format",
			student: &models.Student{
				FullName:    "Nguyen Van A",
				DateOfBirth: now.AddDate(-1, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "4A",
				Email:       "invalid-email",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedErr: ErrEmailFormat,
		},

		{
			name: "future date of birth",
			student: &models.Student{
				FullName:    "Nguyen Van A",
				DateOfBirth: now.AddDate(1, 0, 0),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "4A",
				Email:       "vult@gmail.com",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedErr: ErrValidDOB,
		},

		{
			name: "class empty",
			student: &models.Student{
				FullName:    "Nguyen Van A",
				DateOfBirth: now.AddDate(-1, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "",
				Email:       "vult@gmail.com",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedErr: ErrStudentClass,
		},

		{
			name: "err class format",
			student: &models.Student{
				FullName:    "Nguyen Van A",
				DateOfBirth: now.AddDate(-1, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "10_A3",
				Email:       "vult@gmail.com",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedErr: ErrClassFormat,
		},

		{
			name: "score invalid",
			student: &models.Student{
				FullName:    "Nguyen Van A",
				DateOfBirth: now.AddDate(-1, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "4A",
				Email:       "vult@gmail.com",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   11,
					},
				},
			},
			expectedErr: ErrScore,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateStudent(tc.student)
			if err != tc.expectedErr {
				t.Errorf("[%s] got %v, want %v", tc.name, err, tc.expectedErr)
			}
		})
	}
}

func TestNormalizeStudent(t *testing.T) {
	expectedStudent := &models.Student{
		FullName:    "Nguyen Van A",
		DateOfBirth: time.Now().AddDate(-26, 2, 3),
		Gender:      "male",
		Address:     "Ho Chi Minh City",
		Class:       "4A",
		Email:       "vult@gmail.com",
		Scores: []*models.SubjectScore{
			{
				Subject: "Toan",
				Score:   6.5,
			},
		},
	}

	tests := []struct {
		name                string
		input               *models.Student
		expectedStudentData *models.Student
	}{
		{
			name: "trim space in full name",
			input: &models.Student{
				FullName:    "   Nguyen Van A ",
				DateOfBirth: time.Now().AddDate(-26, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "4A",
				Email:       "vult@gmail.com",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedStudentData: expectedStudent,
		},
		{
			name: "lowercase and trim email",
			input: &models.Student{
				FullName:    "Nguyen Van A",
				DateOfBirth: time.Now().AddDate(-26, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       "4A",
				Email:       " vUlt@GmaiL.com ",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedStudentData: expectedStudent,
		},
		{
			name: "trim class",
			input: &models.Student{
				FullName:    "Nguyen Van A",
				DateOfBirth: time.Now().AddDate(-26, 2, 3),
				Gender:      "male",
				Address:     "Ho Chi Minh City",
				Class:       " 4A ",
				Email:       " vult@gmaiL.com ",
				Scores: []*models.SubjectScore{
					{
						Subject: "Toan",
						Score:   6.5,
					},
				},
			},
			expectedStudentData: expectedStudent,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			originalFullname := tc.input.FullName
			orginalEmail := tc.input.Email
			originalClass := tc.input.Class

			result := normalizeStudent(tc.input)

			if result.FullName != tc.expectedStudentData.FullName {
				t.Errorf("Fullname: got %q, want %q", result.FullName, tc.expectedStudentData.FullName)
			}

			if result.Email != tc.expectedStudentData.Email {
				t.Errorf("Email: got %q, want %q", result.Email, tc.expectedStudentData.Email)
			}

			if result.Class != tc.expectedStudentData.Class {
				t.Errorf("Class: got %q, want %q", result.Class, tc.expectedStudentData.Class)
			}

			// make sure input does not being mutated
			if tc.input.FullName != originalFullname {
				t.Errorf("input.FullName bị mutate: got %q, want %q", tc.input.FullName, tc.input.FullName)
			}
			if tc.input.Email != orginalEmail {
				t.Errorf("input.Email bị mutate: got %q, want %q", tc.input.Email, tc.expectedStudentData.Email)
			}
			if tc.input.Class != originalClass {
				t.Errorf("input.Class bị mutate: got %q, want %q", tc.input.Class, tc.expectedStudentData.Class)
			}

			// Kiểm tra trả về pointer mới
			if result == tc.input {
				t.Error("normalizeStudent() trả về cùng pointer với input")
			}
		})
	}

}

func TestAddStudent(t *testing.T) {
	validStudent := &models.Student{
		FullName:    "Nguyen Van A",
		DateOfBirth: time.Now().AddDate(-26, 2, 3),
		Gender:      "male",
		Address:     "Ho Chi Minh City",
		Class:       "4A",
		Email:       "vult@gmail.com",
		Scores: []*models.SubjectScore{
			{
				Subject: "Toan",
				Score:   6.5,
			},
		},
	}

	tests := []struct {
		name        string
		input       *models.Student
		mockRepo    *MockStudentRepository
		expectedErr error
	}{
		{
			name:        "nil student",
			input:       nil,
			mockRepo:    &MockStudentRepository{},
			expectedErr: ErrStudentData,
		},
		{
			name:        "validation failed",
			input:       &models.Student{},
			mockRepo:    &MockStudentRepository{},
			expectedErr: ErrNameRequired,
		},
		{
			name:  "email already existed",
			input: validStudent,
			mockRepo: &MockStudentRepository{
				GetStudentByEmailFn: func(email string) (*models.Student, error) {
					return &models.Student{Email: email}, nil
				},
			},
			expectedErr: ErrEmailExisted,
		},
		{
			name:  "success create student",
			input: validStudent,
			mockRepo: &MockStudentRepository{
				GetStudentByEmailFn: func(email string) (*models.Student, error) {
					return nil, errors.New("not found")
				},
				AddStudentFn: func(student *models.Student) error {
					if student.ID == "" {
						t.Errorf("expected ID to be generated")
					}
					return nil
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &StudentService{
				repo: tc.mockRepo,
			}

			err := service.AddStudent(tc.input)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("[%s] got %v, want %v", tc.name, err, tc.expectedErr)
			}
		})
	}
}
