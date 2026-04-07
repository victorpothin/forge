# Skill: Context Extract

## Purpose
Read all available project materials and produce a structured summary of what exists — without interpreting intent or suggesting improvements.

## Trigger
Use this skill at the start of LAYER 1 before writing anything into the Context section of FORGE.md.

## Instructions for AI

1. Read all files provided: source code, config files, documentation, migration files, test files
2. For each area, describe only what you observe — not what you think should be there
3. Structure your output using these sections:
   - **Project**: name, purpose, domain
   - **Tech Stack**: languages, frameworks, runtimes, databases
   - **Current State**: what works, what is incomplete, what appears broken
   - **Integrations**: external services, APIs, queues, storage
   - **Entry Points**: how the application starts, how it's invoked
4. If something is ambiguous, write it as a question — do not guess
5. Stop after extraction. Do not make recommendations.

## Output Format

```
## Project
[description]

## Tech Stack
[list]

## Current State
[observations]

## Integrations
[list]

## Entry Points
[list]

## Open Questions
- [question 1]
- [question 2]
```

## What NOT to do
- Do not suggest improvements
- Do not infer intent from code
- Do not describe what should exist — only what does
