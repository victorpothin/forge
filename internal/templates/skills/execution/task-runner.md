# Skill: Task Runner

## Purpose
Execute one task from the approved plan. Produce its defined output. Stop. Report. Wait for confirmation before the next task.

## Trigger
Use this skill at the start of each task in LAYER 5.

## Instructions for AI

Before starting the task:
1. State the task name and number
2. State what the output will be when done
3. State the scope (files/modules that will be touched)
4. State which locked paths are relevant (even if not violated)

During execution:
1. Work only within the declared scope
2. If you discover that the task requires touching something outside scope, stop and report — do not silently expand scope
3. If you discover a conflict with a locked path, stop and report — do not work around it
4. If the task is harder than expected, report the blocker — do not invent a different solution

After execution:
1. List exactly what was created or changed
2. List what was NOT done (if scope had to be limited)
3. State whether tests are needed (pass to testing skill)
4. Wait for user confirmation before proceeding to next task

## Output Format

```
## Executing T{n} — {Task name}

**Output target:** {what will exist when done}
**Scope:** {files/modules}
**Locked path check:** {none / flags}

---

[implementation]

---

## T{n} Complete

**Created/Changed:**
- {file}: {what changed}

**Not done:**
- {anything explicitly deferred}

**Tests needed:** yes / no
**Status:** awaiting confirmation to proceed to T{n+1}
```

## Rules
- One task per execution call
- Never start T(n+1) in the same response as T(n)
- If unsure about scope, ask — don't assume
