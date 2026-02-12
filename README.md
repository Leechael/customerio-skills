# cio

CLI for the [Customer.io App API](https://customer.io/docs/api/app/).

## Install

Download a binary from [Releases](https://github.com/Leechael/customerio-skills/releases), or build from source:

```bash
make build    # → bin/cio
```

## Setup

```bash
export CUSTOMERIO_API_TOKEN="your-app-api-key"
```

## Usage

```bash
cio segments ls
cio customers get user123
cio campaigns ls --jq '.campaigns[0].name'
cio send email --body '{"to":"u@e.com","transactional_message_id":"1"}'
```

Use `--region eu` for the EU datacenter. Use `--jq` to filter JSON output.

Run `cio --help` for all commands, or `cio <command> --help` for subcommand details.

## Available Commands

| Command | Description |
|---------|-------------|
| `customers` | Manage customers |
| `segments` | Manage segments |
| `campaigns` | Manage campaigns |
| `broadcasts` | Manage broadcasts |
| `newsletters` | Manage newsletters |
| `transactional` | Manage transactional messages |
| `send` | Send email, push, SMS |
| `collections` | Manage collections |
| `exports` | Manage exports |
| `objects` | Manage objects |
| `messages` | Manage messages |
| `webhooks` | Manage reporting webhooks |
| `sender-identities` | Manage sender identities |
| `snippets` | Manage snippets |
| `esp-suppression` | Manage ESP suppressions |
| `imports` | Manage imports |
| `info` | General information |
| `workspaces` | Manage workspaces |

## Development

```bash
make test           # Unit tests + BDD tests
make test-unit      # Unit tests only
make lint           # golangci-lint
make all            # Cross-compile all platforms → dist/
make help           # Show all targets
```

## Claude Code Skill

This repo includes a [Claude Code skill](skills/customerio/SKILL.md) for AI-assisted Customer.io management. Install with:

```bash
npx skills add Leechael/customerio-skills
```

## Release

Maintainers can trigger a release by commenting on any Issue or PR:

```
!release patch    # v0.0.1
!release minor    # v0.1.0
!release major    # v1.0.0
```

Or use the **Release Command** workflow dispatch in GitHub Actions.
