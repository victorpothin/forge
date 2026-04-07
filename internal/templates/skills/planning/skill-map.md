# Skill: Planning Skill Map

## Purpose
For each task in the plan, identify which FORGE skills will be used during execution. Skills are listed here, not created. Creation happens in the execution layer.

## Trigger
Use this skill after `decompose.md` produces the task list.

## Instructions for AI

1. Review each task and identify what type of work it requires
2. Map the work type to existing skills in the `skills/` directory
3. If a task requires a skill that doesn't exist yet, mark it as "to create" with a description
4. Do not create skills here — only list them

## Skill categories to check

| Work type | Where to look |
|---|---|
| Reading/understanding code | `skills/context/` |
| Analyzing what to build | `skills/problem/` |
| Checking restrictions | `skills/locked-path/` |
| Generating/running tests | `skills/testing/` |
| Writing documentation | `skills/docs/` |
| Technical skills (language/framework specific) | `skills/execution/` or custom |

## Output Format

```
## Skill Map

| Task | Skills Used | Status |
|------|------------|--------|
| T1   | skills/execution/task-runner.md | existing |
| T2   | skills/execution/task-runner.md, skills/testing/unit.md | existing |
| T3   | skills/testing/acceptance.md | existing |
| T4   | skills/execution/[new: dotnet-migration-runner] | to create |

## Skills to Create
- `skills/execution/dotnet-migration-runner.md` — runs EF Core migrations and validates schema output
```

## Rules
- Only list skills that will actually be used — no speculative mapping
- If a task needs no special skill beyond basic coding, write "general execution"
- Skills to create should be minimal — one skill per distinct repeatable pattern
