# Project Status — Student Management System

**Last updated:** 2026-06-09  
**Current phase:** Postman testing (Phase 2) — Phase 1 bug fixes complete

---

## What Is Complete

### Architecture & Structure
- [x] Layered architecture: Handlers → Services → Repositories → Models
- [x] Domain split: `student/` and `teacher/` across all layers
- [x] Repository interface defined for both domains (ready for DB swap)
- [x] In-memory repository with JSON file persistence
- [x] Config loader via environment variables (`PORT`, `STUDENT_DATA_FILE`, `TEACHER_DATA_FILE`)

### Student Domain
- [x] Model: `Student` struct with `SubjectScore`, `Rank` enum, CSV serialization
- [x] Filter model: `FilterStudents` with all fields
- [x] Repository interface + in-memory implementation (CRUD, scores, pagination, bulk, filter)
- [x] Service: AddStudent, UpdateStudent, DeleteStudent, GetAllStudents, GetStudentByID, GetStudentByEmail
- [x] Service: AddScore, UpdateScore, DeleteScore, GetScoresByStudentID, GetScoresBySubject
- [x] Service: FilterStudents (partial — see bugs), GetStudentsPaginated, BulkAddStudents
- [x] Predicates: ByName, ByClass, ByYear, ByGender, ByAddress, ByAvgScore, ByRank
- [x] Handler + Router for all student endpoints
- [x] CSV import (bulk-upload) and export

### Teacher Domain
- [x] Model: `Teacher` struct with status enum, CSV serialization
- [x] Filter model: `FilterTeachers` with all fields
- [x] Repository interface + in-memory implementation (CRUD, pagination, bulk, filter)
- [x] Service: AddTeacher, UpdateTeacher, DeleteTeacher, GetAllTeachers, GetTeacherByID
- [x] Service: GetTeacherByEmail, GetTeacherByEmployeeID, GetTeacherByStatus, GetTeacherAssignedBySubject, GetTeacherByAssignedClass
- [x] Service: FilterTeachers, GetTeachersPaginated, BulkAddTeachers
- [x] Predicates: ByTeacherName, ByTeacherGender, ByTeacherEmail, ByEmployeeID, ByStatus, ByClassAssigned, BySubjectTaught, ByHireDateRange
- [x] Handler + Router for all teacher endpoints
- [x] CSV import (bulk-upload) and export

### Utilities
- [x] Validation: IsValidEmail, IsValidName, IsValidClass, IsValidSubject, IsValidEmployeeID, IsValidPhoneNumber, IsValidSubjectScore, IsValidScores
- [x] Calculation: CalAvgScore, CalcStudentRankBaseOnAvgScore
- [x] CSV helpers: ReadCSV, WriteCSV

### Tests
- [x] Student predicate tests: ByName, ByYear, ByAvgScore, ByRank, AndStudent
- [x] Student service tests: validateStudent (partial cases)
- [x] Mock StudentRepository for unit tests

---

## Bugs Found (Code Audit — ✅ ALL FIXED 2026-06-09)

These were confirmed code bugs found by reading the source. **All have been fixed and verified** (`go build ./...` and `go test ./...` pass). Kept here for record.

### BUG-01 — `year_of_birth` filter is silently ignored
**File:** `internal/services/student/student_services.go` — `filterPredicateStudents()`  
**Problem:** The `ByYear` predicate exists but is never added to the predicate chain. Filtering by `year_of_birth` has zero effect — returns all students regardless.  
**Fix:** Add `if filter.YearOfBirth != 0 { predicate = append(predicate, predicates.ByYear(filter.YearOfBirth)) }`

### BUG-02 — `min_score` alone or `max_score` alone does nothing
**File:** `internal/services/student/student_services.go` — `filterPredicateStudents()`  
**Problem:** Condition uses `&&` — score filter only fires when BOTH min AND max are provided. Sending only `min_score=7` returns all students.  
**Fix:** Change `filter.MinAvgScore > 0 && filter.MaxAvgScore > 0` to `filter.MinAvgScore > 0 || filter.MaxAvgScore > 0`

