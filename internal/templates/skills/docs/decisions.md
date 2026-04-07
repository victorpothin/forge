# Skill: Docs Decisions

## Purpose
Record decisions made during execution that are worth preserving — choices that weren't obvious, tradeoffs that were made, or things that future developers will ask "why was this done this way?"

## Trigger
Use alongside `implemented.md` in LAYER 7.

## Difference from Locked Path
- **Locked Path** = decisions made BEFORE execution, by the user
- **Execution Decisions** = decisions made DURING execution, by the AI with user confirmation

## What qualifies as a decision worth recording

| Worth recording | Not worth recording |
|---|---|
| Chose approach A over B because of constraint X | Used a for loop instead of LINQ |
| Had to deviate from task scope for reason Y | Named a variable `result` |
| Found that the original plan had a flaw — here's what was done instead | Applied standard error handling |
| Accepted a known tradeoff: faster now, will need rework when Z | Followed existing code conventions |

## Format

Use the ADR (Architecture Decision Record) pattern — lightweight version:

```
## Decision — {title}

**Context:** {what situation led to this decision}
**Decision:** {what was chosen}
**Reason:** {why — constraint, tradeoff, limitation}
**Consequences:** {what this means going forward — what becomes easier, what becomes harder}
```

## Example

```
## Decision — Used in-memory cache instead of Redis for session storage

Context: Task T3 required fast session lookups. Redis was considered but not available in the dev environment.
Decision: Used IMemoryCache for now.
Reason: Unblocks development. Redis integration was not in scope for this session.
Consequences: Sessions are not shared across instances. Must replace with Redis before multi-instance deploy.
```

## Output Format

```
## Execution Decisions

### Decision — {title}
Context: ...
Decision: ...
Reason: ...
Consequences: ...

### Decision — {title}
...
```
