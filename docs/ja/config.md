# 設定リファレンス

## 設定ファイルの場所

- **デフォルト**: `~/.config/gitface/config.json`
- 初回実行時に自動生成

---

## 構造

```json
{
  "lang": "ja",
  "active_profile_id": "",
  "profiles": [
    {
      "id": "personal",
      "name": "個人 (GitHub)",
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

### トップレベルフィールド

| フィールド | 型 | 説明 |
|---|---|---|
| `lang` | string | 言語: `"en"`、`"zh-CN"`、`"ja"`、`"ko"` |
| `active_profile_id` | string | 予約 |
| `profiles` | array | プロファイルリスト |
| `providers` | array | プロバイダーリスト（デフォルト: GitHub, Gitee, GitLab, Bitbucket） |

### Profile フィールド

| フィールド | 説明 |
|---|---|
| `id` | ユニークID。`switch <id>` で使用 |
| `name` | TUIの表示名 |
| `git_name` | `git config user.name` に書き込み |
| `git_email` | `git config user.email` に書き込み |
| `provider_id` | プロバイダーの `id` に関連付け。リモートURLのホストをプロバイダーの `host` に書き換え |
| `ssh_identity_file` | SSH秘密鍵のパス。`git config core.sshCommand` で自動注入 |

### Provider フィールド

| フィールド | 説明 |
|---|---|
| `id` | ユニークID（例: `github`） |
| `name` | 表示名（例: `GitHub`） |
| `host` | Gitホストドメイン（例: `github.com`、`gitee.com`） |

---

## リポジトリごとの設定

各リポジトリは `.git/gitf.json` にプロバイダーごとのリモートパスマッピングを保存できます。このファイルは初回身分切替時に自動作成され、git は無視します（`.git/` 内にあるため）。

### 構造

```json
{
  "paths": {
    "github": "Fish-Ring/GitFace",
    "gitee": "Ableand/git-face"
  }
}
```

### フィールド

| フィールド | 説明 |
|---|---|
| `paths` | プロバイダーIDからそのプロバイダー上のリポジトリパスへのマップ |

### 仕組み

身分切替時に、GitFace はリモートURLを書き換えます。リポジトリ設定は各プロバイダーに使用するパスを指定します。例如、GitHub リポジトリが `Fish-Ring/GitFace` で Gitee ミラーが `Ableand/git-face` の場合、ここで両方のパスを設定します。

TUI で `[G]` を押して編集するか、手動で `.git/gitf.json` を編集してください。

---

## SSH設定

### 自動設定（推奨）
GitFace は `git config core.sshCommand "ssh -i <鍵パス>"` でSSH鍵を自動注入します。**`~/.ssh/config` の手動設定は不要**です。身分切替時に自動注入され、切替解除時に自動クリアされます。

### 手動 ~/.ssh/config（オプション）
従来の方式で `~/.ssh/config` を設定することもできます：

```
Host github-personal
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_personal
```

ただし、GitFace のプロバイダーシステムはリモートホストを実際のドメイン（例: `github.com`）に書き換えます。鍵の選択は `core.sshCommand` で自動処理されます。

---

## 身分切替時の操作

身分に切り替えると、GitFace は以下を自動実行：

1. `git config user.name` ← プロファイルの `git_name`
2. `git config user.email` ← プロファイルの `git_email`
3. リモートURLのホスト名をプロバイダーのドメインに書き換え
4. `git config core.sshCommand "ssh -i <ssh_identity_file>"` ← 鍵の自動注入

---

> ホーム: [README](README.md)
> インストール: [Install Guide](install.md)
> 使い方: [Usage Guide](usage.md)
