package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/student-management/internal/models"
	"github.com/student-management/internal/services"
)

func TestAddStudent(t *testing.T) {
	// HELPER CREATE VALID STUDENT
	validStudent := func() *models.Student {
		return &models.Student{
			FullName:    "Le Trong Vu",
			DateOfBirth: time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
			Gender:      "male",
			Address:     "Ho Chi Minh",
			Class:       "DQT4",
			Email:       "letrongvu.work@gmail.com",
			Scores:      []*models.SubjectScore{},
		}
	}

	// HELPER CREATE DEFAULT MOCK REPO
	happyRepo := func() *MockStudentRepository {
		return &MockStudentRepository{
			GetStudentByEmailFn: func(studentEmail string) (*models.Student, error) {
				return nil, errors.New("email not found")
			},
			AddStudentFn: func(student *models.Student) error {
				return nil
			},
		}
	}
	// HAPPY PATH test cases
	tests := []struct {
		name        string
		input       *models.Student
		setupRepo   func() *MockStudentRepository
		expectedErr error
	}{
		/*------ EACH VALIDATION PASSED -------*/
		{
			name:        "success - valid student",
			input:       validStudent(),
			setupRepo:   happyRepo,
			expectedErr: nil,
		},
		// CASE 1 - EMPTY ID
		{
			name: "success - auto generate UUID when student ID is empty",
			input: func() *models.Student {
				s := validStudent()
				s.ID = ""
				return s
			}(),
			setupRepo:   happyRepo,
			expectedErr: nil,
		},
		// CASE 2: EMAIL UPPERCASE
		{
			name: "success - auto lowercase email",
			input: func() *models.Student {
				s := validStudent()
				s.Email = "LeTrongVu.Work@Gmaik.COM"
				return s
			}(),
			setupRepo:   happyRepo,
			expectedErr: nil,
		},

		/*------ NULL - EMPTY INPUT -------*/
		{
			name:        "fail - student nil",
			input:       nil,
			setupRepo:   happyRepo,
			expectedErr: services.ErrStudentData,
		},
		// CASE 1 - FAILED (EMPTY FULLNAME)
		{
			name: "fail - Empty fullname",
			input: func() *models.Student {
				s := validStudent()
				s.FullName = ""
				return s
			}(),
			setupRepo:   happyRepo,
			expectedErr: services.ErrStudentData,
		},
		// CASE 2 - FULLNAME CONTAINS SPACE
		{
			name: "fail - fullname contain space",
			input: func() *models.Student {
				s := validStudent()
				s.FullName = "  "
				return s
			}(),
			setupRepo:   happyRepo,
			expectedErr: services.ErrStudentData,
		},
		// CASE 3 - EMPTY EMAIL
		{
			name: "fail - empty email",
			input: func() *models.Student {
				s := validStudent()
				s.Email = ""
				return s
			}(),
			setupRepo:   happyRepo,
			expectedErr: services.ErrStudentData,
		},

		/*------ EMAIL VALIDATION -------*/
		{
			name:  "fail - email existed in DB",
			input: validStudent(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByEmailFn: func(studentEmail string) (*models.Student, error) {
						return &models.Student{Email: studentEmail}, nil // <- return student -> email existed in DB
					},
				}
			},
			expectedErr: services.ErrEmailExisted,
		},
		// CASE 1: INCORRECT EMAIL FORMAT
		{
			name: "fail - incorrect email format",
			input: func() *models.Student {
				s := validStudent()
				s.Email = "not-an-email"
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByEmailFn: func(studentEmail string) (*models.Student, error) {
						return nil, errors.New("not found") // <- email does not existed yet, but incorrect format
					},
				}
			},
			expectedErr: services.ErrEmailFormat,
		},
		// CASE 2: EMAIL DOES NOT CONTAIN DOMAIN
		{
			name: "fail - email does not contain domain",
			input: func() *models.Student {
				s := validStudent()
				s.Email = "letrongvu.work@"
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByEmailFn: func(studentEmail string) (*models.Student, error) {
						return nil, errors.New("not found")
					},
				}
			},
			expectedErr: services.ErrEmailFormat,
		},
		// CASE 3: EMAIL DOES NOT CONTAIN '@'
		{
			name: "fail - email does not contain '@' ",
			input: func() *models.Student {
				s := validStudent()
				s.Email = "letrongvu.work.gmail.com"
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByEmailFn: func(studentEmail string) (*models.Student, error) {
						return nil, errors.New("not found")
					},
				}
			},
			expectedErr: services.ErrEmailFormat,
		},

		/*------- ERROR PATH (NAME, DOB, CLASS, SCORE VALIDATION--------*/
		// CASE 1: NAME FORMAT
		{
			name: "fail - name contain special character",
			input: func() *models.Student {
				s := validStudent()
				s.FullName = "Le@Trong`Vu~"
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrNameFormat,
		},
		{
			name: "fail - name contain number",
			input: func() *models.Student {
				s := validStudent()
				s.FullName = "Le Trong 1Vu"
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrNameFormat,
		},
		// CASE 2: DOB IN FUTURE (tomorrow)
		{
			name: "fail - DOB in future",
			input: func() *models.Student {
				s := validStudent()
				s.DateOfBirth = time.Now().Add(24 * time.Hour)
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrValidDOB,
		},
		// CASE 3: DOB IN FUTURE (next hour)
		{
			name: "fail - DOB in today but (next few hour)",
			input: func() *models.Student {
				s := validStudent()
				s.DateOfBirth = time.Now().Add(time.Hour)
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrValidDOB,
		},
		// CASE 4: EMPTY STUDENT CLASS
		{
			name: "fail - empty class",
			input: func() *models.Student {
				s := validStudent()
				s.Class = ""
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrStudentClass,
		},
		// CASE 5: CLASS CONTAIN SPACE
		{
			name: "fail - class contain space",
			input: func() *models.Student {
				s := validStudent()
				s.Class = "   "
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrStudentClass,
		},
		// CASE 6: INVALID CLASS FORMAT
		{
			name: "fail - class format",
			input: func() *models.Student {
				s := validStudent()
				s.Class = "INVALID!"
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrClassFormat,
		},

		/* ------- ERROR SCORE ------- */
		// CASE 1: NEGATIVE SCORE
		{
			name: "fail - negative score",
			input: func() *models.Student {
				s := validStudent()
				s.Scores = []*models.SubjectScore{
					{Subject: "Toan", Score: -5},
				}
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrScore,
		},
		// CASE 2: SCORE GREATER THAN 10
		{
			name: "fail - score greater than 10",
			input: func() *models.Student {
				s := validStudent()
				s.Scores = []*models.SubjectScore{
					{Subject: "Toan", Score: 11},
				}
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrScore,
		},

		/* ----- REPO ERROR ------- */
		// ERROR repo.GetStudentByEmail
		{
			name:  "fail - repo.getStudentByEmail DB error",
			input: validStudent(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByEmailFn: func(studentEmail string) (*models.Student, error) {
						return nil, errors.New("db connection failed")
					},
				}
			},
			expectedErr: errors.New("db connection failed"),
		},
		// ERROR repo.AddStudent
		{
			name:  "fail - repo.addStudent DB error",
			input: validStudent(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByEmailFn: func(studentEmail string) (*models.Student, error) {
						return nil, errors.New("db connection failed")
					},
					AddStudentFn: func(student *models.Student) error {
						return errors.New("DB connection failed")
					},
				}
			},
			expectedErr: errors.New("db connection failed"),
		},
	}
	/* ------ RUN LOOP ------ */
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := services.NewStudentService(tc.setupRepo())
			err := svc.AddStudent(tc.input)

			if tc.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error [%v], got nil", tc.expectedErr)
					return
				}
				if !errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error [%v], got [%v]", tc.expectedErr, err)
				}
			}

		})
	}
}

