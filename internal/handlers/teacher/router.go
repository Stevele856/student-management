package TeacherHandler

import "net/http"

func NewRouter(h TeacherHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /teachers", h.GetTeachers)
	mux.HandleFunc("GET /teachers/{id}", h.GetTeacherByID)
	mux.HandleFunc("GET /teachers/employee/{employee_id}", h.GetTeacherByEmployeeID)

	mux.HandleFunc("POST /teachers", h.AddTeacher)
	mux.HandleFunc("PUT /teachers/{id}", h.UpdateTeacher)
	mux.HandleFunc("DELETE /teachers/{id}", h.DeleteTeacher)

	return mux
}
