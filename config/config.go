package config

import "os"

type Config struct {
	Port        string
	StudentData string
	TeacherData string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	studentData := os.Getenv("STUDENT_DATA_FILE")
	if studentData == "" {
		studentData = "static/students.json"
	}

	teacherData := os.Getenv("TEACHER_DATA_FILE")
	if teacherData == "" {
		teacherData = "static/teachers.json"
	}

	return &Config{
		Port:     port,
		StudentData: studentData,
		TeacherData: teacherData,
	}
}
