# Skill: Problem Prioritize

## Purpose
Order the identified problems by real impact on the project, not by technical convenience or what seems easiest to tackle first.

## Trigger
Use this skill when filling LAYER 2 of FORGE.md.

## Instructions for AI

For each problem identified in context, evaluate:

1. **Impact**: If this problem is not solved, what breaks or is blocked?
2. **Dependency**: Does solving this unblock other problems?
3. **Scope**: Is this one problem or multiple problems lumped together? Split if needed.
4. **User intent**: Did the user explicitly mention this, or are you inferring it?

Assign each problem a priority based on:
- P1: Blocking — nothing else can proceed without this
- P2: High — significant friction or risk if left unsolved
- P3: Medium — worth solving but not urgent
- P4: Low — nice to have, do last or skip

## Output Format

```
## Problem List

### P1 — [Problem name]
[One sentence describing the problem]
Impact: [what breaks if this isn't solved]
Dependency: [what this blocks]

### P2 — [Problem name]
[description]
Impact: [...]
Dependency: [...]
```

## Rules
- Only list problems the user asked to solve or that directly block their goals
- Do not add problems you think should be solved — that is out of scope
- If two problems are actually the same, merge them
- If a problem is actually a symptom of another, list the root cause as the real problem
