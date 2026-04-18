package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/student-management/internal/models"
	"github.com/student-management/internal/services"
	"github.com/student-management/pkg/utils"
)

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(service *services.StudentService) *StudentHandler {
	return &StudentHandler{
		service: service,
	}
}

/* --------WRITE JSON--------- */

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// Map service errors → HTTP status codes
func serviceErrToStatus(err error) int {
	switch {
	case errors.Is(err, services.ErrStudentNotFound),
		errors.Is(err, services.ErrIDRequired),
		errors.Is(err, services.ErrStudentEmail),
		errors.Is(err, services.ErrSubjectEmpty):
		return http.StatusNotFound

	case errors.Is(err, services.ErrEmailExisted),
		errors.Is(err, services.ErrSubjectAlreadyExisted),
		errors.Is(err, services.ErrDublicatedSubject):
		return http.StatusConflict

	case errors.Is(err, services.ErrStudentData),
		errors.Is(err, services.ErrNameFormat),
		errors.Is(err, services.ErrEmailFormat),
		errors.Is(err, services.ErrClassFormat),
		errors.Is(err, services.ErrSubjectFormat),
		errors.Is(err, services.ErrValidDOB),
		errors.Is(err, services.ErrScore),
		errors.Is(err, services.ErrMaxScore),
		errors.Is(err, services.ErrIDRequired),
		errors.Is(err, services.ErrNameRequired),
		errors.Is(err, services.ErrStudentClass),
		errors.Is(err, services.ErrInvalidYear),
		errors.Is(err, services.ErrInvalidGender),
		errors.Is(err, services.ErrAddressTooShort),
		errors.Is(err, services.ErrInvalidMinMax),
		errors.Is(err, services.ErrStudentRank):
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

// GET STUDENTS
func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.service.GetAllStudents()

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, students)
}

func (h *StudentHandler) GetStudentByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	student, err := h.service.GetStudentByEmail(email)
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

// GET STUDENTS IF HAVE EMAIL QUERY PARAM
func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email != "" {
		h.GetStudentByEmail(w, r)
		return
	}

	//CHECK FILTER PARAM
	q := r.URL.Query()

	if q.Get("name") != "" || q.Get("class") != "" || q.Get("gender") != "" || q.Get("year_of_birth") != "" ||
		q.Get("address") != "" || q.Get("min_score") != "" || q.Get("max_score") != "" || q.Get("rank") != "" {
		h.FilterStudents(w, r)
		return
	}

	h.GetAllStudents(w, r)
}

func (h *StudentHandler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
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

// POST STUDENT
func (h *StudentHandler) AddStudent(w http.ResponseWriter, r *http.Request) {
	student := models.Student{}
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.service.AddStudent(&student); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, &student)
}

// PUT STUDENT
func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	student := models.Student{}
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()
	student.ID = id

	if err := h.service.UpdateStudent(&student); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, &student)
}

// DELETE STUDENT
func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteStudent(id); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "student delete successfully"})
}

/* -----------SCORES--------------- */
func (h *StudentHandler) GetScoresByStudentID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	scores, err := h.service.GetScoresByStudentID(id)

	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, scores)
}

func (h *StudentHandler) GetScoresBySubject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	subject := r.PathValue("subject")

	score, err := h.service.GetScoresBySubject(id, subject)

	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, score)

}

func (h *StudentHandler) AddScore(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	score := models.SubjectScore{}
	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.service.AddScore(id, &score); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, &score)
}

func (h *StudentHandler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	subject := r.PathValue("subject")

	score := models.SubjectScore{}
	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	defer r.Body.Close()
	score.Subject = subject

	if err := h.service.UpdateScore(id, &score); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, &score)
}

func (h *StudentHandler) DeleteScore(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	subject := r.PathValue("subject")

	if err := h.service.DeleteScore(id, subject); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "scores deleted successfully"})
}

