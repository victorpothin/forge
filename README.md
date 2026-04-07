<div align="center">

```
  ███████╗ █████╗ ██╗  ██╗  ██████╗
  ██╔════╝██╔══██╗██║  ██║ ██╔════╝
  █████╗  ███████║███████║ ██║
  ██╔══╝  ██╔══██║██╔══██║ ██║
  ██║     ██║  ██║██║  ██║ ╚██████╗
  ╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝  ╚═════╝
```

**Focused · Ordered · Restricted · Guided Execution**

A structured method for working with AI on software projects.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/victorpothin/forge)](https://github.com/victorpothin/forge/releases)

**AI-agnostic.** Works with Qwen, Claude, Gemini, GPT, or any coding assistant.

</div>

FORGE ensures the AI understands your intent, respects your restrictions, and executes only what was asked — nothing more.

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
├── FORGE.md                  # Master template
├── .forgerc.json             # Configuration
├── AI-ADAPTER.md             # AI model behavior guide
├── LICENSE                   # MIT License
├── go.mod                    # Go module definition
├── main.go                   # CLI entry point
│
├── cmd/                      # CLI commands (cobra)
│   ├── root.go
│   ├── init.go
│   └── doctor.go
│
├── internal/
│   └── templates/            # Embedded FORGE templates + skills
│       ├── FORGE.md
│       ├── skills/
│       └── templates.go
│
├── skills/                   # Source skills (copied to cli/embed on build)
│   ├── context/
│   ├── problem/
│   ├── locked-path/
│   ├── planning/
│   ├── execution/
│   ├── testing/
│   └── docs/
│
└── examples/
    ├── dotnet-api/
    └── rust-backend/
```

## CLI

### Install

**From source (works now):**
```bash
git clone https://github.com/victorpothin/forge.git
cd forge && go build -o forge .
sudo mv forge /usr/local/bin/
```

**Via go install (requires a release tag):**
```bash
go install github.com/victorpothin/forge@latest
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
