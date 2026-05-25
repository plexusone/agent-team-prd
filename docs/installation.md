# Installation

## Prerequisites

- Go 1.21 or later (for building from source)
- Git (for cloning the repository)

## Install with Go

The simplest way to install prdtool:

```bash
go install github.com/plexusone/agent-team-prd/cmd/prdtool@latest
```

This installs the `prdtool` binary to your `$GOPATH/bin` directory.

## Build from Source

Clone and build both the CLI and MCP server:

```bash
git clone https://github.com/plexusone/agent-team-prd.git
cd agent-team-prd

# Build CLI
go build -o bin/prdtool ./cmd/prdtool

# Build MCP server (for AI assistant integrations)
go build -o bin/prdtool-mcp ./cmd/prdtool-mcp
```

## Verify Installation

```bash
prdtool --version
prdtool --help
```

## Shell Completion

Generate shell completion scripts for your shell:

=== "Bash"

    ```bash
    prdtool completion bash > /etc/bash_completion.d/prdtool
    ```

=== "Zsh"

    ```bash
    prdtool completion zsh > "${fpath[1]}/_prdtool"
    ```

=== "Fish"

    ```bash
    prdtool completion fish > ~/.config/fish/completions/prdtool.fish
    ```

=== "PowerShell"

    ```powershell
    prdtool completion powershell | Out-String | Invoke-Expression
    ```

## MCP Server Installation

The MCP server (`prdtool-mcp`) is required for AI assistant integrations. Build it alongside the CLI:

```bash
go build -o bin/prdtool-mcp ./cmd/prdtool-mcp
```

Place the binary in a location accessible to your AI assistant. See [Claude Code Integration](integrations/claude-code.md) or [Kiro IDE Integration](integrations/kiro-ide.md) for setup instructions.

## Updating

To update to the latest version:

```bash
go install github.com/plexusone/agent-team-prd/cmd/prdtool@latest
```

Or if building from source:

```bash
cd agent-team-prd
git pull
go build -o bin/prdtool ./cmd/prdtool
go build -o bin/prdtool-mcp ./cmd/prdtool-mcp
```
