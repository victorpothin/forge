---
name: forge-execution
description: Execute tasks one at a time with scope enforcement for FORGE Layer 5. Runs a single task, reports output, stops for confirmation. Use when implementing a planned task, writing code from a task list, or enforcing scope boundaries.
---

# FORGE — Execution Layer Skills

Skills for LAYER 5: execute approved tasks one at a time with strict scope control.

| Skill | File | When to use |
|---|---|---|
| Task-runner | [task-runner.md](task-runner.md) | Execute one task from the approved plan, produce output, stop |
| Scope-guard | [scope-guard.md](scope-guard.md) | Detect and prevent out-of-scope actions during task execution (passive) |

## How to use

1. **Scope-guard** is always active — it monitors for scope violations during execution
2. Run **task-runner.md** for each task in the approved order
3. After each task completes, report what was done and **wait for user confirmation**
4. Only proceed to the next task after explicit approval

## AI-agnostic rules

These skills work with **any AI model** (Qwen, Claude, Gemini, GPT). Rules:

- **One task per execution** — never start T(n+1) in the same response as T(n)
- **Work only within declared scope** — if something outside scope is needed, stop and report
- **If a locked path conflict is discovered, stop and report** — don't work around it
- **If the task is harder than expected, report the blocker** — don't invent a different solution
- **After completion**: list what was created/changed, what was NOT done, whether tests are needed
- **Gate behavior**: Each task reports completion and waits for confirmation before the next
