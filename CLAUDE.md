# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.

---

## Project Context

- **App:** Student Management System — manages student and teacher data (CRUD, filtering)
- **Backend:** Go
- **Frontend:** React + Tailwind CSS + Axios
- **Entry point:** `cmd/app/main.go`
- **Config:** `config/config.go`

## Structure Convention

### Backend
- Features are split by domain: `student/` and `teacher/` under each layer
- Layers: `handlers` → `services` → `repositories` → `models`
- Handler packages: `studentHandler` and `teacherHandler`
- Filter/predicate logic lives in `internal/predicates/`
- Views live in `internal/views/`
- Shared utilities in `pkg/utils/`

### Frontend (planned)
- React component structure should mirror backend domains: `student/` and `teacher/`
- Tailwind CSS for styling — no separate CSS files unless necessary
- Axios for all API calls to the Go backend
- Keep components small and focused — one responsibility per component

## Current State
- Repository layer uses in-memory storage (`*_repo_memory.go`)
- All API endpoints are being tested and verified via Postman before any migration
- Keep repository interfaces (`student_repo.go`, `teacher_repo.go`) clean for easy swapping later
- Frontend not started yet

## Roadmap
1. ✅ Build and structure the Go backend
2. 🔄 Test and verify all API endpoints via Postman (current phase)
3. ⏳ Migrate repository layer to MySQL
4. ⏳ Build frontend with React + Tailwind CSS + Axios
5. ⏳ Deploy and bring the project live

## Patterns to Follow
- Match existing domain-split folder structure when adding new features
- New filters/predicates go in `internal/predicates/`
- Don't introduce new dependencies without flagging first
- Don't jump ahead of the current phase — changes should align with the active roadmap step
- When building frontend, mirror the backend domain structure in React component folders