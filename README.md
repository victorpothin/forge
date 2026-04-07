<div align="center">

# 🔥 FORGE

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
│   ├── edit.go
│   ├── skill.go
│   ├── doctor.go
│   └── update.go
│
├── internal/
│   ├── templates/            # Embedded FORGE templates + skills
│   │   ├── FORGE.md
│   │   ├── skills/
│   │   └── templates.go
│   └── ui/                   # CLI UI helpers
│       └── ui.go
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

**One command — just works:**
```bash
curl -fsSL https://raw.githubusercontent.com/victorpothin/forge/main/install.sh | bash
```

Installs the binary and configures your PATH automatically. Requires Go.

**Via go install:**
```bash
go install github.com/victorpothin/forge@latest
```

**From source:**
```bash
git clone https://github.com/victorpothin/forge.git
cd forge && make install
```

### Usage

```bash
# Initialize a project (interactive wizard)
forge init

# Non-interactive — choose AI, layers, and confirm
forge init --ai qwen --layers context,problem,execution -y
forge init --ai claude -m claude-sonnet-4 --gate-mode strict --force -y

# Manage skills
forge skill add                          # Interactive import
forge skill add --from ./my-skill --layer execution
forge skill add --from github.com/user/my-skill --layer planning
forge skill list                         # Show skills by layer

# Manage layers after init
forge edit                               # Interactive toggle
forge edit --add testing,docs            # Enable layers
forge edit --remove context,documentation # Disable layers

# Check project health
forge doctor
forge doctor --dir /path/to/project

# Keep CLI up to date
forge update                             # Check and prompt
forge update -y                          # Auto-update
```

### Commands

| Command | Description |
|---|---|
| `forge init` | Initialize FORGE in a project (interactive or flags) |
| `forge skill add` | Import a skill from local dir, URL, or git repo |
| `forge skill list` | List all skills organized by layer |
| `forge edit` | Add or remove FORGE layers from a project |
| `forge doctor` | Check if FORGE is properly set up |
| `forge update` | Check for updates and self-update the CLI |
| `forge --version` | Show CLI version |

### Init flags

| Flag | Short | Description |
|---|---|---|
| `--ai <model>` | `-a` | AI model: `qwen`, `claude`, `gpt`, `gemini`, `custom` |
| `--model <name>` | `-m` | Model name (e.g. `qwen-code`, defaults to AI-specific value) |
| `--gate-mode <mode>` | `-g` | Gate mode: `strict` (default) or `auto` |
| `--layers <list>` | `-L` | Comma-separated layers (default: all 7) |
| `--dir <path>` | `-d` | Target directory (default: `.`) |
| `--yes` | `-y` | Skip confirmation |
| `--force` | | Overwrite existing FORGE files |

### Skill flags

| Flag | Short | Description |
|---|---|---|
| `--from <source>` | | Source: local path, URL, or `github.com/user/repo` |
| `--layer <layer>` | `-l` | Target: `context`, `problem`, `locked-path`, `planning`, `execution`, `testing`, `docs` |
| `--dir <path>` | `-d` | Target project directory (default: `.`) |
| `--force` | | Overwrite existing skill |

### Edit flags

| Flag | Short | Description |
|---|---|---|
| `--add <layers>` | | Comma-separated layers to enable |
| `--remove <layers>` | | Comma-separated layers to disable |
| `--dir <path>` | `-d` | Target project directory (default: `.`) |
| `--yes` | `-y` | Skip confirmation |

### Manual setup (no CLI)

1. Copy `FORGE.md` and `.forgerc.json` into your project
2. Edit `.forgerc.json` to set your AI model (`qwen`, `claude`, `gemini`, `gpt`, `custom`)
3. Copy `skills/` to your AI's skill directory (see AI-ADAPTER.md)
4. Start a session and reference `FORGE.md`

## License

MIT — see [LICENSE](LICENSE) for details.
