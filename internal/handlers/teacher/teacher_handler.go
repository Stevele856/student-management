package TeacherHandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/student-management/internal/models"
	teacherModels "github.com/student-management/internal/models/teacher"
	teacherRepo "github.com/student-management/internal/repositories/teacher"
	"github.com/student-management/internal/services/teacher"
)

type TeacherHandler struct {
	service *teacher.TeacherService
}

func NewTeacherHandler(service *teacher.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		service: service,
	}
}

// Write json
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// Map service error -> HTTP status code
func serviceErrToStatus(err error) int {
	switch {
	case errors.Is(err, teacherRepo.ErrTeacherNotFound),
		errors.Is(err, teacher.ErrTeacherIDRequired):
		return http.StatusNotFound // 404

	case errors.Is(err, teacher.ErrTeacherEmailExisted),
		errors.Is(err, teacher.ErrTeacherEmployeeIDExisted),
		errors.Is(err, teacher.ErrSubjectDuplicated),
		errors.Is(err, teacher.ErrClassDuplicate):
		return http.StatusConflict // 409

	case errors.Is(err, teacher.ErrTeacherData),
		errors.Is(err, teacher.ErrNameRequired),
		errors.Is(err, teacher.ErrEmailRequired),
		errors.Is(err, teacher.ErrNameFormat),
		errors.Is(err, teacher.ErrEmailFormat),
		errors.Is(err, teacher.ErrValidDOB),
		errors.Is(err, teacher.ErrTeacherGender),
		errors.Is(err, teacher.ErrTeacherStatus),
		errors.Is(err, teacher.ErrEmployeeID),
		errors.Is(err, teacher.ErrFormatPhoneNumber),
		errors.Is(err, teacher.ErrSubjectFormat),
		errors.Is(err, teacher.ErrClassFormat),
		errors.Is(err, teacher.ErrSubjectRequired),
		errors.Is(err, teacher.ErrTeacherHireDate),
		errors.Is(err, teacher.ErrClassRequired),
		errors.Is(err, teacherRepo.ErrInvalidPage),
		errors.Is(err, teacherRepo.ErrInvalidPageSize):
		return http.StatusBadRequest // 400

	default:
		return http.StatusInternalServerError // default fallback (500 interal error)
	}
}

// GET TEACHER
func (h *TeacherHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if len(q) == 0{
		h.GetAllTeachers(w,r)
	}
	// any query params -> filter by all provided conditions (AND semantics in predicate chain)
	h.FilterTeachers(w,r)
}

// GET ALL TEACHERS
func (h *TeacherHandler) GetAllTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.service.GetAllTeachers()

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, &teachers)
}

// GET TEACHER BY EMAIL
func (h *TeacherHandler) GetTeacherByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	teacher, err := h.service.GetTeacherByEmail(email)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	if teacher == nil {
		writeError(w, http.StatusNotFound, "teacher not found")
		return
	}

	writeJSON(w, http.StatusOK, teacher)
}

// GET TEACHER BY SUBJECT
func (h *TeacherHandler) GetTeacherBySubject(w http.ResponseWriter, r *http.Request) {
	subject := r.URL.Query().Get("subject")

	teachers, err := h.service.GetTeacherAssignedBySubject(subject)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, teachers)
}

// GET TEACHER BY CLASS
func (h *TeacherHandler) GetTeacherByClass(w http.ResponseWriter, r *http.Request) {
	classAssigned := r.URL.Query().Get("class")

	teachers, err := h.service.GetTeacherByAssignedClass(classAssigned)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, teachers)
}

// GET TEACHER BY STATUS
func (h *TeacherHandler) GetTeacherByStatus(w http.ResponseWriter, r *http.Request) {
	status := teacherModels.TeacherStatus(r.URL.Query().Get("status"))

	teachers, err := h.service.GetTeacherByStatus(status)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, teachers)
}

// FILTER TEACHER
func (h *TeacherHandler) FilterTeachers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := &teacherModels.FilterTeachers{
		Name:       q.Get("name"),
		Gender:     models.Gender(q.Get("gender")),
		Email:      q.Get("email"),
		EmployeeID: q.Get("employee_id"),
		Status:     teacherModels.TeacherStatus(q.Get("status")),
	}

	// if repeated query params:
	// ?class=A1&class=A2 or ?subject=math&subject=english
	filter.ClassAssigned = q["class"]
	filter.SubjectTaught = q["subject"]

	const dateLayout = "2006-01-02"

	if v := q.Get("hire_date_from"); v != "" {
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid hire_date_from format, expected YYYY-MM-DD")
			return
		}
		filter.HireDateFrom = &t
	}

	if v := q.Get("hire_date_to"); v != "" {
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid hire_date_to format, expected YYYY-MM-DD")
			return
		}
		filter.HireDateTo = &t
	}

	teachers, err := h.service.FilterTeachers(filter)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, teachers)
}

// GET TEACHER BY ID
func (h *TeacherHandler) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	teacher, err := h.service.GetTeacherByID(id)

	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	if teacher == nil {
		writeError(w, http.StatusNotFound, "teacher not found")
		return
	}

	writeJSON(w, http.StatusOK, teacher)
}

// GET TEACHER BY EMPL ID
func (h *TeacherHandler) GetTeacherByEmployeeID(w http.ResponseWriter, r *http.Request) {
	employeeID := r.PathValue("employee_id")

	teacher, err := h.service.GetTeacherByEmployeeID(employeeID)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	if teacher == nil {
		writeError(w, http.StatusNotFound, "teacher not found")
		return
	}

	writeJSON(w, http.StatusOK, teacher)
}

// ADD TEACHER
func (h *TeacherHandler) AddTeacher(w http.ResponseWriter, r *http.Request) {
	teacher := teacherModels.Teacher{}
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.service.AddTeacher(&teacher); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, &teacher)
}

// UPDATE TEACHER
func (h *TeacherHandler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	teacher := teacherModels.Teacher{}
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()
	teacher.ID = id

	if err := h.service.UpdateTeacher(&teacher); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, &teacher)

}

// DELETE TEACHER
func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteTeacher(id); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "teacher delete successfully"})
}
