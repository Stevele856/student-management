package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/student-management/internal/services"
)

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(service *services.StudentService) *StudentHandler{
	return &StudentHandler{
		service: service,
	}
}

/* --------WRITE JSON--------- */

func writeJSON(w http.ResponseWriter, status int, data any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string){
	writeJSON(w, status, map[string]string{"error": msg})
}

// Map service errors → HTTP status codes
func serviceErrToStatus(err error) int {
	switch {
	case errors.Is(err, services.ErrStudentNotFound),
		errors.Is(err, services.ErrStudentID),
		errors.Is(err, services.ErrStudentEmail),
		errors.Is(err, services.ErrSubjectEmpty):
		return http.StatusNotFound
 
	case errors.Is(err, services.ErrEmailExisted),
		errors.Is(err, services.ErrSubjectAlreadyExisted),
		errors.Is(err, services.ErrDublicatedSubject):
		return http.StatusConflict
 
	case errors.Is(err, services.ErrStudentInfo),
		errors.Is(err, services.ErrNameFormat),
		errors.Is(err, services.ErrEmailFormat),
		errors.Is(err, services.ErrClassFormat),
		errors.Is(err, services.ErrSubjectFormat),
		errors.Is(err, services.ErrValidDOB),
		errors.Is(err, services.ErrScore),
		errors.Is(err, services.ErrMaxScore),
		errors.Is(err, services.ErrIDRequired),
		errors.Is(err, services.ErrNameRequired),
		errors.Is(err, services.ErrStudentClass):
		return http.StatusBadRequest
 
	default:
		return http.StatusInternalServerError
	}
}

// GET 

func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request){
	students, err := h.service.GetAllStudents()

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, students)
}

func (h *StudentHandler) GetStudentByEmail(w http.ResponseWriter, r *http.Request){
	email := r.URL.Query().Get("email")

	student, err := h.service.GetStudentByEmail(email)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	if student == nil{
		writeError(w, http.StatusNotFound, "student not found")
		return
	}

	writeJSON(w, http.StatusOK, student)
}

// GET STUDENTS IF HAVE EMAIL QUERY PARAM
func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request){
	email := r.URL.Query().Get("email")
	if email != ""{
		h.GetStudentByEmail(w,r)
		return
	}	

	h.GetAllStudents(w,r)
}

func (h *StudentHandler) GetStudentByID(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")

	student, err := h.service.GetStudentByID(id)

	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	if student == nil {
		writeError(w, http.StatusNotFound, "student not found")
		return
	}

	writeJSON(w, http.StatusOK, student)
}

