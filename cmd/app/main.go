package main

import (
	"log"
	"net/http"

	"github.com/student-management/config"

	studentHandler "github.com/student-management/internal/handlers/student"
	teacherHandler "github.com/student-management/internal/handlers/teacher"
	studentRepo "github.com/student-management/internal/repositories/student"
	teacherRepo "github.com/student-management/internal/repositories/teacher"
	studentService "github.com/student-management/internal/services/student"
	teacherService "github.com/student-management/internal/services/teacher"
)

func main() {
	cfg := config.Load()

	studentRepository, err := studentRepo.NewStudentMemoryRepo(cfg.StudentData)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}

	studentService := studentService.NewStudentService(studentRepository)
	studentHandlers := studentHandler.NewStudentHandler(studentService)
	studentRouter := studentHandler.NewRouter(*studentHandlers)

	teacherRepository, err := teacherRepo.NewTeacherMemoryRepo(cfg.TeacherData)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}

	teacherService := teacherService.NewTeacherService(teacherRepository)
	teacherHandlers := teacherHandler.NewTeacherHandler(teacherService)
	teacherRouter := teacherHandler.NewRouter(*teacherHandlers)

	rootMux := http.NewServeMux()
	rootMux.Handle("/students", studentRouter)
	rootMux.Handle("/teachers", teacherRouter)
	rootMux.Handle("/students/", studentRouter)
	rootMux.Handle("/teachers/", teacherRouter)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: rootMux,
	}

	log.Printf("server starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}

}
