---
name: forge-locked-path
description: Register and enforce technical decisions and conceptual directions that are off-limits for FORGE Layer 3. Use when recording restrictions, locking decisions, or preventing the AI from suggesting forbidden approaches.
---

# FORGE — Locked Path Layer Skills

Skills for LAYER 3: register restrictions that the AI must never question or work around.

| Skill | File | When to use |
|---|---|---|
| Decision | [decision.md](decision.md) | Register a technical decision already made — no alternatives explored |
| Thought | [thought.md](thought.md) | Register a conceptual direction that is off-limits — a philosophy or mindset restriction |

## How to use

1. Run **decision.md** for each technical choice the user has already made
2. Run **thought.md** for each conceptual/philosophical direction to exclude
3. Present all locked paths to the user for confirmation
4. **Stop at GATE 3** — wait for user confirmation before proceeding to Layer 4

## AI-agnostic rules

These skills work with **any AI model** (Qwen, Claude, Gemini, GPT). Rules:

- **Never suggest alternatives** to a locked decision
- **Never say "you could also consider X"** when X contradicts a locked path
- **DECISION vs THOUGHT**: A decision closes one door (technical choice). A thought closes a corridor (entire direction/mindset)
- **If a locked path appears to cause a problem, raise it explicitly** — don't quietly work around it
- **Gate behavior**: All locked paths must be confirmed by the user before planning begins
