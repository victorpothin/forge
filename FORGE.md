# FORGE — Project Template

> **AI-agnostic.** Works with Qwen, Claude, Gemini, GPT, or any coding assistant.
> Copy this file into your project. Fill each layer in order. **Do not advance** without user validation.

---

## AI Instructions

**Read these files first:**
- `.forgerc.json` — AI model, gate mode, active layers
- `.forge-memory` — persistent project rules and context (load only what's relevant)

**Rules that apply to ALL AI models:**
1. Execute layers **in order** (1→7). Never skip.
2. **Stop at each `[ GATE ]`** and wait for user confirmation.
3. Do NOT infer or fill sections without explicit user input when required.
4. Do NOT suggest alternatives when a `[DECISION]` or `[THOUGHT]` is locked.
5. Documentation (Layer 7) must reflect **what was actually built**, not what was planned.

**About `.forge-memory`:**
- This file stores persistent project rules (e.g., "PostgreSQL only", "no microservices")
- **Read it before starting each session** — apply memories as implicit Locked Paths
- Only load memories relevant to the current context — don't dump everything
- Never delete or modify this file — use `forge memory` commands

**About the `forge` CLI:**
You can run `forge` commands directly from the project root. Use them to
manage layers, skills, memories, and verify project health.

### Available commands

| Command | When to use |
|---|---|
| `forge memory list` | Check what rules are stored before giving recommendations |
| `forge memory add -t <tag> "text"` | Save a rule the user confirmed (e.g. after Layer 3 Locked Path) |
| `forge memory delete <id>` | Remove a rule the user asked to delete |
| `forge skill add --from <src> --layer <l>` | Import an external skill the user provided |
| `forge skill list` | See what skills are available for this project |
| `forge edit --add <layers>` / `--remove <layers>` | Enable or disable layers at the user's request |
| `forge doctor` | Quick health check if the user reports issues |

### Rules for running commands
- Always run from the project root (`.` or `--dir <path>`)
- Use `-y` flag to skip confirmation **only** when the user explicitly agreed
- Never run destructive commands (`memory clear`, `edit --remove`) without confirming with the user first
- After running a command, report the output back to the user

---

## LAYER 1 — CONTEXT

> **Instruction:** Read all project files, code, and materials the user provided. Synthesize your understanding below. Do not infer intent — only describe what exists. Then **stop and wait for validation**.

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

> **Instruction:** Based on the validated context, list the problems the user needs to solve. **Do NOT suggest solutions.** List only problems, ordered by priority. If there are no problems, state that explicitly — nothing will be executed.

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

> **Instruction:** Record all restrictions provided by the user. These are non-negotiable. **Never suggest alternatives** inside a locked path. **Never explore approaches** that contradict these entries.

### [DECISION] — Technical decisions already made
<!-- Example: [DECISION] Using PostgreSQL. No other database will be considered. -->

### [THOUGHT] — Conceptual directions that are off-limits
<!-- Example: [THOUGHT] No microservices. The system will remain a monolith for now. -->

---

**[ GATE 3 ]** — User confirms all locked paths are correctly registered before proceeding.

---

## LAYER 4 — PLANNING

> **Instruction:** Decompose the problems into ordered tasks. For each task, list which skills will be used. Respect dependencies between tasks. **Do not begin implementation here.**

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
<!-- Describe the order and reasoning for sequencing -->

---

**[ GATE 4 ]** — User approves the task plan and skill list before execution begins.

---

## LAYER 5 — EXECUTION

> **Instruction:** Execute tasks in the approved order, **one at a time**. After each task, stop and report what was done. **Do not proceed** to the next task without confirmation. **Do not do anything outside the task scope.**

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

> **Instruction:** For each completed task, generate and run the appropriate tests. Tests are part of task delivery — a task is only `done` when its tests pass.

### T1 Tests
- [ ] Unit:
- [ ] Acceptance:

### T2 Tests
- [ ] Unit:
- [ ] Acceptance:

---

## LAYER 7 — DOCUMENTATION

> **Instruction:** Generate documentation based on **what was actually implemented**, not what was planned. If something in the plan was not built, it should not appear here.

### What Was Built
<!-- AI fills this based on execution output -->

### Decisions Made During Execution
<!-- AI fills this with any decisions taken during Tasks that are worth recording -->

### How to Run / Use
<!-- AI fills this -->

### Known Limitations
<!-- AI fills this -->