func (h *StudentHandler) FilterStudents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := &models.FilterStudents{
		Name:    q.Get("name"),
		Class:   q.Get("class"),
		Gender:  q.Get("gender"),
		Address: q.Get("address"),
	}

	if v := q.Get("year_of_birth"); v != "" {
		year, err := strconv.Atoi(v)

		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid year_of_birth")
			return
		}
		filter.YearOfBirth = year
	}

	if v := q.Get("min_score"); v != "" {
		minScore, err := strconv.ParseFloat(v, 64)

		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid min_score")
			return
		}
		filter.MinAvgScore = minScore
	}

	if v := q.Get("max_score"); v != "" {
		maxScore, err := strconv.ParseFloat(v, 64)

		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid max_score")
			return
		}
		filter.MaxAvgScore = maxScore
	}

	if v := q.Get("rank"); v != "" {
		filter.StudentRank = models.Rank(v)
	}

	students, err := h.service.FilterStudents(filter)
	if err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, students)

}

/* -----------BULK UPLOAD CSV--------------- */
func (h *StudentHandler) BulkAddStudents(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with max 10MB size
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	// Get CSV file from form
	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file not found in form")
		return
	}
	defer file.Close()

	// Read CSV
	_, rows, err := utils.ReadCSV(file)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read CSV: "+err.Error())
		return
	}

	if len(rows) == 0 {
		writeError(w, http.StatusBadRequest, "CSV is empty")
		return
	}

	// Group rows by student ID
	studentRows := make(map[string][][]string)
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		studentID := row[0]
		if studentID == "" {
			writeError(w, http.StatusBadRequest, "student ID cannot be empty")
			return
		}
		studentRows[studentID] = append(studentRows[studentID], row)
	}

	// Convert rows to Student objects
	students := make([]*models.Student, 0, len(studentRows))
	for _, rowGroup := range studentRows {
		student, err := models.StudentFromCSVRows(rowGroup)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse student from CSV: %v", err))
			return
		}
		students = append(students, student)
	}

	// Call service to bulk add students
	if err := h.service.BulkAddStudents(students); err != nil {
		writeError(w, serviceErrToStatus(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": fmt.Sprintf("successfully added %d students", len(students)),
		"count":   len(students),
	})
}

/* -----------CSV EXPORT/DOWNLOAD--------------- */
func (h *StudentHandler) ExportStudents(w http.ResponseWriter, r *http.Request) {
	// Check for filter parameters
	q := r.URL.Query()

	var students []*models.Student
	var err error

	// If filter parameters exist, use filtered results
	if q.Get("name") != "" || q.Get("class") != "" || q.Get("gender") != "" ||
		q.Get("year_of_birth") != "" || q.Get("address") != "" ||
		q.Get("min_score") != "" || q.Get("max_score") != "" || q.Get("rank") != "" {

		filter := &models.FilterStudents{
			Name:    q.Get("name"),
			Class:   q.Get("class"),
			Gender:  q.Get("gender"),
			Address: q.Get("address"),
		}

		if v := q.Get("year_of_birth"); v != "" {
			year, err := strconv.Atoi(v)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid year_of_birth")
				return
			}
			filter.YearOfBirth = year
		}

		if v := q.Get("min_score"); v != "" {
			minScore, err := strconv.ParseFloat(v, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid min_score")
				return
			}
			filter.MinAvgScore = minScore
		}

		if v := q.Get("max_score"); v != "" {
			maxScore, err := strconv.ParseFloat(v, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid max_score")
				return
			}
			filter.MaxAvgScore = maxScore
		}

		if v := q.Get("rank"); v != "" {
			filter.StudentRank = models.Rank(v)
		}

		students, err = h.service.FilterStudents(filter)
	} else {
		// Get all students
		students, err = h.service.GetAllStudents()
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(students) == 0 {
		writeError(w, http.StatusNotFound, "no students found to export")
		return
	}

	// Convert students to CSV rows
	csvRows := [][]string{}
	for _, student := range students {
		csvRows = append(csvRows, student.ToCSVRows()...)
	}

	// Set response headers for file download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"students.csv\"")
	w.WriteHeader(http.StatusOK)

	// Write CSV to response
	if err := utils.WriteCSV(w, models.CSVHeader(), csvRows); err != nil {
		// Headers already sent, log the error instead
		fmt.Fprintf(w, "\nerror writing CSV: %v", err)
		return
	}
}
