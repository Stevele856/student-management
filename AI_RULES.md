# AI Rules

## Assumptions About Current State

- This is a small layered Go HTTP service using `net/http`, `http.ServeMux`, service structs, repository interfaces, in-memory JSON-backed repositories, model packages, predicate-based filtering, and utility validators.
- The student feature is the only fully routed API path. The teacher feature has models, service, repository, and predicate scaffolding, but no HTTP handlers/router and `TeacherService.FilterTeachers` is unfinished.
- `go test ./...` currently fails because `internal/services/teacher/teacher_services.go` has a missing return in `FilterTeachers`; student service tests and predicate tests pass.

## 1. Architecture Rules

- Keep the existing layered structure:
  - `cmd/app`: application wiring only.
  - `config`: environment/default configuration only.
  - `internal/handlers`: HTTP routing, request parsing, and response writing.
  - `internal/services/<entity>`: normalization, validation, business rules, UUID generation, and repository orchestration.
  - `internal/repositories/<entity>`: repository interfaces and JSON-backed in-memory implementations.
  - `internal/models/<entity>`: data structs, JSON tags, enum-like types, and CSV conversion for the entity.
  - `internal/predicate`: reusable filter predicates.
  - `pkg/utils`: shared format validators, score/rank helpers, and CSV helpers.
- Do not introduce a new framework, ORM, web router, dependency injection container, database driver, or generic architecture layer unless the existing pattern cannot support the required change.
- Keep `cmd/app/main.go` as the composition root: load config, create repository, create service, create handler/router, start `http.Server`.
- New entity features must follow the existing student/teacher package shape before adding shortcuts elsewhere.
- Do not put business rules in handlers or repositories.

## 2. Layer Responsibilities

- Handlers must:
  - Read path values with `r.PathValue`.
  - Read query parameters from `r.URL.Query()`.
  - Decode request JSON with `json.NewDecoder(r.Body).Decode`.
  - Close request bodies after successful decoding.
  - Convert primitive query strings to typed filter fields.
  - Call exactly one service method for the operation.
  - Convert service errors to HTTP responses.
- Services must:
  - Trim and normalize input before validation.
  - Own all required-field checks, format validation, uniqueness checks, duplicate checks, score limits, rank validation, and ID generation.
  - Use repository interfaces, never concrete repository implementations.
  - Return package-level sentinel errors for expected business and validation failures.
- Repositories must:
  - Own map storage and JSON file load/save only.
  - Persist after mutating operations.
  - Return empty slices and `nil` error for empty list queries.
  - Avoid HTTP concepts, request parsing, UUID generation, and validation.
- Models must:
  - Define structs, JSON field names, enum-like constants, and entity-local serialization helpers such as CSV row conversion.
  - Avoid service or repository imports.
- Predicates must:
  - Contain pure filter functions only.
  - Avoid persistence, validation, HTTP, and mutation.

## 3. Naming Conventions

- Keep package folder names singular and entity-based: `student`, `teacher`.
- Use import aliases when model and service/repository package names collide:
  - `studentModels`, `teacherModels`
  - `studentRepo`, `teacherRepo`
- Keep constructor names in the current form: `NewStudentService`, `NewTeacherService`, `NewStudentMemoryRepo`, `NewTeacherMemoryRepo`, `NewStudentHandler`, `NewRouter`.
- Repository interfaces must be named `<Entity>Repository`.
- In-memory JSON repositories must stay named `InMemo<Entity>Repo` to match existing code.
- Service methods must use the current verbs: `Add<Entity>`, `Update<Entity>`, `Delete<Entity>`, `Get<Entity>By...`, `GetAll<Entity>`, `Filter<Entities>`.
- JSON tags and query parameters must remain snake_case, for example `full_name`, `date_of_birth`, `year_of_birth`, `min_score`, `max_score`, `student_rank`.
- Error variables must be package-level `Err...` values.

## 4. Error Handling Standards

- Expected validation and business failures must be returned as sentinel errors from services.
- Handlers must use `writeError(w, status, err.Error())` and preserve the existing JSON error shape: `{"error":"..."}`.
- Student HTTP error mapping must go through `serviceErrToStatus`.
- Add any new student service error to `serviceErrToStatus` in the same change that introduces it.
- Not-found, duplicate, invalid-format, and required-field cases must not fall through to HTTP 500.
- Repository errors must not be exposed directly from handlers unless the service intentionally passes them through.
- Teacher repository not-found and pagination errors already use sentinel errors in `internal/repositories/teacher/errors.go`; keep using `errors.Is` for those.
- When wrapping technical persistence errors, use `%w` so callers can inspect them.

## 5. Validation Rules

