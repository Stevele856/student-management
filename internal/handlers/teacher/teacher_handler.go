package teacherHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/student-management/internal/models"
	teacherModels "github.com/student-management/internal/models/teacher"
	teacherRepo "github.com/student-management/internal/repositories/teacher"
	"github.com/student-management/internal/services/teacher"
	"github.com/student-management/pkg/utils"
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
	case errors.Is(err, teacherRepo.ErrTeacherNotFound):
		return http.StatusNotFound // 404

	case errors.Is(err, teacher.ErrTeacherEmailExisted),
		errors.Is(err, teacher.ErrTeacherEmployeeIDExisted),
		errors.Is(err, teacher.ErrSubjectDuplicated),
		errors.Is(err, teacher.ErrClassDuplicate),
		errors.Is(err, teacherRepo.ErrTeacherAlreadyExists):
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
		errors.Is(err, teacherRepo.ErrInvalidPageSize),
		errors.Is(err, teacher.ErrTeacherIDRequired):
		return http.StatusBadRequest // 400

	default:
		return http.StatusInternalServerError // default fallback (500 interal error)
	}
}

// GET TEACHER
func (h *TeacherHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if q.Get("page") != "" || q.Get("page_size") != "" {
		h.GetTeachersPaginated(w, r)
		return
	}
	if len(q) == 0 {
		h.GetAllTeachers(w, r)
		return
	}
	// any query params -> filter by all provided conditions (AND semantics in predicate chain)
	h.FilterTeachers(w, r)
}

func (h *TeacherHandler) GetTeachersPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page := 1
	pageSize := 10

	if v := q.Get("page"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid page")
			return
		}
		page = parsed
	}

	if v := q.Get("page_size"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid page_size")
			return
		}
		pageSize = parsed
	}

	teachers, total, err := h.service.GetTeachersPaginated(page, pageSize)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"data":      teachers,
	})
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

/* -----------BULK UPLOAD CSV--------------- */
func (h *TeacherHandler) BulkAddTeachers(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file not found in form")
		return
	}
	defer file.Close()

	_, rows, err := utils.ReadCSV(file)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read CSV: "+err.Error())
		return
	}

	if len(rows) == 0 {
		writeError(w, http.StatusBadRequest, "CSV is empty")
		return
	}

	teacherRows := make(map[string][][]string)
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		teacherID := row[0]
		if teacherID == "" {
			writeError(w, http.StatusBadRequest, "teacher ID cannot be empty")
			return
		}
		teacherRows[teacherID] = append(teacherRows[teacherID], row)
	}

	teachers := make([]*teacherModels.Teacher, 0, len(teacherRows))
	for _, rowGroup := range teacherRows {
		t, err := teacherModels.TeacherFromCSVRows(rowGroup)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse teacher from CSV: %v", err))
			return
		}
		teachers = append(teachers, t)
	}

	if err := h.service.BulkAddTeachers(teachers); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": fmt.Sprintf("successfully added %d teachers", len(teachers)),
		"count":   len(teachers),
	})
}

/* -----------CSV EXPORT/DOWNLOAD--------------- */
func (h *TeacherHandler) ExportTeachers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var teachers []*teacherModels.Teacher
	var err error

	hasFilter := q.Get("name") != "" || q.Get("gender") != "" || q.Get("email") != "" ||
		q.Get("employee_id") != "" || q.Get("status") != "" || q.Get("hire_date_from") != "" ||
		q.Get("hire_date_to") != "" || len(q["class"]) > 0 || len(q["subject"]) > 0

	if hasFilter {
		filter := &teacherModels.FilterTeachers{
			Name:       q.Get("name"),
			Gender:     models.Gender(q.Get("gender")),
			Email:      q.Get("email"),
			EmployeeID: q.Get("employee_id"),
			Status:     teacherModels.TeacherStatus(q.Get("status")),
		}
		filter.ClassAssigned = q["class"]
		filter.SubjectTaught = q["subject"]

		const dateLayout = "2006-01-02"
		if v := q.Get("hire_date_from"); v != "" {
			t, parseErr := time.Parse(dateLayout, v)
			if parseErr != nil {
				writeError(w, http.StatusBadRequest, "invalid hire_date_from format, expected YYYY-MM-DD")
				return
			}
			filter.HireDateFrom = &t
		}
		if v := q.Get("hire_date_to"); v != "" {
			t, parseErr := time.Parse(dateLayout, v)
			if parseErr != nil {
				writeError(w, http.StatusBadRequest, "invalid hire_date_to format, expected YYYY-MM-DD")
				return
			}
			filter.HireDateTo = &t
		}

		teachers, err = h.service.FilterTeachers(filter)
	} else {
		teachers, err = h.service.GetAllTeachers()
	}

	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	if len(teachers) == 0 {
		writeError(w, http.StatusNotFound, "no teachers found to export")
		return
	}

	csvRows := [][]string{}
	for _, teacher := range teachers {
		csvRows = append(csvRows, teacher.ToCSVRows()...)
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"teachers.csv\"")
	w.WriteHeader(http.StatusOK)

	if err := utils.WriteCSV(w, teacherModels.CSVHeader(), csvRows); err != nil {
		fmt.Fprintf(w, "\nerror writing CSV: %v", err)
		return
	}
}

/*
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
*/
