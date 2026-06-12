> [简体中文](../../README.md) · [繁體中文](../zh-TW/README.md) · [日本語](../ja/README.md) · [한국어](../ko/README.md)

# GitFace

Git Multi-Identity & Remote Manager — Switch identity, git provider, and SSH key with one key. TUI + CLI dual interface.

[Download Latest Release](https://github.com/Fish-Ring/GitFace/releases/latest)

---

## Features

- **Profile Switching** — One-key switch `user.name`/`user.email` + remote URL + SSH key
- **Hosting Providers** — Built-in GitHub / Gitee / GitLab / Bitbucket, freely add your own
- **SSH Key Auto-Injection** — Automatically sets `core.sshCommand` on switch, no `~/.ssh/config` required
- **SSH Key Scanning** — `[Ctrl+O]` scan `~/.ssh/` directory and fill key path instantly
- **Repo Config** — Per-repo remote path mappings, edit via `[G]` in TUI
- **Conventional Commit** — Interactive commit type selection with auto add/commit/push
- **Release Tag** — Create and push release tags from TUI after commit or via CLI
- **Dual Interface** — TUI + CLI subcommands
- **Mouse Scroll** — All lists support scroll wheel navigation
- **Multi-language** — Auto-detected from env or config (zh-CN / zh-TW / en / ja / ko)
- **Zero-config** — Auto-generates config on first run

---

## Quick Install

```powershell
# Build from source (requires Go)
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

> Platform-specific install: [Install Guide](install.md)

---

## Quick Start

```bash
# Run inside any Git repository
gitf                # Launch TUI
gitf status         # Show git status
gitf switch work    # Quick switch profile
```

---

## Documentation

[Install](install.md) · [Usage](usage.md) · [Config](config.md)

---

## License

MIT
