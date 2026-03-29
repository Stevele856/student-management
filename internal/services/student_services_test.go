package services

import (
	"testing"
	"time"

	"github.com/student-management/internal/models"
)

func TestAddStudent(t *testing.T) {
	// HELPER CREATE VALID STUENT
	validStudent := func() *models.Student {
		return &models.Student{
			FullName:    "Le Trong Vu",
			DateOfBirth: time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
			Gender: "male",
			Address: "Ho Chi Minh",
			Class: "DQT4",
			Email: "letrongvu.work@gmail.com",
			Scores: []*models.SubjectScore{},
		}
	}

	// HELPER CREATE DEFAULT MOCK REPO
	happyRepo := func () *MockStudentRepository  {
		
	}
}
