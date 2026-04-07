# Skill: Planning Order

## Purpose
Sequence tasks so that dependencies are respected, locked paths are not violated, and each task can begin as soon as the previous one is done.

## Trigger
Use this skill after `decompose.md` and `skill-map.md` are complete.

## Instructions for AI

1. For each task, identify what it requires to be true before it can start (input dependency)
2. For each task, identify what becomes true after it finishes (output dependency)
3. Build a dependency chain — not a flat list
4. Check every task against the locked path list — if a task would require going through a locked path, flag it
5. If two tasks have no dependency on each other, mark them as parallelizable

## Dependency types

| Type | Meaning |
|---|---|
| HARD | Task B cannot start until Task A is done |
| SOFT | Task B is easier after A, but can technically proceed |
| PARALLEL | A and B can be done at the same time |
| BLOCKED | Task requires going through a locked path — must be discussed |

## Output Format

```
## Execution Order

T1 → T2 → T3
         ↘
          T4 (parallel with T3)

## Dependency Table

| Task | Depends On | Type | Notes |
|------|-----------|------|-------|
| T1   | —         | —    | First task, no dependencies |
| T2   | T1        | HARD | Needs schema from T1 |
| T3   | T2        | HARD | Needs endpoint from T2 |
| T4   | T2        | PARALLEL | Can run while T3 is being built |

## Flags
- none
```

## Rules
- Never reorder tasks based on what seems more interesting — follow dependencies
- If ordering reveals a circular dependency, raise it before proceeding
- Keep the order as flat as possible — deep chains mean fragile plans
