# FORGE

**Focused, Ordered, Restricted, Guided Execution**

FORGE is a structured method for working with AI on software projects. It ensures the AI understands your intent, respects your restrictions, and executes only what was asked — nothing more.

**AI-agnostic.** Works with Qwen, Claude, Gemini, GPT, or any coding assistant. Mix models across layers if you want.

## Why FORGE?

Most AI workflows fail because the AI:
- Solves problems you didn't ask it to solve
- Ignores decisions you've already made
- Goes off-scope during execution
- Documents what was planned instead of what was built

FORGE fixes this by enforcing a strict flow with gates between layers.

## The 7 Layers

| Layer | Purpose |
|---|---|
| **Context** | AI extracts and synthesizes project understanding, then you validate |
| **Problem** | Explicit prioritized list of what needs solving — if no problem, nothing executes |
| **Locked Path** | Decisions and thoughts that are off-limits — technical and conceptual |
| **Planning** | Task breakdown with skill mapping and dependency ordering |
| **Execution** | Task-by-task implementation, strictly scoped |
| **Testing** | Validation built into each task delivery |
| **Documentation** | Generated from what was actually implemented |

## Gate Rule

**No layer starts before the previous one is validated.**

The AI must stop and confirm with you before advancing. This is not optional.

## Repository Structure

```
forge/
├── README.md
├── FORGE.md                  # Master template — start here
├── .forgerc.json             # Configuration — AI model, gates, layers
├── AI-ADAPTER.md             # How each AI model should behave with FORGE
├── LICENSE                   # MIT License
│
├── cli/                      # Go CLI — forge init, forge doctor
│   ├── cmd/
│   └── internal/templates/
│
├── skills/
│   ├── context/              # Layer 1: extract & validate project understanding
│   ├── problem/              # Layer 2: identify & prioritize problems
│   ├── locked-path/          # Layer 3: register restrictions
│   ├── planning/             # Layer 4: decompose, map skills, order tasks
│   ├── execution/            # Layer 5: run tasks, enforce scope
│   ├── testing/              # Layer 6: unit & acceptance tests
│   └── docs/                 # Layer 7: document what was built
│
└── examples/
    ├── dotnet-api/           # Qwen-only workflow example
    └── rust-backend/         # Mixed model workflow example (Qwen + Claude)
```

## CLI

### Install

```bash
go install github.com/forge-cli/forge@latest
```

Or from source:

```bash
git clone https://github.com/victorpothin/forge.git && cd forge/cli && go build -o forge . && sudo mv forge /usr/local/bin/
```

### Usage

```bash
# Interactive wizard
forge init

# Non-interactive — just specify the AI
forge init --ai qwen -y

# Override defaults
forge init --ai claude --model claude-sonnet-4-20250514 --gate-mode strict -y

# Force overwrite existing files
forge init --ai gpt --force -y

# Check project health
forge doctor
forge doctor --dir /path/to/project
```

### Commands

| Command | Description |
|---|---|
| `forge init` | Initialize FORGE in a project (interactive or flags) |
| `forge doctor` | Check if FORGE is properly set up |
| `forge --version` | Show CLI version |

### Init flags

| Flag | Short | Description |
|---|---|---|
| `--ai <model>` | `-a` | AI model: `qwen`, `claude`, `gpt`, `gemini`, `custom` |
| `--model <name>` | `-m` | Model name (e.g. `qwen-code`, defaults to AI-specific value) |
| `--gate-mode <mode>` | `-g` | Gate mode: `strict` (default) or `auto` |
| `--dir <path>` | `-d` | Target directory (default: `.`) |
| `--yes` | `-y` | Skip confirmation |
| `--force` | | Overwrite existing FORGE files |

### Manual setup (no CLI)

1. Copy `FORGE.md` and `.forgerc.json` into your project
2. Edit `.forgerc.json` to set your AI model (`qwen`, `claude`, `gemini`, `gpt`, `custom`)
3. Copy `skills/` to your AI's skill directory (see AI-ADAPTER.md)
4. Start a session and reference `FORGE.md`

## License

MIT — see [LICENSE](LICENSE) for details.
