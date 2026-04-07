# FORGE — Example: .NET REST API

> This is a filled example of the FORGE template applied to a .NET 8 minimal API project.

---

## LAYER 1 — CONTEXT

### Project
A REST API for a task management system. Users can create, list, update, and delete tasks. The API is consumed by a React frontend (not part of this scope).

### Tech Stack
- .NET 8 Minimal API
- Entity Framework Core 8 with PostgreSQL
- FluentValidation for input validation
- xUnit + Testcontainers for tests
- Docker Compose for local development

### Current State
- Project scaffolded, runs on `dotnet run`
- Two endpoints exist: `GET /tasks` and `POST /tasks`
- No authentication
- No input validation on `POST /tasks`
- Tests directory exists but is empty

### Integrations
- PostgreSQL via connection string in `appsettings.json`
- No external services

### Entry Points
- `dotnet run` → starts on `http://localhost:5000`
- `docker compose up` → starts with PostgreSQL

---

**[ GATE 1 ]** ✅ — Context validated by user.

---

## LAYER 2 — PROBLEM

### Priority 1
`POST /tasks` accepts any input including empty strings and null values, which causes a database exception instead of a proper validation error.

### Priority 2
There are no tests. Changes cannot be verified without manual testing.

---

**[ GATE 2 ]** ✅ — Problem list confirmed. P1 is blocking, P2 follows.

---

## LAYER 3 — LOCKED PATH

### [DECISION] — Technical decisions already made
```
[DECISION] FluentValidation is already installed. Use it for all validation. Do not implement manual validation.
Reason: Already a dependency, team knows it.
Scope: Project-wide.
```

### [THOUGHT] — Conceptual directions that are off-limits
```
[THOUGHT] Do not add authentication in this session.
Reason: Auth is a separate workstream and will be handled later.
Implication: All endpoints remain open. Do not add any auth middleware or guards.
```

---

**[ GATE 3 ]** ✅ — Locked paths confirmed.

---

## LAYER 4 — PLANNING

### Task List

| # | Task | Depends On | Skills |
|---|---|---|---|
| T1 | Add FluentValidation validator for CreateTask request | — | execution/task-runner |
| T2 | Register validator and return 400 on invalid input | T1 | execution/task-runner, execution/scope-guard |
| T3 | Write unit tests for the validator | T1 | testing/unit |
| T4 | Write integration test for POST /tasks with invalid input | T2 | testing/acceptance |

### Skills Required

| Skill | Status |
|---|---|
| skills/execution/task-runner.md | existing |
| skills/execution/scope-guard.md | existing |
| skills/testing/unit.md | existing |
| skills/testing/acceptance.md | existing |

---

**[ GATE 4 ]** ✅ — Plan approved.

---

## LAYER 5 — EXECUTION

### T1 — Add FluentValidation validator for CreateTask request
- Status: `done`
- Output: `CreateTaskRequestValidator.cs` created in `/Validators`
- Notes: Validates that Title is not empty, MaxLength 200. Description is optional.

### T2 — Register validator and return 400 on invalid input
- Status: `done`
- Output: Validator registered in `Program.cs`, endpoint returns `400 + validation errors` on invalid input
- Notes: Used `Results.ValidationProblem()` for RFC 7807 compliance

### T3 — Unit tests for the validator
- Status: `done`
- Output: `CreateTaskRequestValidatorTests.cs` — 4 tests, all passing

### T4 — Integration test for POST /tasks with invalid input
- Status: `done`
- Output: `TaskEndpointsTests.cs` — 2 acceptance tests, all passing

---

## LAYER 6 — TESTING

### T1 / T3 Tests
- [x] Unit: empty title → validation fails
- [x] Unit: title too long (201 chars) → validation fails
- [x] Unit: valid title + no description → validation passes
- [x] Unit: valid title + valid description → validation passes

### T2 / T4 Tests
- [x] Acceptance: POST /tasks with empty body → 400 + error details
- [x] Acceptance: POST /tasks with valid body → 201 + created task

---

## LAYER 7 — DOCUMENTATION

### What Was Built
Input validation for `POST /tasks` using FluentValidation. The endpoint now returns `400 Bad Request` with structured validation errors (RFC 7807 format) when the request is invalid. All validation logic is isolated in `CreateTaskRequestValidator`.

### How to Use
```
POST /tasks
Content-Type: application/json

{ "title": "", "description": "..." }

→ 400 Bad Request
{
  "type": "https://tools.ietf.org/html/rfc7807",
  "errors": {
    "Title": ["'Title' must not be empty."]
  }
}
```

### Dependencies
- FluentValidation.AspNetCore (already installed)
- No new dependencies added

### Known Limitations
- Validation messages are in English — no localization
- Description field has no max length constraint (deferred)
