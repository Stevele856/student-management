package handlers

import "net/http"

func NewRouter(h StudentHandler) *http.ServeMux{
	mux := http.NewServeMux()

	mux.HandleFunc("GET /students", h.GetAllStudents)
	mux.HandleFunc("GET /students/{id}", h.GetStudentByID)
	return mux
}