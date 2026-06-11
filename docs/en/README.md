# GitFace

Git Multi-Identity & Remote Manager — Switch identity, git provider, and SSH key with one key. TUI + CLI dual interface.

[Download Latest Release](https://github.com/Fish-Ring/GitFace/releases/latest)

---

## Features

- **Profile Switching** — One-key switch `user.name`/`user.email` + remote URL + SSH key
- **Hosting Providers** — Built-in GitHub / Gitee / GitLab / Bitbucket, freely add your own
- **SSH Key Auto-Injection** — Automatically sets `core.sshCommand` on switch, no `~/.ssh/config` required
- **SSH Key Scanning** — `[Ctrl+O]` scan `~/.ssh/` directory and fill key path instantly
- **Conventional Commit** — Interactive commit type selection with auto add/commit/push
- **Dual Interface** — TUI + CLI subcommands
- **Mouse Scroll** — All lists support scroll wheel navigation
- **Multi-language** — Auto-detected from env or config (zh-CN / en)
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

| Document | Link |
|---|---|
| Install | [Install Guide](install.md) |
| Usage | [Usage Guide](usage.md) |
| Config | [Config Reference](config.md) |

---

## License

MIT
