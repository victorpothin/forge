# FORGE — AI Adapter Instructions

> This file ensures consistent behavior across all AI models (Qwen, Claude, Gemini, GPT, etc.).
> It is a reference for how each model should interpret and execute the FORGE method.

---

## Core Rules (All Models)

These rules apply regardless of which AI is executing FORGE:

1. **Read `.forgerc.json` first** — it defines gate mode, active layers, and model assignments
2. **Execute layers in order (1→7)** — never skip a layer
3. **Stop at each `[ GATE ]`** — wait for explicit user confirmation before advancing
4. **No inference** — describe what exists, not what you think should exist
5. **No alternatives** when a path is locked — never suggest "you could also consider X"
6. **One task at a time** in execution — never batch tasks
7. **Document reality** — Layer 7 reflects what was built, not what was planned

---

## Model-Specific Notes

### Qwen Code

- **Skills system**: Uses `.qwen/skills/` (project) or `~/.qwen/skills/` (personal)
- **Invocation**: Skills are model-invoked based on description. Ensure SKILL.md `description` field is specific
- **Sub-agents**: Can use `/agent` for parallel exploration when needed
- **Approval modes**: `plan`, `default`, `auto_edit`, `yolo` — FORGE works best with `plan` or `default`
- **Setup**: Copy `skills/` into `.qwen/skills/` in your project, or reference via `FORGE.md`

### Claude (Claude Code / CLI)

- **Skills system**: Uses `.claude/skills/` directory with `SKILL.md` files
- **Invocation**: Claude discovers skills from the project directory
- **Custom instructions**: Claude respects `.claude/CLAUDE.md` for project-level instructions
- **Setup**: Copy `FORGE.md` into `.claude/` and reference it, or place in project root

### Gemini (Gemini CLI)

- **Skills system**: Uses `.gemini/` directory for project context
- **Invocation**: Provide `FORGE.md` as context at session start
- **Setup**: Reference `FORGE.md` explicitly: "Follow the FORGE method from FORGE.md"

### GPT (ChatGPT / GPT CLI)

- **Skills system**: No native skills — rely on structured prompting
- **Invocation**: Paste `FORGE.md` content or provide as a file upload
- **Setup**: Start session with: "We will use the FORGE method. Here are the rules: [paste FORGE.md]"
- **Note**: Without native skill gating, the AI must be reminded to stop at gates

---

## Gate Enforcement

The gate mechanism is the most important part of FORGE. Here's how each model should handle it:

### Strict Gate Mode (default)

```
AI: [Completes Layer N]
AI: "[ GATE N ] — I have completed Layer N. Do you validate this output? 
      Reply 'confirmed' to proceed, or tell me what to correct."
AI: [WAITS. Does not proceed until user responds.]
```

### What NOT to do at a gate

❌ "Should I continue to the next layer?" (implies eagerness)
❌ "Here's the next layer analysis..." (didn't wait)
❌ "I think this looks good, moving on..." (self-validated)

### What TO do at a gate

✅ "Layer N complete. [summary]. Do you validate this? I will wait for your response before proceeding."
✅ Then **stop generating output** until the user responds.

---

## Mixing Models Across Layers

You can use different AI models for different layers via `.forgerc.json`:

```json
{
  "layer_ai": {
    "context": "qwen",
    "problem": "claude",
    "locked_path": "qwen",
    "planning": "claude",
    "execution": "qwen",
    "testing": "qwen",
    "documentation": "claude"
  }
}
```

**When to mix:**
- Claude excels at: reasoning, planning, documentation
- Qwen excels at: code execution, file operations, testing
- GPT excels at: creative problem identification

**Rules for mixing:**
- Each model must read the full `FORGE.md` state from previous layers
- The output of one model becomes the input context for the next
- Gates are enforced regardless of model switches

---

## Starting a FORGE Session

### Qwen Code
```
1. Place FORGE.md in project root
2. Copy skills/ to .qwen/skills/ (or reference them)
3. Create .forgerc.json in project root
4. Start session: "Start a FORGE session. Read FORGE.md and .forgerc.json. Begin Layer 1."
```

### Claude Code
```
1. Place FORGE.md in project root or .claude/
2. Copy skills/ to .claude/skills/
3. Create .forgerc.json
4. Start session: "@FORGE.md — begin Layer 1 (Context)."
```

### Any other model
```
1. Place FORGE.md in project root
2. Create .forgerc.json
3. Start session: "We will use the FORGE method. Read FORGE.md and .forgerc.json. Begin Layer 1."
```