### BUG-03 — `UpdateStudent` returns 500 when changing to a new email
**File:** `internal/services/student/student_service_helper.go` — `ensureStudentUnique()`  
**Problem:** Compares against service-level `ErrStudentNotFound` but repo returns `studentRepo.ErrStudentNotFound`. These are different error values. When a new email is unique (nobody has it), the repo returns "not found" which is misread as an unexpected error → function returns the repo error → handler maps it to 500.  
**Fix:** Change `!errors.Is(err, ErrStudentNotFound)` to `!errors.Is(err, studentRepo.ErrStudentNotFound)`

### BUG-04 — `POST /students` and `POST /teachers` response missing generated UUID
**File:** `internal/handlers/student/student_handler.go` `AddStudent()`, `internal/handlers/teacher/teacher_handler.go` `AddTeacher()`  
**Problem:** The service calls `normalizeStudent(student)` which creates a copy (new pointer). The ID and timestamp are set on the copy, not on the original struct the handler holds. The response body shows the original un-modified struct — no ID, no `CreatedAt`.  
**Fix:** Service should mutate the caller's struct directly instead of working on a copy, OR the handler should call a `GetStudentByID` after add to fetch the full stored record.

### BUG-05 — `PUT /students/{id}` and `PUT /teachers/{id}` response missing `UpdatedAt`
**File:** Same normalization issue as BUG-04. The `student = normalizeStudent(student)` in the service creates a copy; `UpdatedAt = time.Now()` is set on the copy. Handler returns the original stale struct.

### BUG-06 — `ErrStudentIDRequired` and `ErrSubjectEmpty` return 404 instead of 400
**File:** `internal/handlers/student/student_handler.go` — `serviceErrToStatus()`  
**Problem:** Both errors are grouped under the 404 `case` but they are validation errors (missing input), not "not found" errors. They should return 400 Bad Request.

### BUG-07 — `ErrTeacherIDRequired` returns 404 instead of 400
**File:** `internal/handlers/teacher/teacher_handler.go` — `serviceErrToStatus()`  
**Same problem as BUG-06 for the teacher domain.**

### BUG-08 — `ErrScoreRange` not mapped → returns 500 instead of 400
**File:** `internal/handlers/student/student_handler.go` — `serviceErrToStatus()`  
**Problem:** `ErrScoreRange` (returned when `min_score` or `max_score` is outside 0–10) is defined in the service but missing from the handler's error map. Result: sending `?min_score=-1` returns 500.

### BUG-09 — Address filter uses exact match instead of substring
**File:** `internal/predicates/student_filters.go` — `ByAddress()`  
**Problem:** Uses `strings.EqualFold(s.Address, address)` (exact match). A user searching `?address=Hanoi` won't find students whose address is "123 Hanoi Street". Should be `strings.Contains(toLower(s.Address), toLower(address))` to match ByName behavior.

### BUG-10 — Class not uppercased in student normalization
**File:** `internal/services/student/student_service_helper.go` — `normalizeStudent()`  
**Problem:** `s.Class = strings.TrimSpace(s.Class)` — no `ToUpper`. `IsValidClass` requires uppercase letter (e.g., "10A"). Submitting `"class": "10a"` fails with `ErrClassFormat`. Teacher service correctly uppercases class but student service does not.

### BUG-11 — `BulkAddTeachers` uses wrong error for duplicate ID in payload
**File:** `internal/services/teacher/teacher_services_helper.go` line ~322  
**Problem:** Uses `ErrTeacherIDRequired` ("teacher ID required") when a teacher ID appears twice in the same CSV upload. Error message is misleading — should be a duplicate ID error.  
**Note:** Even the code has a comment flagging this.

---

## Endpoints to Test (Postman Checklist)

Work through these in order. Fix bugs as you find them. Mark [x] when verified.

### Student Endpoints

**Basic CRUD**
- [x] `POST /students` — create with valid body → expect 201 + body with ID
- [x] `POST /students` — missing name → expect 400
- [x] `POST /students` — invalid email → expect 400
- [x] `POST /students` — duplicate email → expect 409
- [x] `POST /students` — future date_of_birth → expect 400
- [x] `POST /students` — invalid class (e.g., "10a" lowercase) → expect 400 *(BUG-10 must be fixed first)*
- [x] `GET /students/{id}` — existing ID → expect 200
- [x] `GET /students/{id}` — non-existent ID → expect 404
- [x] `PUT /students/{id}` — valid update → expect 200 + body with UpdatedAt *(BUG-04/05 must be fixed)*
- [x] `PUT /students/{id}` — change email to a new unique email → expect 200 *(BUG-03 must be fixed)*
- [x] `PUT /students/{id}` — change email to an email used by another student → expect 409
- [x] `DELETE /students/{id}` — existing → expect 200
- [x] `DELETE /students/{id}` — non-existent → expect 404

