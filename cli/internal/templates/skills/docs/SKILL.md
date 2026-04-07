---
name: forge-docs
description: Generate documentation from actual implementation for FORGE Layer 7. Documents what was built, not what was planned. Use when creating docs from completed work, recording execution decisions, or summarizing deliverables.
---

# FORGE — Documentation Layer Skills

Skills for LAYER 7: generate documentation from what was actually implemented.

| Skill | File | When to use |
|---|---|---|
| Implemented | [implemented.md](implemented.md) | Describe what was built, how it works, how to use it |
| Decisions | [decisions.md](decisions.md) | Record decisions made during execution that aren't in the Locked Path |

## How to use

1. Run **implemented.md** — extract what was created/changed from execution output
2. Run **decisions.md** — capture any non-obvious tradeoffs or choices made during execution
3. Generate documentation that reflects **reality**, not the original plan
4. Present to user — documentation is the final layer, no gate needed unless issues found

## AI-agnostic rules

These skills work with **any AI model** (Qwen, Claude, Gemini, GPT). Rules:

- **Document only what exists** — if a planned task was not completed, don't document it
- **If behavior differs from plan, document the actual behavior** — not the plan
- **Keep it short** — readable in 5 minutes; if longer, split into separate docs
- **Locked Path vs Execution Decisions**: Locked Path = decisions made BEFORE by user. Execution Decisions = choices made DURING by AI with user confirmation
- **Known limitations are required** — explicitly list what was deferred or differs from intent
