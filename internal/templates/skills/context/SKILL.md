---
name: forge-context
description: Extract and validate project context for FORGE Layer 1. Reads all project files, code, and materials to build a structured understanding. Use when starting a FORGE session, analyzing a project, or filling the Context layer.
---

# FORGE — Context Layer Skills

Skills for LAYER 1: extract project understanding and validate it before proceeding.

| Skill | File | When to use |
|---|---|---|
| Extract | [extract.md](extract.md) | Read all available materials and build a structured summary of what exists |
| Validate | [validate.md](validate.md) | After extraction, check consistency and flag gaps before GATE 1 |

## How to use

1. Run **extract.md** first — read code, config, docs, tests
2. Run **validate.md** — check for contradictions, missing info, low-confidence sections
3. Present the synthesized context to the user
4. **Stop at GATE 1** — wait for user validation before proceeding to Layer 2

## AI-agnostic rules

These skills work with **any AI model** (Qwen, Claude, Gemini, GPT). Rules:

- **Only describe what exists** — never infer intent or suggest improvements
- **If something is ambiguous, write it as a question** — do not guess
- **Stop after extraction** — do not make recommendations
- **Gate behavior**: If any section is LOW confidence or has open questions, present flags to the user and wait
