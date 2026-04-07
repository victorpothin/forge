---
name: forge-testing
description: Generate and run unit and acceptance tests for FORGE Layer 6. Tests are part of task delivery — a task is only done when tests pass. Use when validating task output, writing tests for new code, or verifying problem resolution.
---

# FORGE — Testing Layer Skills

Skills for LAYER 6: validate that each completed task works correctly through tests.

| Skill | File | When to use |
|---|---|---|
| Unit | [unit.md](unit.md) | Generate unit tests for specific functions, methods, or modules |
| Acceptance | [acceptance.md](acceptance.md) | Verify the task output satisfies the original problem's acceptance criteria |

## How to use

1. After a task is marked complete, run **unit.md** — test the behavior of what was built
2. Run **acceptance.md** — verify the task solves the problem it was created for
3. A task is only `done` when its tests pass — fix code or tests if they fail
4. Report results in the testing layer before proceeding to documentation

## AI-agnostic rules

These skills work with **any AI model** (Qwen, Claude, Gemini, GPT). Rules:

- **Use the project's existing test framework** — do not introduce a new one
- **Tests must pass** — if they don't, fix the code or the test and explain why
- **Acceptance criteria come from the Problem layer** — don't invent new criteria
- **"Partially solved" is valid** — document what remains and create a follow-up task
- **Gate behavior**: If tests fail, the task is not done — report, fix, retest before advancing
