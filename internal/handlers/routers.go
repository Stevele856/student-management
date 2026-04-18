package handlers

import "net/http"

func NewRouter(h StudentHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /students", h.GetStudents)
	mux.HandleFunc("GET /students/export", h.ExportStudents)
	mux.HandleFunc("GET /students/{id}", h.GetStudentByID)

	mux.HandleFunc("POST /students", h.AddStudent)
	mux.HandleFunc("POST /students/bulk-upload", h.BulkAddStudents)
	mux.HandleFunc("PUT /students/{id}", h.UpdateStudent)
	mux.HandleFunc("DELETE /students/{id}", h.DeleteStudent)

	mux.HandleFunc("GET /students/{id}/scores", h.GetScoresByStudentID)
	mux.HandleFunc("GET /students/{id}/scores/{subject}", h.GetScoresBySubject)

	mux.HandleFunc("POST /students/{id}/scores", h.AddScore)
	mux.HandleFunc("PUT /students/{id}/scores/{subject}", h.UpdateScore)
	mux.HandleFunc("DELETE /students/{id}/scores/{subject}", h.DeleteScore)
	return mux
}
