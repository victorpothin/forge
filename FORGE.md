# FORGE — Project Template

> Copy this file into your project. Fill each layer in order. Do not advance to the next layer without validating the current one.

---

## LAYER 1 — CONTEXT

> **AI instruction:** Read all project files, existing code, and any materials the user provided. Synthesize your understanding below. Do not infer intent — only describe what exists. Then stop and wait for validation.

### Project
<!-- AI fills this -->

### Tech Stack
<!-- AI fills this -->

### Current State
<!-- AI fills this: what exists, what works, what doesn't -->

### Integrations and Dependencies
<!-- AI fills this -->

---

**[ GATE 1 ]** — User validates context before proceeding. If anything is wrong, correct it here before moving on.

---

## LAYER 2 — PROBLEM

> **AI instruction:** Based on the validated context, list the problems the user needs to solve. Do NOT suggest solutions here. List only problems, ordered by priority. If there are no problems, state that explicitly — nothing will be executed.

### Priority 1
<!-- Most critical problem -->

### Priority 2
<!-- Second problem — only tackle if Priority 1 is resolved -->

### Priority N
<!-- Continue as needed -->

> **Rule:** If this section is empty or marked "no problems", execution stops. There is nothing to build.

---

**[ GATE 2 ]** — User confirms problem list and priority order before proceeding.

---

## LAYER 3 — LOCKED PATH

> **AI instruction:** Record all restrictions provided by the user. These are non-negotiable. Never suggest alternatives inside a locked path. Never explore approaches that contradict these entries.

### [DECISION] — Technical decisions already made
<!-- Example: [DECISION] Using PostgreSQL. No other database will be considered. -->

### [THOUGHT] — Conceptual directions that are off-limits
<!-- Example: [THOUGHT] No microservices. The system will remain a monolith for now. -->

---

**[ GATE 3 ]** — User confirms all locked paths are correctly registered before proceeding.

---

## LAYER 4 — PLANNING

> **AI instruction:** Decompose the problems into ordered tasks. For each task, list which skills will be used or need to be created. Respect dependencies between tasks. Do not begin implementation here.

### Task List

| # | Task | Depends On | Skills |
|---|---|---|---|
| T1 | | — | |
| T2 | | T1 | |
| T3 | | T2 | |

### Skills Required

| Skill | Layer | Status |
|---|---|---|
| | | existing / to create |

### Execution Order
<!-- AI describes the order and reasoning for sequencing -->

---

**[ GATE 4 ]** — User approves the task plan and skill list before execution begins.

---

## LAYER 5 — EXECUTION

> **AI instruction:** Execute tasks in the approved order, one at a time. After each task, stop and report what was done. Do not proceed to the next task without confirmation. Do not do anything outside the task scope.

### T1 — [Task name]
- Status: `pending` / `in progress` / `done`
- Output:
- Notes:

### T2 — [Task name]
- Status: `pending` / `in progress` / `done`
- Output:
- Notes:

---

## LAYER 6 — TESTING

> **AI instruction:** For each completed task, generate and run the appropriate tests. Tests are part of task delivery — a task is only done when its tests pass.

### T1 Tests
- [ ] Unit:
- [ ] Acceptance:

### T2 Tests
- [ ] Unit:
- [ ] Acceptance:

---

## LAYER 7 — DOCUMENTATION

> **AI instruction:** Generate documentation based on what was actually implemented, not what was planned. If something in the plan was not built, it should not appear here.

### What Was Built
<!-- AI fills this based on execution output -->

### Decisions Made During Execution
<!-- AI fills this with any decisions taken during Tasks that are worth recording -->

### How to Run / Use
<!-- AI fills this -->

### Known Limitations
<!-- AI fills this -->