**Listing & Pagination**
- [x] `GET /students` — no params → expect 200 + full list
- [x] `GET /students?page=1&page_size=5` → expect 200 + paginated response with `total`
- [x] `GET /students?page=0` → expect 400
- [x] `GET /students/paginated?page=1&page_size=5` → same as above

**Filters** *(Fix BUG-01, BUG-02 first)*
- [x] `GET /students?name=Nguyen` → expect matching students
- [x] `GET /students?class=10A` → expect matching students
- [x] `GET /students?year_of_birth=2005` → expect matching *(BUG-01)*
- [x] `GET /students?gender=male` → expect matching
- [x] `GET /students?address=Hanoi` → expect substring match *(BUG-09)*
- [x] `GET /students?min_score=7` → expect students with avg >= 7 *(BUG-02)*
- [x] `GET /students?max_score=5` → expect students with avg <= 5 *(BUG-02)*
- [x] `GET /students?min_score=5&max_score=8` → expect students with 5 <= avg <= 8
- [x] `GET /students?min_score=-1` → expect 400 *(BUG-08)*
- [x] `GET /students?rank=excellent` → expect students with avg >= 9
- [x] `GET /students?name=Nguyen&class=10A` → combined filter (AND semantics)
- [x] `GET /students?name=Nguyen&page=1` → should paginate, ignores name (current design) — note behavior

**Scores**
- [x] `POST /students/{id}/scores` — add new subject → expect 201
- [x] `POST /students/{id}/scores` — duplicate subject → expect 409
- [x] `POST /students/{id}/scores` — score > 10 → expect 400
- [x] `POST /students/{id}/scores` — 11th score → expect 400
- [x] `GET /students/{id}/scores` → expect 200 + array
- [x] `GET /students/{id}/scores/{subject}` → expect 200 + score object
- [x] `GET /students/{id}/scores/{subject}` — non-existent subject → expect 404
- [x] `PUT /students/{id}/scores/{subject}` — update score → expect 200
- [x] `DELETE /students/{id}/scores/{subject}` → expect 200
- [x] `DELETE /students/{id}/scores/{subject}` — non-existent → expect 404

**CSV**
- [x] `POST /students/bulk-upload` — valid CSV → expect 201 + count
- [x] `POST /students/bulk-upload` — duplicate email in CSV → expect 409
- [x] `GET /students/export` — no params → download all as CSV
- [x] `GET /students/export?class=10A` → filtered CSV download

---

### Teacher Endpoints

