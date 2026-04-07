---
name: forge-planning
description: Decompose problems into ordered tasks with skill mappings for FORGE Layer 4. Breaks problems into concrete, testable tasks with dependency ordering. Use when creating a task plan, mapping skills to tasks, or sequencing work.
---

# FORGE — Planning Layer Skills

Skills for LAYER 4: decompose problems into executable tasks with dependencies and skill mappings.

| Skill | File | When to use |
|---|---|---|
| Decompose | [decompose.md](decompose.md) | Break problems into concrete, scoped, testable tasks |
| Skill-map | [skill-map.md](skill-map.md) | For each task, identify which FORGE skills will be used or need creation |
| Order | [order.md](order.md) | Sequence tasks respecting dependencies and locked paths |

## How to use

1. Run **decompose.md** — break each problem into smallest testable units of work
2. Run **skill-map.md** — map each task to existing or new skills
3. Run **order.md** — build dependency chain, identify parallel tasks, check locked paths
4. Present the full plan to the user
5. **Stop at GATE 4** — wait for user approval before execution begins

## AI-agnostic rules

These skills work with **any AI model** (Qwen, Claude, Gemini, GPT). Rules:

- **Do NOT begin implementation** — this layer is planning only
- **Tasks must be testable** — if output is not independently verifiable, merge or split
- **Tasks must not span multiple problems** — one problem → N tasks, but no cross-problem tasks
- **Flag locked path conflicts** — if a task conflicts with a locked path, flag it, don't block
- **Gate behavior**: Full plan (task list, skill map, execution order) must be approved by user before Layer 5
