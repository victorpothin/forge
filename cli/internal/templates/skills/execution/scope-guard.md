# Skill: Scope Guard

## Purpose
Detect when the AI is about to do something outside the current task's declared scope and stop it before it happens.

## Trigger
Apply this skill continuously during execution. It is a passive guard, not a separate step.

## What counts as scope violation

| Type | Example |
|---|---|
| File outside declared scope | Task touches `/config` but scope was `/handlers` |
| Refactoring unrequested code | "Cleaning up" code nearby while implementing the task |
| Adding features not in the task | Task is "add endpoint", AI also adds pagination |
| Changing behavior of adjacent code | Task adds a field, AI also changes how an existing field works |
| Solving a problem not in the problem list | AI notices something "wrong" and fixes it |

## When a violation is detected

Stop. Do not proceed with the out-of-scope action. Report using this format:

```
## Scope Guard — Alert

I was about to {describe the out-of-scope action}.

This is outside the scope of T{n} ({task name}).

Options:
A) Skip this and continue with T{n} as scoped
B) Add a new task to the plan for this (requires returning to LAYER 4)
C) Expand T{n} scope — requires your approval

How would you like to proceed?
```

## Why this exists
The AI naturally tends to "improve" code while working on it. This is usually well-intentioned but causes:
- Unreviewable diffs
- Unintended behavior changes
- Loss of trust in what was actually delivered

Scope guard keeps each task's output predictable and reviewable.
