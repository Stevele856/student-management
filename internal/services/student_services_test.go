package services

import (
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
