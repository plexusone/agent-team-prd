# Product Requirements Agent Team

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/agent-team-prd/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/agent-team-prd/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/agent-team-prd/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/agent-team-prd/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/agent-team-prd/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/agent-team-prd/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/agent-team-prd
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/agent-team-prd
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/agent-team-prd
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/agent-team-prd
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fagent-team-prd
 [loc-svg]: https://tokei.rs/b1/github/plexusone/agent-team-prd
 [repo-url]: https://github.com/plexusone/agent-team-prd
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/agent-team-prd/blob/master/LICENSE

CLI and AI assistant integrations for creating, validating, scoring, and managing Product Requirements Documents.

## Features

- **Structured PRDs**: JSON-based schema for consistent, machine-readable documents
- **Quality Scoring**: Automated rubric-based scoring across 10 categories
- **Multiple Views**: PM-focused and Executive summary views
- **AI Integration**: Claude Code (MCP) and Kiro IDE (Power) support
- **Validation**: Schema validation, ID format checking, traceability verification

## Quick Start

```bash
# Install
go install github.com/agentplexus/agent-team-prd/cmd/prdtool@latest

# Create a PRD
prdtool init --title "User Authentication" --owner "Jane PM"
prdtool add problem --statement "Users cannot securely access accounts"
prdtool add goal --statement "Reduce password tickets by 50%"
prdtool add req --description "Support OAuth 2.0 login" --priority must

# Validate and score
prdtool validate
prdtool score
```

## AI Assistant Integration

```bash
# Claude Code
prdtool deploy --target claude

# Kiro IDE
prdtool deploy --target kiro-power --output ~/.kiro/powers/prdtool
```

## Documentation

Full documentation: https://agentplexus.github.io/agent-team-prd/

- [Installation](https://agentplexus.github.io/agent-team-prd/installation/)
- [CLI Reference](https://agentplexus.github.io/agent-team-prd/cli/commands/)
- [Claude Code Integration](https://agentplexus.github.io/agent-team-prd/integrations/claude-code/)
- [Kiro IDE Integration](https://agentplexus.github.io/agent-team-prd/integrations/kiro-ide/)
- [PRD Schema Reference](https://agentplexus.github.io/agent-team-prd/reference/prd-schema/)

## Build from Source

```bash
git clone https://github.com/agentplexus/agent-team-prd.git
cd agent-team-prd
go build -o bin/prdtool ./cmd/prdtool
go build -o bin/prdtool-mcp ./cmd/prdtool-mcp
```

## License

MIT
