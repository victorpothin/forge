# Skill: Problem Gate

## Purpose
Enforce the core FORGE rule: if there is no problem, there is nothing to execute.

## Trigger
Use this skill after the problem list is filled and before advancing to LAYER 3.

## Instructions for AI

Check the problem list and apply this decision tree:

```
Is the problem list empty?
├── YES → Output: "No problems identified. Execution will not proceed. FORGE stops here."
│         Stop. Do not advance.
└── NO  → Are all problems already solved based on current state?
          ├── YES → Output: "All problems are already resolved in the current state."
          │         List what solved each one. Stop. Do not advance.
          └── NO  → Proceed to GATE 2.
```

## Output Format (when proceeding)

```
## Gate 2 Check — PASS

Problems confirmed:
- P1: [name] — active, not yet solved
- P2: [name] — active, not yet solved

Ready to proceed to Locked Path layer.
```

## Output Format (when stopping)

```
## Gate 2 Check — STOP

Reason: [no problems / all resolved]
[explanation]

FORGE execution halted. Nothing will be built.
```

## Why this exists
The AI tends to start building things even when there is nothing to build. This gate makes the absence of work an explicit, documented outcome — not a failure.