// SEPERATED GENERATE ID - NEED CHECK STRUCT STATE AFTER CALLED, NOT ONLY CHECK ERROR
func TestAddStudent_AutoGenerateID(t *testing.T) {
	student := &models.Student{
		ID:          "",
		FullName:    "Le Trong Vu",
		DateOfBirth: time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
		Gender:      "male",
		Address:     "Ho Chi Minh",
		Class:       "DQT4",
		Email:       "letrongvu.work@gmail.com",
		Scores:      []*models.SubjectScore{},
	}

	repo := &MockStudentRepository{
		GetStudentByEmailFn: func(email string) (*models.Student, error) {
			return nil, errors.New("not found")
		},
		AddStudentFn: func(s *models.Student) error {
			return nil
		},
	}

	svc := services.NewStudentService(repo)
	err := svc.AddStudent(student)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// After call AddStudent, auto generate ID
	if student.ID == "" {
		t.Error("expected ID to be auto-generated, got empty string")
	}
}

func TestUpdateStudent(t *testing.T) {
	validStudent := func() *models.Student {
		return &models.Student{
			ID:          "existing-id-123",
			FullName:    "Le Trong Vu",
			DateOfBirth: time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
			Gender:      "male",
			Address:     "Ho Chi Minh",
			Class:       "DQT4",
			Email:       "letrongvu.work@gmail.com",
			Scores:      []*models.SubjectScore{},
		}
	}
	// HAPPY REPO - EXISTED STUDENT, EMAIL NOT BEING USED BY OTHERS
	happyRepo := func() *MockStudentRepository {
		return &MockStudentRepository{
			GetStudentByIDFn: func(studentID string) (*models.Student, error) {
				return &models.Student{ID: studentID}, nil
			},
			GetStudentByEmailFn: func(studentEmail string) (*models.Student, error) {
				return nil, errors.New("not found") // Email not being used
			},
			UpdateStudentFn: func(student *models.Student) error {
				return nil // update succesfully
			},
		}
	}

	tests := []struct {
		name        string
		input       *models.Student
		setupRepo   func() *MockStudentRepository
		expectedErr error
	}{
		/* ----- HAPPY PATH ----- */
		{
			name:        "success - valid student",
			input:       validStudent(),
			setupRepo:   happyRepo,
			expectedErr: nil,
		},
		{
			name:  "success - same email with student is updating",
			input: validStudent(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByIDFn: func(id string) (*models.Student, error) {
						return &models.Student{ID: id}, nil
					},
					GetStudentByEmailFn: func(email string) (*models.Student, error) {
						// return that student
						return &models.Student{
							ID:    "existing-id-123",
							Email: email,
						}, nil
					},
					UpdateStudentFn: func(s *models.Student) error {
						return nil
					},
				}
			},
			expectedErr: nil,
		},
		/* ----- NULL/ EMPTY ----- */
		{
			name:        "fail - student nil",
			input:       nil,
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrStudentData,
		},
		// CASE 1: EMPTY ID
		{
			name: "fail - empty ID",
			input: func() *models.Student {
				s := validStudent()
				s.ID = ""
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrIDRequired,
		},
		// CASE 2: EMPTY FULLNAME
		{
			name: "fail - Empty FullName",
			input: func() *models.Student {
				s := validStudent()
				s.FullName = ""
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrStudentData,
		},
		// CASE 3: EMPTY EMAIL
		{
			name: "fail - Empty Email",
			input: func() *models.Student {
				s := validStudent()
				s.Email = ""
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrStudentData,
		},
		/* ----- FORMAT VALIDATION ----- */

		// CASE 1: WRONG NAME FORMAT
		{
			name: "fail - incorrect name format",
			input: func() *models.Student {
				s := validStudent()
				s.FullName = "Le@Trong#Vu"
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrNameFormat,
		},
		// CASE 2: WRONG EMAIL FORMAT
		{
			name: "fail - wroincorrectng email format",
			input: func() *models.Student {
				s := validStudent()
				s.Email = "not-an-email"
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrEmailFormat,
		},
		// CASE 3: DOB IN FUTURE
		{
			name: "fail - DOB IN FUTURE",
			input: func() *models.Student {
				s := validStudent()
				s.DateOfBirth = time.Now().Add(24 * time.Hour)
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrValidDOB,
		},
		// CASE 4: DOB IN FUTURE (NEXT HOUR)
		{
			name: "fail - DOB IN FUTURE (NEXT HOUR)",
			input: func() *models.Student {
				s := validStudent()
				s.DateOfBirth = time.Now().Add(time.Hour)
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrValidDOB,
		},
		// CASE 5: EMPTY CLASS
		{
			name: "fail - empty class",
			input: func() *models.Student {
				s := validStudent()
				s.Class = ""
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrStudentClass,
		},
		// CASE 6: WRONG CLASS FORMAT
		{
			name: "fail -  incorrect class format",
			input: func() *models.Student {
				s := validStudent()
				s.Class = "INVALID!!!"
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrClassFormat,
		},
		// CASE 7: SCORE NEGATIVE
		{
			name: "fail - negative score",
			input: func() *models.Student {
				s := validStudent()
				s.Scores = []*models.SubjectScore{
					{Subject: "Math", Score: -5},
				}
				return s
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{}
			},
			expectedErr: services.ErrScore,
		},
		// CASE 8: SCORE GREATER THAN 10
		{
			name: "fail - score greater than 10",
			input: func() *models.Student {
				s := validStudent()
				s.Scores = []*models.SubjectScore{{Subject: "Math", Score: 11}}
				return s
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedErr: services.ErrScore,
		},

		/* ------ REPO CASES ------- */

		{
			name:  "fail - student does not existed",
			input: validStudent(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByIDFn: func(id string) (*models.Student, error) {
						return nil, nil // not found
					},
				}
			},
			expectedErr: services.ErrStudentNotFound,
		},
		{
			name:  "fail - email dublicated by others student",
			input: validStudent(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByIDFn: func(id string) (*models.Student, error) {
						return &models.Student{ID: id}, nil
					},
					GetStudentByEmailFn: func(email string) (*models.Student, error) {
						// return student using this email
						return &models.Student{
							ID:    "another-id-999",
							Email: email,
						}, nil
					},
				}
			},
			expectedErr: services.ErrEmailExisted,
		},
		/*------- REPO ERROR ----- */
		// CASE 1: REPO GetStudentByID
		{
			name:  "fail - repo.GetStudentByID error DB",
			input: validStudent(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByIDFn: func(id string) (*models.Student, error) {
						return nil, errors.New("db connection failed")
					},
				}
			},
			expectedErr: errors.New("db connection failed"),
		},
		// CASE 2: REPO UpdateStudent
		{
			name:  "fail - repo.UpdateStudent error DB",
			input: validStudent(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					GetStudentByIDFn: func(id string) (*models.Student, error) {
						return &models.Student{ID: id}, nil
					},
					GetStudentByEmailFn: func(email string) (*models.Student, error) {
						return nil, errors.New("not found")
					},
					UpdateStudentFn: func(s *models.Student) error {
						return errors.New("db connection failed")
					},
				}
			},
			expectedErr: errors.New("db connection failed"),
		},
	}
	// ── Run loop 
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := services.NewStudentService(tc.setupRepo())

			err := svc.UpdateStudent(tc.input)

			if tc.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error [%v], got nil", tc.expectedErr)
					return
				}
				if !errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error [%v], got [%v]", tc.expectedErr, err)
				}
			}
		})
	}
}

/* ANNOYNYMOUS FUNC
input := func() *modes.Student{
	s := validStudent() // coppy from base
	s.ID = "" 			// Replace field that need to test
	return s
}(), 					// <- () call func
*/
