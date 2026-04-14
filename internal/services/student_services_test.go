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
			name:        "valid student",
			student:     baseStudent(),
			expectedErr: nil,
		},

		{
			name: "empty fullname",
			student: func() *models.Student {
				s := baseStudent()
				s.FullName = ""
				return s
			}(),
			expectedErr: ErrNameRequired,
		},

		{
			name: "invalid name format",
			student: func() *models.Student {
				s := baseStudent()
				s.FullName = "Ng@e6n Va*n 4"
				return s
			}(),
			expectedErr: ErrNameFormat,
		},

		{
			name: "invalid email format",
			student: func() *models.Student {
				s := baseStudent()
				s.Email = "invalid-email"
				return s
			}(),
			expectedErr: ErrEmailFormat,
		},

		{
			name: "empty email",
			student: func() *models.Student {
				s := baseStudent()
				s.Email = ""
				return s
			}(),
			expectedErr: ErrEmailRequired,
		},

		{
			name: "future date of birth",
			student: func() *models.Student {
				s := baseStudent()
				s.DateOfBirth = now.AddDate(1, 1, 1)
				return s
			}(),
			expectedErr: ErrValidDOB,
		},

		{
			name: "class empty",
			student: func() *models.Student {
				s := baseStudent()
				s.Class = ""
				return s
			}(),
			expectedErr: ErrStudentClass,
		},

		{
			name: "error class format",
			student: func() *models.Student {
				s := baseStudent()
				s.Class = "10*A$2"
				return s
			}(),
			expectedErr: ErrClassFormat,
		},

		{
			name: "score invalid",
			student: func() *models.Student {
				s := baseStudent()
				s.Scores = []*models.SubjectScore{{Subject: "Toan", Score: 11}}
				return s
			}(),
			expectedErr: ErrScore,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateStudent(tc.student)
			assertError(t, err, tc.expectedErr)
		})
	}
}

func TestNormalizeStudent(t *testing.T) {
	expectedStudent := baseStudent()

	tests := []struct {
		name                string
		input               *models.Student
		expectedStudentData *models.Student
	}{
		{
			name: "trim space in full name",
			input: func() *models.Student {
				s := baseStudent()
				s.FullName = "  Nguyen Van A  "
				return s
			}(),
			expectedStudentData: expectedStudent,
		},
		{
			name: "lowercase and trim email",
			input: func() *models.Student {
				s := baseStudent()
				s.Email = "  VuLt@gmail.com "
				return s
			}(),
			expectedStudentData: expectedStudent,
		},
		{
			name: "trim class",
			input: func() *models.Student {
				s := baseStudent()
				s.Class = " 4A "
				return s
			}(),
			expectedStudentData: expectedStudent,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// SAVE ORIGINAL
			originalFullname := tc.input.FullName
			orginalEmail := tc.input.Email
			originalClass := tc.input.Class

			result := normalizeStudent(tc.input)

			assertField(t, "Fullname", result.FullName, tc.expectedStudentData.FullName)
			assertField(t, "Email", result.Email, tc.expectedStudentData.Email)
			assertField(t, "Fullname", result.Class, tc.expectedStudentData.Class)

			// make sure input does not being mutated
			mutateField(t, "Fullname", tc.input.FullName, originalFullname)
			mutateField(t, "Email", tc.input.Email, orginalEmail)
			mutateField(t, "Email", tc.input.Class, originalClass)

			// Check return new pointer
			if result == tc.input {
				t.Error("normalizeStudent() return same pointer with input")
			}
		})
	}

}

func TestAddStudent(t *testing.T) {
	validStudent := makeValidStudent("")

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
				GetStudentByEmailFn: returnGetStudentByEmail(&models.Student{Email: validStudent.Email}),
			},
			expectedErr: ErrEmailExisted,
		},
		{
			name:  "success create student",
			input: validStudent,
			mockRepo: &MockStudentRepository{
				GetStudentByEmailFn: returnDBError("not found"),
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
			service := &StudentService{repo: tc.mockRepo}
			err := service.AddStudent(tc.input)

			assertError(t, err, tc.expectedErr)
		})
	}
}

