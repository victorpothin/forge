---
name: forge-problem
description: Identify and prioritize problems for FORGE Layer 2. Lists problems ordered by impact, with a gate that stops execution if no problems exist. Use when defining what needs to be solved, filling the Problem layer, or checking if work should proceed.
---

# FORGE — Problem Layer Skills

Skills for LAYER 2: identify and prioritize what actually needs to be solved.

| Skill | File | When to use |
|---|---|---|
| Prioritize | [prioritize.md](prioritize.md) | Order problems by real impact — not by what seems easiest |
| Gate | [gate.md](gate.md) | Verify there is a real problem before allowing execution to proceed |

## How to use

1. Run **prioritize.md** — evaluate each problem by impact, dependency, scope, and user intent
2. Run **gate.md** — enforce the core FORGE rule: no problem = nothing to execute
3. Present the prioritized list to the user
4. **Stop at GATE 2** — wait for user confirmation before proceeding to Layer 3

## AI-agnostic rules

These skills work with **any AI model** (Qwen, Claude, Gemini, GPT). Rules:

- **Do NOT suggest solutions** — list only problems
- **Do NOT add problems you think should be solved** — only problems the user asked for or that directly block their goals
- **If the problem list is empty, state it explicitly** — execution stops, this is not a failure
- **Gate behavior**: The gate skill makes the decision to proceed or halt — no automatic advancement
