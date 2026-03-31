package services_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/student-management/internal/models"
	"github.com/student-management/internal/predicate"
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

func TestDeleteStudent(t *testing.T) {

    // ── Helper 
    happyRepo := func() *MockStudentRepository {
        return &MockStudentRepository{
            GetStudentByIDFn: func(id string) (*models.Student, error) {
                return &models.Student{ID: id}, nil // existed
            },
            DeleteStudentFn: func(id string) error {
                return nil // delete succesfully
            },
        }
    }

    tests := []struct {
        name        string
        input       string // studentID (different with AddStudent and UpdateStudent)
        setupRepo   func() *MockStudentRepository
        expectedErr error
    }{
        // ── HAPPY PATH 
        {
            name:        "success - delete succesfully",
            input:       "existing-id-123",
            setupRepo:   happyRepo,
            expectedErr: nil,
        },
        {
            name:        "success - ID has trim space",
            input:       "  existing-id-123  ", // service TrimSpace
            setupRepo:   happyRepo,
            expectedErr: nil,
        },

        // ── EMPTY INPUT 
        {
            name:        "fail - empty ID",
            input:       "",
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrStudentID,
        },
        {
            name:        "fail - ID has only space",
            input:       "   ",
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrStudentID,
        },

        // ── REPO CASES 
        {
            name:  "fail - student not existing",
            input: "non-existing-id",
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
            name:  "fail - repo.GetStudentByID DB error",
            input: "existing-id-123",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
        {
            name:  "fail - repo.DeleteStudent DB error",
            input: "existing-id-123",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{ID: id}, nil
                    },
                    DeleteStudentFn: func(id string) error {
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

            err := svc.DeleteStudent(tc.input)

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



func TestGetAllStudents(t *testing.T) {
    tests := []struct {
        name        string
        setupRepo   func() *MockStudentRepository
        expected    []*models.Student
        expectedErr error
    }{
        {
            name: "success - return list of students",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetAllStudentsFn: func() ([]*models.Student, error) {
                        return []*models.Student{
                            {ID: "1", FullName: "Le Trong Vu"},
                            {ID: "2", FullName: "Jenifer"},
                        }, nil
                    },
                }
            },
            expected: []*models.Student{
                {ID: "1", FullName: "Le Trong Vu"},
                {ID: "2", FullName: "Jenifer"},
            },
            expectedErr: nil,
        },
        {
            name: "success - return empty list",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetAllStudentsFn: func() ([]*models.Student, error) {
                        return []*models.Student{}, nil
                    },
                }
            },
            expected:    []*models.Student{},
            expectedErr: nil,
        },
        {
            name: "fail - repo returns db error",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetAllStudentsFn: func() ([]*models.Student, error) {
                        return nil, errors.New("db connection failed")
                    },
                }
            },
            expected:    nil,
            expectedErr: errors.New("db connection failed"),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            svc := services.NewStudentService(tc.setupRepo())

            result, err := svc.GetAllStudents()

            if tc.expectedErr == nil {
                if err != nil {
                    t.Errorf("expected no error, got: %v", err)
                }
                if len(result) != len(tc.expected) {
                    t.Errorf("expected %d students, got %d", len(tc.expected), len(result))
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

// ─────────────────────────────────────────────────────────────────────────────

func TestGetStudentByID(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        setupRepo   func() *MockStudentRepository
        expectedErr error
    }{
        {
            name:  "success - student found",
            input: "existing-id-123",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{ID: id}, nil
                    },
                }
            },
            expectedErr: nil,
        },
        {
            name:        "fail - empty ID",
            input:       "",
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrStudentID,
        },
        {
            name:        "fail - whitespace only ID",
            input:       "   ",
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrStudentID,
        },
        {
            name:  "fail - repo returns db error",
            input: "existing-id-123",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            svc := services.NewStudentService(tc.setupRepo())

            _, err := svc.GetStudentByID(tc.input)

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

func TestGetStudentByEmail(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        setupRepo   func() *MockStudentRepository
        expectedErr error
    }{
        {
            name:  "success - student found",
            input: "letrongvu.work@gmail.com",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByEmailFn: func(email string) (*models.Student, error) {
                        return &models.Student{Email: email}, nil
                    },
                }
            },
            expectedErr: nil,
        },
        {
            name:        "fail - empty email",
            input:       "",
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrStudentEmail,
        },
        {
            name:        "fail - whitespace only email",
            input:       "   ",
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrStudentEmail,
        },
        {
            name:        "fail - invalid email format",
            input:       "not-an-email",
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrEmailFormat,
        },
        {
            name:  "fail - repo returns db error",
            input: "letrongvu.work@gmail.com",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByEmailFn: func(email string) (*models.Student, error) {
                        return nil, errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            svc := services.NewStudentService(tc.setupRepo())

            _, err := svc.GetStudentByEmail(tc.input)

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

func TestAddScore(t *testing.T) {

    // ── Helpers ───────────────────────────────────────────────────────────
    validScore := func() *models.SubjectScore {
        return &models.SubjectScore{
            Subject: "Math",
            Score:   8.5,
        }
    }

    happyRepo := func() *MockStudentRepository {
        return &MockStudentRepository{
            GetStudentByIDFn: func(id string) (*models.Student, error) {
                return &models.Student{
                    ID:     id,
                    Scores: []*models.SubjectScore{}, // chưa có môn nào
                }, nil
            },
            AddScoreFn: func(id string, score *models.SubjectScore) error {
                return nil
            },
        }
    }

    tests := []struct {
        name        string
        inputID     string
        inputScore  *models.SubjectScore
        setupRepo   func() *MockStudentRepository
        expectedErr error
    }{
        // ── HAPPY PATH ────────────────────────────────────────────────────
        {
            name:        "success - valid score",
            inputID:     "existing-id-123",
            inputScore:  validScore(),
            setupRepo:   happyRepo,
            expectedErr: nil,
        },
        {
            name:    "success - score is 0",
            inputID: "existing-id-123",
            inputScore: &models.SubjectScore{
                Subject: "Math",
                Score:   0, // 0 là hợp lệ
            },
            setupRepo:   happyRepo,
            expectedErr: nil,
        },
        {
            name:    "success - score is 10",
            inputID: "existing-id-123",
            inputScore: &models.SubjectScore{
                Subject: "Math",
                Score:   10, // 10 là hợp lệ
            },
            setupRepo:   happyRepo,
            expectedErr: nil,
        },

        // ── EMPTY / NIL INPUT ─────────────────────────────────────────────
        {
            name:        "fail - empty student ID",
            inputID:     "",
            inputScore:  validScore(),
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrIDRequired,
        },
        {
            name:        "fail - whitespace only student ID",
            inputID:     "   ",
            inputScore:  validScore(),
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrIDRequired,
        },
        {
            name:        "fail - score is nil",
            inputID:     "existing-id-123",
            inputScore:  nil,
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrStudentData,
        },

        // ── SUBJECT VALIDATION ────────────────────────────────────────────
        {
            name:    "fail - empty subject",
            inputID: "existing-id-123",
            inputScore: &models.SubjectScore{
                Subject: "",
                Score:   8.5,
            },
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrSubjectFormat,
        },
        {
            name:    "fail - invalid subject format",
            inputID: "existing-id-123",
            inputScore: &models.SubjectScore{
                Subject: "Math@#$",
                Score:   8.5,
            },
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrSubjectFormat,
        },

        // ── SCORE VALIDATION ──────────────────────────────────────────────
        {
            name:    "fail - score below 0",
            inputID: "existing-id-123",
            inputScore: &models.SubjectScore{
                Subject: "Math",
                Score:   -1,
            },
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrScore,
        },
        {
            name:    "fail - score above 10",
            inputID: "existing-id-123",
            inputScore: &models.SubjectScore{
                Subject: "Math",
                Score:   11,
            },
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrScore,
        },

        // ── BUSINESS RULES ────────────────────────────────────────────────
        {
            name:       "fail - student not found",
            inputID:    "non-existing-id",
            inputScore: validScore(),
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, nil
                    },
                }
            },
            expectedErr: services.ErrStudentNotFound,
        },
        {
            name:       "fail - max 10 scores reached",
            inputID:    "existing-id-123",
            inputScore: validScore(),
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        // student đã có đủ 10 môn
                        scores := make([]*models.SubjectScore, 10)
                        for i := range scores {
                            scores[i] = &models.SubjectScore{
                                Subject: fmt.Sprintf("Subject%d", i),
                                Score:   8,
                            }
                        }
                        return &models.Student{ID: id, Scores: scores}, nil
                    },
                }
            },
            expectedErr: services.ErrMaxScore,
        },
        {
            name:       "fail - duplicate subject",
            inputID:    "existing-id-123",
            inputScore: validScore(), // Subject: "Math"
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{
                            ID: id,
                            Scores: []*models.SubjectScore{
                                {Subject: "Math", Score: 9}, // đã có Math rồi
                            },
                        }, nil
                    },
                }
            },
            expectedErr: services.ErrSubjectAlreadyExisted,
        },
        {
            name:       "fail - duplicate subject case insensitive",
            inputID:    "existing-id-123",
            inputScore: &models.SubjectScore{Subject: "MATH", Score: 8}, // MATH == math
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{
                            ID: id,
                            Scores: []*models.SubjectScore{
                                {Subject: "math", Score: 9},
                            },
                        }, nil
                    },
                }
            },
            expectedErr: services.ErrSubjectAlreadyExisted,
        },

        // ── REPO ERROR ────────────────────────────────────────────────────
        {
            name:       "fail - repo.GetStudentByID db error",
            inputID:    "existing-id-123",
            inputScore: validScore(),
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
        {
            name:       "fail - repo.AddScore db error",
            inputID:    "existing-id-123",
            inputScore: validScore(),
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{
                            ID:     id,
                            Scores: []*models.SubjectScore{},
                        }, nil
                    },
                    AddScoreFn: func(id string, score *models.SubjectScore) error {
                        return errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            svc := services.NewStudentService(tc.setupRepo())

            err := svc.AddScore(tc.inputID, tc.inputScore)

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

func TestUpdateScore(t *testing.T) {

    validScore := func() *models.SubjectScore {
        return &models.SubjectScore{Subject: "Math", Score: 9}
    }

    happyRepo := func() *MockStudentRepository {
        return &MockStudentRepository{
            GetStudentByIDFn: func(id string) (*models.Student, error) {
                return &models.Student{ID: id}, nil
            },
            UpdateScoreFn: func(id string, score *models.SubjectScore) error {
                return nil
            },
        }
    }

    tests := []struct {
        name        string
        inputID     string
        inputScore  *models.SubjectScore
        setupRepo   func() *MockStudentRepository
        expectedErr error
    }{
        // ── HAPPY PATH ────────────────────────────────────────────────────
        {
            name:        "success - valid score update",
            inputID:     "existing-id-123",
            inputScore:  validScore(),
            setupRepo:   happyRepo,
            expectedErr: nil,
        },

        // ── EMPTY / NIL INPUT ─────────────────────────────────────────────
        {
            name:        "fail - empty student ID",
            inputID:     "",
            inputScore:  validScore(),
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrIDRequired,
        },
        {
            name:        "fail - score is nil",
            inputID:     "existing-id-123",
            inputScore:  nil,
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrStudentData,
        },

        // ── VALIDATION ────────────────────────────────────────────────────
        {
            name:    "fail - invalid subject format",
            inputID: "existing-id-123",
            inputScore: &models.SubjectScore{
                Subject: "Math@#$",
                Score:   8,
            },
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrSubjectFormat,
        },
        {
            name:    "fail - score above 10",
            inputID: "existing-id-123",
            inputScore: &models.SubjectScore{
                Subject: "Math",
                Score:   11,
            },
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrScore,
        },

        // ── REPO ERROR ────────────────────────────────────────────────────
        {
            name:       "fail - student not found",
            inputID:    "non-existing-id",
            inputScore: validScore(),
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, nil
                    },
                }
            },
            expectedErr: services.ErrStudentNotFound,
        },
        {
            name:       "fail - repo.UpdateScore db error",
            inputID:    "existing-id-123",
            inputScore: validScore(),
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{ID: id}, nil
                    },
                    UpdateScoreFn: func(id string, score *models.SubjectScore) error {
                        return errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            svc := services.NewStudentService(tc.setupRepo())

            err := svc.UpdateScore(tc.inputID, tc.inputScore)

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

func TestDeleteScore(t *testing.T) {

    happyRepo := func() *MockStudentRepository {
        return &MockStudentRepository{
            GetStudentByIDFn: func(id string) (*models.Student, error) {
                return &models.Student{ID: id}, nil
            },
            DeleteScoreFn: func(id, subject string) error {
                return nil
            },
        }
    }

    tests := []struct {
        name        string
        inputID     string
        inputSubject string
        setupRepo   func() *MockStudentRepository
        expectedErr error
    }{
        // ── HAPPY PATH ────────────────────────────────────────────────────
        {
            name:         "success - valid delete",
            inputID:      "existing-id-123",
            inputSubject: "Math",
            setupRepo:    happyRepo,
            expectedErr:  nil,
        },

        // ── EMPTY INPUT ───────────────────────────────────────────────────
        {
            name:         "fail - empty student ID",
            inputID:      "",
            inputSubject: "Math",
            setupRepo:    func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr:  services.ErrIDRequired,
        },
        {
            name:         "fail - empty subject",
            inputID:      "existing-id-123",
            inputSubject: "",
            setupRepo:    func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr:  services.ErrSubjectEmpty,
        },
        {
            name:         "fail - invalid subject format",
            inputID:      "existing-id-123",
            inputSubject: "Math@#$",
            setupRepo:    func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr:  services.ErrSubjectFormat,
        },

        // ── REPO CASES ────────────────────────────────────────────────────
        {
            name:         "fail - student not found",
            inputID:      "non-existing-id",
            inputSubject: "Math",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, nil
                    },
                }
            },
            expectedErr: services.ErrStudentNotFound,
        },
        {
            name:         "fail - repo.DeleteScore db error",
            inputID:      "existing-id-123",
            inputSubject: "Math",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{ID: id}, nil
                    },
                    DeleteScoreFn: func(id, subject string) error {
                        return errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            svc := services.NewStudentService(tc.setupRepo())

            err := svc.DeleteScore(tc.inputID, tc.inputSubject)

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

func TestGetScoresByStudentID(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        setupRepo   func() *MockStudentRepository
        expectedErr error
    }{
        {
            name:  "success - scores found",
            input: "existing-id-123",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{ID: id}, nil
                    },
                    GetScoresByStudentIDFn: func(id string) ([]*models.SubjectScore, error) {
                        return []*models.SubjectScore{
                            {Subject: "Math", Score: 8},
                        }, nil
                    },
                }
            },
            expectedErr: nil,
        },
        {
            name:        "fail - empty student ID",
            input:       "",
            setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr: services.ErrIDRequired,
        },
        {
            name:  "fail - student not found",
            input: "non-existing-id",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, nil
                    },
                }
            },
            expectedErr: services.ErrStudentNotFound,
        },
        {
            name:  "fail - repo db error",
            input: "existing-id-123",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            svc := services.NewStudentService(tc.setupRepo())

            _, err := svc.GetScoresByStudentID(tc.input)

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

// ─────────────────────────────────────────────────────────────────────────────

func TestGetScoresBySubject(t *testing.T) {
    tests := []struct {
        name         string
        inputID      string
        inputSubject string
        setupRepo    func() *MockStudentRepository
        expectedErr  error
    }{
        {
            name:         "success - score found",
            inputID:      "existing-id-123",
            inputSubject: "Math",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return &models.Student{ID: id}, nil
                    },
                    GetScoresBySubjectFn: func(id, subject string) (*models.SubjectScore, error) {
                        return &models.SubjectScore{Subject: subject, Score: 8}, nil
                    },
                }
            },
            expectedErr: nil,
        },
        {
            name:         "fail - empty student ID",
            inputID:      "",
            inputSubject: "Math",
            setupRepo:    func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr:  services.ErrIDRequired,
        },
        {
            name:         "fail - empty subject",
            inputID:      "existing-id-123",
            inputSubject: "",
            setupRepo:    func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr:  services.ErrSubjectEmpty,
        },
        {
            name:         "fail - invalid subject format",
            inputID:      "existing-id-123",
            inputSubject: "Math@#$",
            setupRepo:    func() *MockStudentRepository { return &MockStudentRepository{} },
            expectedErr:  services.ErrSubjectFormat,
        },
        {
            name:         "fail - student not found",
            inputID:      "non-existing-id",
            inputSubject: "Math",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, nil
                    },
                }
            },
            expectedErr: services.ErrStudentNotFound,
        },
        {
            name:         "fail - repo db error",
            inputID:      "existing-id-123",
            inputSubject: "Math",
            setupRepo: func() *MockStudentRepository {
                return &MockStudentRepository{
                    GetStudentByIDFn: func(id string) (*models.Student, error) {
                        return nil, errors.New("db connection failed")
                    },
                }
            },
            expectedErr: errors.New("db connection failed"),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            svc := services.NewStudentService(tc.setupRepo())

            _, err := svc.GetScoresBySubject(tc.inputID, tc.inputSubject)

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

func TestFilterStudents(t *testing.T) {
	// ── Helper ─────────────────────────────────────────────────────────────
	validFilter := func() *models.FilterStudents {
		return &models.FilterStudents{
			Name:        "",
			Class:       "",
			YearOfBirth: 0,
			Gender:      "",
			Address:     "",
			MinAvgScore: 0,
			MaxAvgScore: 10,
			StudentRank: "",
		}
	}

	// Sample students 
	sampleStudents := []*models.Student{
		{
			ID:          "1",
			FullName:    "Le Trong Vu",
			DateOfBirth: time.Date(2005, 5, 15, 0, 0, 0, 0, time.UTC),
			Gender:      "male",
			Address:     "Ho Chi Minh",
			Class:       "DQT4",
			Scores:      []*models.SubjectScore{{Subject: "Toan", Score: 8.5}, {Subject: "Van", Score: 7.0}},
		},
		{
			ID:          "2",
			FullName:    "Nguyen Thi Lan",
			DateOfBirth: time.Date(2006, 3, 20, 0, 0, 0, 0, time.UTC),
			Gender:      "female",
			Address:     "Ha Noi",
			Class:       "DQT4",
			Scores:      []*models.SubjectScore{{Subject: "Toan", Score: 9.2}, {Subject: "Ly", Score: 8.5}},
		},
		{
			ID:          "3",
			FullName:    "Tran Van An",
			DateOfBirth: time.Date(2005, 11, 10, 0, 0, 0, 0, time.UTC),
			Gender:      "male",
			Address:     "Ho Chi Minh",
			Class:       "DQT5",
			Scores:      []*models.SubjectScore{{Subject: "Toan", Score: 6.5}, {Subject: "Hoa", Score: 7.0}},
		},
	}

	happyRepo := func() *MockStudentRepository {
		return &MockStudentRepository{
			FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
				return sampleStudents, nil
			},
		}
	}

	tests := []struct {
		name        string
		input       *models.FilterStudents
		setupRepo   func() *MockStudentRepository
		expectedLen int
		expectedErr error
	}{
		/* ── HAPPY PATH ───────────────────────────────────────────────── */
		{
			name:        "success - no filter (return all)",
			input:       validFilter(),
			setupRepo:   happyRepo,
			expectedLen: 3,
			expectedErr: nil,
		},
		{
			name: "success - filter by Name (partial match)",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.Name = "Trong Vu"
				return f
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
						return []*models.Student{sampleStudents[0]}, nil
					},
				}
			},
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			name: "success - filter by Class",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.Class = "DQT4"
				return f
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
						return []*models.Student{sampleStudents[0], sampleStudents[1]}, nil
					},
				}
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name: "success - filter by YearOfBirth",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.YearOfBirth = 2005
				return f
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
						return []*models.Student{sampleStudents[0], sampleStudents[2]}, nil
					},
				}
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name: "success - filter by Gender",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.Gender = "female"
				return f
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
						return []*models.Student{sampleStudents[1]}, nil
					},
				}
			},
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			name: "success - filter by Address",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.Address = "Ho Chi Minh"
				return f
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
						return []*models.Student{sampleStudents[0], sampleStudents[2]}, nil
					},
				}
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name: "success - filter by MinAvgScore and MaxAvgScore",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.MinAvgScore = 8.0
				f.MaxAvgScore = 9.5
				return f
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
						// chỉ trả về student có điểm trung bình trong khoảng
						return []*models.Student{sampleStudents[0], sampleStudents[1]}, nil
					},
				}
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name: "success - filter by StudentRank",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.StudentRank = "Excellent" // hoặc giá trị Rank bạn đang dùng
				return f
			}(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
						return []*models.Student{sampleStudents[1]}, nil
					},
				}
			},
			expectedLen: 1,
			expectedErr: nil,
		},

		/* ── INVALID INPUT ─────────────────────────────────────────────── */
		{
			name:        "fail - filter is nil",
			input:       nil,
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedLen: 0,
			expectedErr: services.ErrStudentData, // hoặc ErrInvalidFilter nếu bạn có
		},
		{
			name: "fail - invalid MinAvgScore > MaxAvgScore",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.MinAvgScore = 9.0
				f.MaxAvgScore = 5.0
				return f
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedLen: 0,
			expectedErr: services.ErrStudentData, // nếu bạn có error này, còn không thì dùng ErrStudentData
		},
		{
			name: "fail - invalid Gender",
			input: func() *models.FilterStudents {
				f := validFilter()
				f.Gender = "unknown"
				return f
			}(),
			setupRepo:   func() *MockStudentRepository { return &MockStudentRepository{} },
			expectedLen: 0,
			expectedErr: services.ErrInvalidGender, // tùy bạn có định nghĩa hay không
		},

		/* ── REPO ERROR ────────────────────────────────────────────────── */
		{
			name:  "fail - repo returns database error",
			input: validFilter(),
			setupRepo: func() *MockStudentRepository {
				return &MockStudentRepository{
					FilterStudentsFn: func(p predicate.PredicateStudent) ([]*models.Student, error) {
						return nil, errors.New("db connection failed")
					},
				}
			},
			expectedLen: 0,
			expectedErr: errors.New("db connection failed"),
		},
	}

	// ── Run test cases ─────────────────────────────────────────────────────
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := services.NewStudentService(tc.setupRepo())

			result, err := svc.FilterStudents(tc.input)

			if tc.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					return
				}
				if len(result) != tc.expectedLen {
					t.Errorf("expected %d students, got %d", tc.expectedLen, len(result))
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
