# Usage Guide

## CLI Commands

| Command | Description |
|---|---|
| `gitf` | Launch TUI |
| `gitf tui` | Launch TUI (explicit) |
| `gitf switch <id>` | Switch to profile by ID |
| `gitf status` | Show current identity & remote |
| `gitf tag <version>` | Create and push a release tag |
| `gitf edit` | Open config in editor |
| `gitf help` | Show help |
| `gitf version` | Show version |

### CLI Examples

```bash
# Show repo status
gitf status

# Quick switch identity
gitf switch personal
gitf switch company

# Edit config
gitf edit

# Create a release tag
gitf tag v1.0.0
```

---

## TUI Interface

Run `gitf` (no arguments) to launch the TUI.

### Layout

```
┌──────────────────────────────────────────────┐
│ GitFace v1.0  │  main [5]                    │
├──────────────────────────────────────────────┤

  [Status]
  • Branch:  main (2 files unstaged)
  • Remote:  git@github.com:user/repo.git
  • Identity:  iol <iol@personal.com> [Personal]

  ──────────────────────────────────────────────

  [Actions]
  [1] Personal (GitHub)
  [2] Company (GitLab)
  [C] Commit
  [A] Accounts
  [P] Manage
  [S] Settings

  [Y]Copy  [R]Refresh  [Q]Quit
```

### Hotkeys (Main Screen)

| Key | Action |
|---|---|
| `[1]` `[2]` ... | Switch to profile by number |
| `[C]` | Start conventional commit |
| `[A]` | Account management |
| `[P]` | Provider management |
| `[S]` | Settings (language toggle + edit config) |
| `[Y]` | Copy identity info to clipboard |
| `[R]` | Refresh status |
| `[Q]` | Quit |
| ↑ ↓ / scroll wheel | Navigate |

### Account Management

Press `[A]` to manage profiles:

| Key | Action |
|---|---|
| `[N]` | New account |
| `[D]` | Delete mode (press [Enter] to confirm) |
| `[Enter]` | Edit selected account |
| `[Esc]` | Back |
| ↑ ↓ / scroll wheel | Navigate |

### Account Form (6 fields)

| Field | Description |
|---|---|
| ID | Unique identifier for `switch <id>` |
| Name | Display name |
| Git Name | `git config user.name` |
| Git Email | `git config user.email` |
| Provider | Hosting provider — type name or `[Ctrl+P]` to pick |
| SSH Key | Path to private key — type path or `[Ctrl+O]` to scan |

### Provider Management

Press `[P]` to manage providers:

| Key | Action |
|---|---|
| `[N]` | New provider |
| `[D]` | Delete mode |
| `[Enter]` | Edit selected provider |
| `[Esc]` | Back |

Default providers: GitHub, Gitee, GitLab, Bitbucket. You can add/modify/delete any.

### Settings

Press `[S]` to open settings:

- **Language** — ← → / Tab to toggle, Enter confirms
- **Edit Config** — Open config file in editor (auto-reloads on save)

### SSH Key Scanning

1. Focus the SSH Key field in account form
2. Press `[Ctrl+O]`
3. Select a key from the scan list with ↑↓ or scroll wheel
4. Press `[Enter]` to fill the path

### Provider Picker

1. Focus the Provider field in account form
2. Press `[Ctrl+P]`
3. Select a provider with ↑↓ or scroll wheel
4. Press `[Enter]` to fill the name

### Outside Git Repository

Only provider management, account management, settings, and quit are available.

---

## Conventional Commit Flow

1. Press `[C]` to select commit type
2. Choose type:
   - `f` → `feat` (new feature)
   - `x` → `fix` (bug fix)
   - `d` → `docs` (documentation)
   - `r` → `refactor` (code refactor)
3. Enter description, press `Enter`
4. Auto `git add .` → `git commit -m "type: desc"` → `git push`
5. View output, press any key to return
6. After successful commit, you are prompted to create a release tag (Yes/No)
7. Select Yes, enter a version like `v1.0.0`, and the tag is created and pushed

---

## Multi-language

Auto-detected from `LANG` env:

```bash
export LANG=en_US.UTF-8   # English
export LANG=zh_CN.UTF-8   # 中文 (简体)
export LANG=zh_TW.UTF-8   # 中文 (繁體)
export LANG=ja_JP.UTF-8   # 日本語
export LANG=ko_KR.UTF-8   # 한국어
```

Override in config:

```json
{ "lang": "en" }      // English
{ "lang": "zh-CN" }   // Chinese (Simplified)
{ "lang": "zh-TW" }   // Chinese (Traditional)
{ "lang": "ja" }      // Japanese
{ "lang": "ko" }      // Korean
```

Priority: config > `LANG` > `LC_ALL` > `LANGUAGE` > English

---

> Home: [README](README.md)  
> Install: [Install Guide](install.md)  
> Config: [Config Reference](config.md)
