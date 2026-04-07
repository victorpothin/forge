# Skill: Locked Path — Thought

## Purpose
Register a conceptual direction the user does not want explored. Unlike [DECISION], this is not about a specific technical choice — it's about an entire way of thinking or architectural philosophy that is off-limits.

## Trigger
Use this skill when filling the [THOUGHT] entries in LAYER 3.

## Format

```
[THOUGHT] {the approach or mindset that is locked out}
Reason: {why — optional}
Implication: {what this means for how solutions should be shaped}
```

## Examples

```
[THOUGHT] No microservices. The system is a monolith and will remain one.
Reason: Team size and ops capacity don't support distributed systems right now.
Implication: All solutions must fit within a single deployable unit.

[THOUGHT] Do not over-engineer for scale we don't have.
Reason: We are a small team and premature optimization has burned us before.
Implication: Simple solutions are preferred over flexible but complex ones.

[THOUGHT] No AI-generated UX patterns. The UI must be designed by a human first.
Reason: Product design is not part of this scope.
Implication: Do not suggest layout changes or UX improvements.
```

## Difference from [DECISION]

| [DECISION] | [THOUGHT] |
|---|---|
| A specific technical choice | A philosophical direction |
| "We use PostgreSQL" | "We don't want a distributed system" |
| Closes one door | Closes a corridor of doors |
| Usually reversible with migration | Usually reflects a team/org constraint |

## AI behavior rule
When a THOUGHT is registered:
- Do not explore or suggest anything in that direction
- If a task naturally leads toward a locked thought, flag it and ask the user how to proceed
- Never say "since you can't do X, you should do Y (which is also X but smaller)"