func TestUpdateStudent(t *testing.T) {

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
			name:        "missing ID",
			input:       makeValidStudent(""),
			mockRepo:    &MockStudentRepository{},
			expectedErr: ErrIDRequired,
		},
		{
			name:  "error GetStudentByID",
			input: makeValidStudent("id-123"),
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn: returnGetStudentByID(nil),
			},
			expectedErr: ErrStudentNotFound,
		},
		{
			name:  "student ID not found",
			input: makeValidStudent("id-123"),
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn: returnGetStudentByID(nil),
			},
			expectedErr: ErrStudentNotFound,
		},
		{
			name:  "email already existed",
			input: makeValidStudent("id-123"),
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn:    returnGetStudentByID(makeValidStudent("id-123")),
				GetStudentByEmailFn: returnGetStudentByEmail(&models.Student{ID: "other-id"}),
			},
			expectedErr: ErrEmailExisted,
		},
		{
			name:  "email existed and same ID",
			input: makeValidStudent("id-123"),
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn:    returnGetStudentByID(makeValidStudent("id-123")),
				GetStudentByEmailFn: returnGetStudentByEmail(&models.Student{ID: "id-123"}),
				UpdateStudentFn:     func(student *models.Student) error { return nil },
			},
			expectedErr: nil,
		},

		{
			name:  "update student successfully", // HAPPY PATH
			input: makeValidStudent("id-123"),
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn:    returnGetStudentByID(makeValidStudent("id-123")),
				GetStudentByEmailFn: returnGetStudentByEmail(nil),
				UpdateStudentFn:     func(student *models.Student) error { return nil },
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &StudentService{repo: tc.mockRepo}
			err := service.UpdateStudent(tc.input)

			assertError(t, err, tc.expectedErr)
		})
	}
}

func TestDeleteStudent(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		mockRepo    *MockStudentRepository
		expectedErr error
	}{
		{
			name:        "WhiteSpace ID",
			input:       "   ",
			mockRepo:    &MockStudentRepository{},
			expectedErr: ErrIDRequired,
		},
		{
			name:        "missing ID",
			input:       "",
			mockRepo:    &MockStudentRepository{},
			expectedErr: ErrIDRequired,
		},
		{
			name:  "error GetStudentByID",
			input: "id-123",
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn: returnGetStudentByID(nil),
			},
			expectedErr: ErrStudentNotFound,
		},

		{
			name:  "student ID not found",
			input: "id-123",
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn: returnGetStudentByID(nil),
			},
			expectedErr: ErrStudentNotFound,
		},
		{
			name:  "delete student successfully", // HAPPY PATH
			input: "id-123",
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn: returnGetStudentByID(makeValidStudent("id-123")),
				DeleteStudentFn:  func(studentID string) error { return nil },
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &StudentService{repo: tc.mockRepo}
			err := service.DeleteStudent(tc.input)
			assertError(t, err, tc.expectedErr)
		})
	}
}

