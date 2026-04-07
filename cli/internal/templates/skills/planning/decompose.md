# Skill: Planning Decompose

## Purpose
Break each problem into concrete, independently executable tasks — small enough to be completed in one session, large enough to produce a meaningful output.

## Trigger
Use this skill first in LAYER 4 planning.

## Instructions for AI

For each problem in the prioritized list:

1. Identify the smallest unit of work that produces a testable output
2. Name each task as an action: verb + object (e.g., "Add JWT validation middleware", not "JWT")
3. State what the task produces — not what it does, but what exists after it's done
4. Mark if a task touches a locked path — flag it, don't block
5. One problem can produce multiple tasks, but tasks should not span multiple problems

## Task format

```
T{n} — {Action: verb + object}
Problem: P{n}
Output: {what exists when this task is done}
Scope: {files, modules, or layers affected}
Locked Path check: none / [FLAG: conflicts with THOUGHT/DECISION X]
```

## Example

```
T1 — Add database migration for user sessions table
Problem: P1
Output: Migration file applied, table exists in schema
Scope: /migrations, /db
Locked Path check: none

T2 — Implement session creation endpoint
Problem: P1
Output: POST /sessions returns token on valid credentials
Scope: /handlers, /services
Locked Path check: none

T3 — Add integration test for session flow
Problem: P1
Output: Test passes end-to-end against real database
Scope: /tests
Locked Path check: none
```

## Rules
- Tasks must not overlap in scope
- A task is too big if it requires more than one session to complete — split it
- A task is too small if its output is not independently testable — merge it
- Do not add tasks for problems not in the list
