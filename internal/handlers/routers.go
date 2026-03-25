package handlers

import "net/http"

func NewRouter(h StudentHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /students", h.GetAllStudents)
	mux.HandleFunc("GET /students/{id}", h.GetStudentByID)

	mux.HandleFunc("POST /students", h.AddStudent)
	mux.HandleFunc("PUT /students/{id}", h.UpdateStudent)
	mux.HandleFunc("DELETE /students/{id}", h.DeleteStudent)
	return mux
}
