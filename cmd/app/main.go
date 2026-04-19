package main

import (
	"log"
	"net/http"

	"github.com/student-management/config"
	"github.com/student-management/internal/handlers"
	"github.com/student-management/internal/repositories/student"
	"github.com/student-management/internal/services/student"
)

func main(){
	cfg := config.Load()

	repo, err := repositories.NewStudentMemoryRepo(cfg.DataFile)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}

	service := services.NewStudentService(repo)

	handler := handlers.NewStudentHandler(service)

	router := handlers.NewRouter(*handler)

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: router,
	}

	log.Printf("server starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}