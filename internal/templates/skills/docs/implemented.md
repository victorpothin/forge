# Skill: Docs Implemented

## Purpose
Generate documentation that accurately reflects what was built. No speculation, no planned features, no "future work" unless explicitly asked.

## Trigger
Use after all tasks are complete and all tests pass.

## Instructions for AI

1. Review the execution output for each completed task
2. For each task, extract:
   - What was created or changed
   - How to use it (commands, API calls, configuration)
   - What it depends on
3. Group related tasks into coherent documentation sections
4. Only document what exists in the code — not what was planned

## Sections to generate

### What Was Built
A plain-language summary of the new capabilities, changes, or fixes. One paragraph per problem solved.

### How to Use
Concrete instructions: commands, API endpoints, configuration values. Code examples where helpful.

### Dependencies and Requirements
What needs to be in place for this to work. Versions, environment variables, setup steps.

### Known Limitations
What this does NOT do. Edge cases that were explicitly deferred. Behaviors that differ from the original intent.

## Output Format

```
## Documentation

### What Was Built
[description]

### How to Use
[instructions + examples]

### Dependencies
[list]

### Known Limitations
- [limitation 1]
- [limitation 2]
```

## Rules
- If a planned task was not completed, do not document it
- If behavior differs from what was in the plan, document the actual behavior
- Keep it short enough to read in 5 minutes — if it needs more, split into separate docs
