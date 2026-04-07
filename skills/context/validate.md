# Skill: Context Validate

## Purpose
After extraction, verify that the synthesized context is consistent, complete, and not based on assumptions.

## Trigger
Use this skill after `extract.md` produces output and before GATE 1.

## Instructions for AI

1. Review the extracted context for internal contradictions
2. Check if any section was filled with assumptions rather than observations
3. For every Open Question from the extract phase, mark it as:
   - `[ANSWERED]` if the answer was found elsewhere in the project
   - `[STILL OPEN]` if it needs user input
4. If critical information is missing (e.g., no entry point found, no database config found), flag it explicitly
5. Present a confidence score per section: HIGH / MEDIUM / LOW

## Output Format

```
## Validation Report

### Project — [HIGH / MEDIUM / LOW]
[notes]

### Tech Stack — [HIGH / MEDIUM / LOW]
[notes]

### Current State — [HIGH / MEDIUM / LOW]
[notes]

### Open Questions
- [ANSWERED] Question X → Answer Y
- [STILL OPEN] Question Z — needs user input

### Flags
- [FLAG] Missing: [what is missing and why it matters]
```

## Gate behavior
If any section is LOW confidence or there are STILL OPEN questions, do not advance to LAYER 2. Present the flags to the user and wait for answers.
