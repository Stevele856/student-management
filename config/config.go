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

	StudentData := os.Getenv("STUDENT_DATA_FILE")
	if StudentData == "" {
		StudentData = "static/students.json"
	}

	TeacherData := os.Getenv("TEACHER_DATA_FILE")
	if TeacherData == "" {
		TeacherData = "static/teachers.json"
	}

	return &Config{
		Port:     port,
		StudentData: StudentData,
		TeacherData: TeacherData,
	}
}