- Use `pkg/utils` validators for shared formats:
  - `IsValidEmail`
  - `IsValidName`
  - `IsValidClass`
  - `IsValidSubject`
  - `IsValidSubjectScore`
  - `IsValidEmployeeID`
  - `IsValidPhoneNumber`
- Student normalization must trim names, class, address, and lowercase email.
- Teacher normalization must trim names, address, phone, subject/class values, lowercase email/status, and uppercase employee ID.
- Scores must stay in the `0..10` range and students must not exceed 10 subject scores.
- Student ranks must remain `excellent`, `good`, `average`, or `weak`.
- Gender must remain `male` or `female` using `internal/models.Gender` where the model already uses it.
- CSV dates must use `studentModels.DateLayout` (`2006-01-02`).
- Service normalizers should copy structs rather than mutating caller-owned input unless existing behavior for that method requires mutation.

## 6. Repository & Database Rules

- The current database is the JSON file configured by `DATA_FILE`, defaulting to `static/students.json`.
- Do not add SQL, migrations, ORM models, or external storage without a clear migration plan and explicit user request.
- In-memory repositories must load JSON in their constructor and save JSON after add/update/delete/bulk-add operations.
- Student repository keys are student IDs stored in `map[string]*studentModels.Student`.
- Teacher repository keys are teacher IDs stored in `map[string]*teacherModels.Teacher`.
- Email and subject lookups must use case-insensitive comparisons where the existing code does.
- List results from map iteration are unordered unless explicitly sorted; pagination must sort by stable ID as teacher pagination already does.
- Bulk student CSV import must validate every student and detect duplicate emails before writing any student.

## 7. API Response Format

- Successful JSON responses must use `writeJSON`.
- Error responses must use `writeError` and the exact `error` key.
- Delete success responses currently use a JSON object with a `message` key; keep that shape.
- Bulk upload success responses currently use `message` and `count`; keep that shape.
- CSV export must set `Content-Type: text/csv` and `Content-Disposition: attachment; filename="students.csv"`.
- Do not introduce a global response envelope unless all existing handlers are intentionally migrated together.
- Routes must stay in `internal/handlers/routers.go` using Go 1.22-style method patterns such as `GET /students/{id}`.

## 8. Testing Expectations

- Use table-driven tests for validation and service behavior, following `internal/services/student/student_services_test.go`.
- Use package-local tests for service internals when testing unexported helpers.
- Use `predicate_test` external package style for predicate behavior.
- Mock repositories with function fields, following `MockStudentRepository`.
- Every new service rule must include at least one success case and one failure case.
- Every new predicate must have direct predicate tests and at least one composition test when used with `predicate.And` or `predicate.AndTeacher`.
- Before considering work done, run `go test ./...`.
- If `go test ./...` fails because of pre-existing teacher incompleteness, state that explicitly and do not hide new failures.

## 9. Scope Control Rules For AI Behavior

- Make the smallest change that fits the current architecture.
- Do not rename packages, exported types, JSON fields, routes, or environment variables unless the user asked for a breaking change.
- Do not rewrite the student API while working on teacher scaffolding.
- Do not change persisted JSON shape casually; existing `static/students.json` data must remain readable.
- Do not modify generated/exported data files such as `exported_students.csv` unless the task is specifically about CSV output.
- Do not touch unrelated dirty files. This repository may contain user edits.
- If a feature is incomplete, finish it inside the existing layer pattern before adding new surfaces.

## 10. Anti-patterns

- Do not put validation only in handlers; services must enforce it.
- Do not let handlers inspect repository errors directly.
- Do not add service methods that depend on concrete `InMemo...Repo` types.
- Do not add filters that parse query strings in repositories or predicates.
- Do not mutate input structs in normalization helpers when tests expect copied values.
- Do not return `nil, nil` for missing singular records unless the service explicitly handles that case.
- Do not return HTTP 500 for user input errors.
- Do not add empty placeholder files or functions; unfinished code currently breaks the teacher package.
- Do not save repository data twice for one mutation.
- Do not rely on Go map iteration order for API responses that need stable ordering.
- Do not duplicate filter-building logic between normal JSON endpoints and CSV export if the logic changes; keep them consistent.

## 11. Definition of Done

- The change fits the existing folder and layer responsibilities.
- New validation lives in the service layer and uses existing `pkg/utils` helpers when applicable.
- New errors are sentinel `Err...` values and are mapped to HTTP status codes when exposed through handlers.
- Repository changes preserve JSON-backed in-memory persistence.
- API responses preserve the current success/error shapes.
- Tests cover changed service rules, predicates, and repository behavior when repository behavior changes.
- `gofmt` has been run on changed Go files.
- `go test ./...` has been run, and any remaining failure is documented with the package and reason.
- No unrelated user edits or data files were reverted or reformatted.
