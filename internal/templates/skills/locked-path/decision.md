# Skill: Locked Path — Decision

## Purpose
Register a technical decision that has already been made. The AI must never suggest alternatives or question this decision.

## Trigger
Use this skill when filling the [DECISION] entries in LAYER 3.

## Format

```
[DECISION] {what was decided}
Reason: {why — optional, but helps the AI understand the boundary}
Scope: {where this applies — project-wide / specific module / specific layer}
```

## Examples

```
[DECISION] PostgreSQL is the only database in this project.
Reason: Already in production, migration cost is too high.
Scope: Project-wide.

[DECISION] All HTTP endpoints use REST, not GraphQL.
Reason: Client team is not ready for GraphQL.
Scope: API layer.

[DECISION] Authentication is handled by an external identity provider.
Reason: Compliance requirement.
Scope: Project-wide. Do not implement custom auth logic.
```

## AI behavior rule
When a DECISION is registered:
- Never suggest the alternative approach
- Never say "you could also consider X"
- If the locked approach appears to cause a problem, raise it explicitly rather than quietly working around it