func TestGetAllStudents(t *testing.T) {
	tests := []struct {
		name             string
		mockRepo         *MockStudentRepository
		expectedStudents []*models.Student
		expectedErr      error
	}{
		{
			name: "repo return error",
			mockRepo: &MockStudentRepository{
				GetAllStudentsFn: func() ([]*models.Student, error) {
					return nil, ErrStudentData
				},
			},
			expectedStudents: nil,
			expectedErr:      ErrStudentData,
		},

		{
			name: "repo empty list",
			mockRepo: &MockStudentRepository{
				GetAllStudentsFn: func() ([]*models.Student, error) {
					return []*models.Student{}, nil
				},
			},
			expectedStudents: []*models.Student{},
			expectedErr:      nil,
		},

		{
			name: "repo return students",
			mockRepo: &MockStudentRepository{
				GetAllStudentsFn: func() ([]*models.Student, error) {
					return []*models.Student{
						makeValidStudent("id-1"),
						makeValidStudent("id-2"),
						makeValidStudent("id-3"),
					}, nil
				},
			},
			expectedStudents: []*models.Student{
				makeValidStudent("id-1"),
				makeValidStudent("id-2"),
				makeValidStudent("id-3"),
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &StudentService{repo: tc.mockRepo}
			students, err := service.GetAllStudents()
			assertError(t, err, tc.expectedErr)

			if len(students) != len(tc.expectedStudents) {
				t.Errorf("len(students): got %d, want %d", len(students), len(tc.expectedStudents))
			}
		})
	}
}

func TestGetStudentByID(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		mockRepo        *MockStudentRepository
		expectedStudent *models.Student
		expectedErr     error
	}{
		{
			name:            "whitespace ID",
			input:           "   ",
			mockRepo:        &MockStudentRepository{},
			expectedStudent: nil,
			expectedErr:     ErrIDRequired,
		},
		{
			name:            "missing ID",
			input:           "",
			mockRepo:        &MockStudentRepository{},
			expectedStudent: nil,
			expectedErr:     ErrIDRequired,
		},

		{
			name:  "error GetStudentByID",
			input: "id-123",
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn: func(studentID string) (*models.Student, error) {
					return nil, ErrStudentData
				},
			},
			expectedStudent: nil,
			expectedErr:     ErrStudentData,
		},

		{
			name:  "get student by ID successfully",
			input: "id-123",
			mockRepo: &MockStudentRepository{
				GetStudentByIDFn: returnGetStudentByID(makeValidStudent("id-123")),
			},
			expectedStudent: makeValidStudent("id-123"),
			expectedErr:     nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &StudentService{repo: tc.mockRepo}
			student, err := service.GetStudentByID(tc.input)
			assertError(t, err, tc.expectedErr)

			if (student == nil) != (tc.expectedStudent == nil) {
				t.Errorf("student: got %v, want %v", student, tc.expectedStudent)
			}
		})
	}
}

func TestGetStudentByEmail(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		mockRepo        *MockStudentRepository
		expectedStudent *models.Student
		expectedErr     error
	}{
		{
			name:            "missing email",
			input:           "",
			mockRepo:        &MockStudentRepository{},
			expectedStudent: nil,
			expectedErr:     ErrEmailRequired,
		},
		{
			name:            "invalid email format",
			input:           "invalid-email",
			mockRepo:        &MockStudentRepository{},
			expectedStudent: nil,
			expectedErr:     ErrEmailFormat,
		},
		{
			name:  "white space email",
			input: "  vult@gmail.com  ",
			mockRepo: &MockStudentRepository{
				GetStudentByEmailFn: returnGetStudentByEmail(makeValidStudent("id-123")),
			},
			expectedStudent: makeValidStudent("id-123"),
			expectedErr:     nil,
		},
		{
			name:  "error return repo",
			input: "vult@gmail.com",
			mockRepo: &MockStudentRepository{
				GetStudentByEmailFn: returnDBError("DB error"),
			},
			expectedStudent: nil,
			expectedErr:     errors.New("DB error"),
		},
		{
			name:  "get student by email return successfully",
			input: "vult@gmail.com",
			mockRepo: &MockStudentRepository{
				GetStudentByEmailFn: returnGetStudentByEmail(makeValidStudent("id-123")),
			},
			expectedStudent: makeValidStudent("id-123"),
			expectedErr:     nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &StudentService{repo: tc.mockRepo}
			student, err := service.GetStudentByEmail(tc.input)

			assertError(t, err, tc.expectedErr)
			if (student == nil) != (tc.expectedStudent == nil) {
				t.Errorf("student: got %v, want %v", student, tc.expectedStudent)
			}
		})
	}
}
