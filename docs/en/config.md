# Configuration Reference

## Config File Location

- **Default**: `~/.config/gitface/config.json`
- Auto-generated on first run

---

## Structure

```json
{
  "lang": "en",
  "active_profile_id": "",
  "profiles": [
    {
      "id": "personal",
      "name": "Personal (GitHub)",
      "git_name": "yourname",
      "git_email": "you@personal.com",
      "provider_id": "github",
      "ssh_identity_file": "~/.ssh/id_ed25519_personal"
    }
  ],
  "providers": [
    {
      "id": "github",
      "name": "GitHub",
      "host": "github.com"
    }
  ]
}
```

### Top-level Fields

| Field | Type | Description |
|---|---|---|
| `lang` | string | Language: `"en"`, `"zh-CN"`, `"zh-TW"`, `"ja"`, or `"ko"` |
| `active_profile_id` | string | Reserved |
| `profiles` | array | Profile list |
| `providers` | array | Provider list (default: GitHub, Gitee, GitLab, Bitbucket) |

### Profile Fields

| Field | Description |
|---|---|
| `id` | Unique identifier, used in `switch <id>` |
| `name` | Display name in TUI |
| `git_name` | Written to `git config user.name` |
| `git_email` | Written to `git config user.email` |
| `provider_id` | Links to a provider's `id`; rewrites remote URL host to provider's `host` |
| `ssh_identity_file` | Path to SSH private key; auto-injects via `git config core.sshCommand` |

### Provider Fields

| Field | Description |
|---|---|
| `id` | Unique identifier (e.g. `github`) |
| `name` | Display name (e.g. `GitHub`) |
| `host` | Git host domain (e.g. `github.com`, `gitee.com`) |

---

## Full Example

```json
{
  "lang": "en",
  "profiles": [
    {
      "id": "personal",
      "name": "Personal (GitHub)",
      "git_name": "iol",
      "git_email": "iol@example.com",
      "provider_id": "github",
      "ssh_identity_file": "~/.ssh/id_ed25519_personal"
    },
    {
      "id": "work",
      "name": "Company (GitLab)",
      "git_name": "iol-work",
      "git_email": "iol@company.com",
      "provider_id": "gitlab",
      "ssh_identity_file": "~/.ssh/id_ed25519_company"
    },
    {
      "id": "client",
      "name": "Client (Bitbucket)",
      "git_name": "iol-dev",
      "git_email": "iol@client.com",
      "provider_id": "bitbucket",
      "ssh_identity_file": "~/.ssh/id_ed25519_client"
    }
  ],
  "providers": [
    { "id": "github", "name": "GitHub", "host": "github.com" },
    { "id": "gitee", "name": "Gitee", "host": "gitee.com" },
    { "id": "gitlab", "name": "GitLab", "host": "gitlab.com" },
    { "id": "bitbucket", "name": "Bitbucket", "host": "bitbucket.org" }
  ]
}
```

---

## Per-Repo Configuration

Each repository can have its own `.git/gitf.json` file that stores per-provider remote path mappings. This file is auto-created on first identity switch and is ignored by git (inside `.git/`).

### Structure

```json
{
  "paths": {
    "github": "Fish-Ring/GitFace",
    "gitee": "Ableand/git-face"
  }
}
```

### Fields

| Field | Description |
|---|---|
| `paths` | Map of provider ID to repo path on that provider |

### How It Works

When you switch identities, GitFace rewrites the remote URL. The repo config tells it what path to use for each provider. For example, if your GitHub repo is at `Fish-Ring/GitFace` but your Gitee mirror is at `Ableand/git-face`, you configure both paths here.

You can edit this from the TUI by pressing `[G]` on the main screen, or manually edit `.git/gitf.json`.

---

## SSH Setup

### Automatic (Recommended)
GitFace auto-injects the SSH key via `git config core.sshCommand`. No `~/.ssh/config` entries needed. When you switch to a profile with `ssh_identity_file` set, the key is used automatically.

### Manual via ~/.ssh/config (Optional)
If you prefer using `~/.ssh/config` directly, you can set it up traditionally:

```
Host github-personal
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_personal
```

But note: GitFace's provider system rewrites the remote host (e.g., to `github.com`), not to an alias. The `core.sshCommand` injection handles key selection instead of SSH config aliasing.

---

## How Switching Works

When you switch to a profile:

1. `git config user.name` ← profile's `git_name`
2. `git config user.email` ← profile's `git_email`
3. Remote URL host is rewritten to the provider's `host` domain
4. `git config core.sshCommand "ssh -i <ssh_identity_file>"` ← key injection

---

> Home: [README](README.md)  
> Install: [Install Guide](install.md)  
> Usage: [Usage Guide](usage.md)