**Basic CRUD**
- [x] `POST /teachers` — valid body → expect 201 + body with ID *(BUG-04)*
- [x] `POST /teachers` — missing phone → expect 400
- [x] `POST /teachers` — invalid employee_id format (not T###) → expect 400
- [x] `POST /teachers` — hire date before min age 25 → expect 400
- [x] `POST /teachers` — duplicate email → expect 409
- [x] `POST /teachers` — duplicate employee_id → expect 409
- [x] `GET /teachers/{id}` — existing → expect 200
- [x] `GET /teachers/{id}` — non-existent → expect 404
- [x] `GET /teachers/employee/{employee_id}` — e.g., `T001` → expect 200
- [x] `GET /teachers/employee/{employee_id}` — non-existent → expect 404
- [x] `PUT /teachers/{id}` — valid update → expect 200 *(BUG-05)*
- [x] `DELETE /teachers/{id}` — existing → expect 200
- [x] `DELETE /teachers/{id}` — non-existent → expect 404

**Listing & Pagination**
- [x] `GET /teachers` — no params → expect 200 + full list
- [x] `GET /teachers?page=1&page_size=5` → paginated
- [x] `GET /teachers/paginated?page=1&page_size=5` → same
- [x] `GET /teachers?page=0` → expect 400

**Filters**
- [x] `GET /teachers?name=Nguyen` → expect matching teachers
- [x] `GET /teachers?status=active` → expect active teachers
- [x] `GET /teachers?status=on-leave` → expect on-leave teachers
- [x] `GET /teachers?status=invalid` → expect 400
- [x] `GET /teachers?gender=female` → expect female teachers
- [x] `GET /teachers?employee_id=T001` → expect single match
- [x] `GET /teachers?subject=Math` → expect teachers teaching Math
- [x] `GET /teachers?class=10A` → expect teachers assigned to 10A
- [x] `GET /teachers?subject=Math&subject=English` → AND semantics (must teach both)
- [x] `GET /teachers?hire_date_from=2020-01-01&hire_date_to=2023-12-31` → date range
- [x] `GET /teachers?hire_date_from=2023-01-01&hire_date_to=2020-01-01` → expect 400 (from > to)
- [x] `GET /teachers?hire_date_from=invalid` → expect 400
- [x] `GET /teachers?name=Nguyen&status=active` → combined filter

**CSV**
- [ ] `POST /teachers/bulk-upload` — valid CSV → expect 201 + count
- [ ] `POST /teachers/bulk-upload` — duplicate employee_id in CSV → expect 409
- [ ] `GET /teachers/export` — no params → download all
- [ ] `GET /teachers/export?status=active` → filtered download

---

## Tasks To Do Next (In Order)

### Phase 1 — Fix Code Bugs ✅ COMPLETE (2026-06-09)

All bugs fixed and verified (`go build ./...` and `go test ./...` pass).

1. [x] **BUG-01** — `ByYear` predicate added to `filterPredicateStudents`
2. [x] **BUG-02** — Score range condition uses `||`
3. [x] **BUG-03** — `ensureStudentUnique` uses `studentRepo.ErrStudentNotFound`
4. [x] **BUG-04 & BUG-05** — Service mutates caller's struct; ID/timestamps now in POST/PUT response
5. [x] **BUG-06** — `ErrStudentIDRequired` and `ErrSubjectEmpty` mapped to 400
6. [x] **BUG-07** — `ErrTeacherIDRequired` mapped to 400
7. [x] **BUG-08** — `ErrScoreRange` mapped to 400 in student `serviceErrToStatus`
8. [x] **BUG-09** — `ByAddress` uses substring (`strings.Contains`)
9. [x] **BUG-10** — `normalizeStudent` uppercases class
10. [x] **BUG-11** — `BulkAddTeachers` returns `ErrTeacherIDDuplicated` for duplicate ID

### Phase 2 — Complete Postman Testing

- Work through every case in the checklist above
- Log any new issues discovered and add to this file
- Verify all error codes match expectations

### Phase 3 — Missing Tests

- [ ] Write teacher predicate tests (`teacher_filter_test.go` is empty)
- [ ] Write teacher service tests (no tests exist for teacher service)
- [ ] Expand student service tests (currently only `validateStudent` is covered)

### Phase 4 — Before MySQL Migration (Improvements Needed)

- [ ] **Add CORS middleware** — required for React frontend; without it, all browser calls will be blocked
- [ ] **Consistent error wrapping in teacher service** — `wrapRepoError` in teacher service passes repo errors through unchanged (inconsistent with student service pattern; repo errors leak out to handler)
- [ ] **Validate `students.json` file path on startup** — if the JSON path is wrong, only a fatal crash reveals it
- [ ] **Add `GET /students/paginated` route awareness to dispatcher** — `GET /students?name=X&page=1` currently ignores the name filter (pagination check fires first); decide if this is intentional and document it
- [ ] **`PasswordHash` has no set/change API** — currently only preserved on update; no endpoint to set it; decide before frontend

### Phase 5 — MySQL Migration (Roadmap Step 3)

- Implement `StudentRepository` interface against a real DB
- Implement `TeacherRepository` interface against a real DB
- Wire up via config/env (swap `NewStudentMemoryRepo` → `NewStudentMySQLRepo`)
- Run all Postman tests again against the DB-backed implementation

### Phase 6 — Frontend (Roadmap Step 4)

- Start React app with Tailwind CSS + Axios
- Mirror domain structure: `src/student/`, `src/teacher/`
- Must have CORS middleware working first (Phase 4)
