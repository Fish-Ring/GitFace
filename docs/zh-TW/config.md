# 配置參考

## 配置檔案路徑

- **預設位置**: `~/.config/gitface/config.json`
- 首次執行自動產生

---

## 結構

```json
{
  "lang": "zh-TW",
  "active_profile_id": "",
  "profiles": [
    {
      "id": "personal",
      "name": "個人專案(GitHub)",
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

### 頂層欄位

| 欄位 | 類型 | 說明 |
|---|---|---|
| `lang` | string | 語言：`"zh-TW"`、`"zh-CN"`、`"en"`、`"ja"` 或 `"ko"` |
| `active_profile_id` | string | 預留 |
| `profiles` | array | 身分列表 |
| `providers` | array | 提供者列表（預設: GitHub, Gitee, GitLab, Bitbucket） |

### Profile 欄位

| 欄位 | 說明 |
|---|---|
| `id` | 唯一識別碼，CLI `switch <id>` 使用 |
| `name` | 顯示名稱（TUI 中顯示） |
| `git_name` | 寫入 `git config user.name` |
| `git_email` | 寫入 `git config user.email` |
| `provider_id` | 關聯的提供者 ID；切身分時會自動將遠端位址主機替換為該提供者的主機網域 |
| `ssh_identity_file` | SSH 私鑰路徑；切身分時自動注入 `git config core.sshCommand` |

### Provider 欄位

| 欄位 | 說明 |
|---|---|
| `id` | 唯一識別碼（例如 `github`） |
| `name` | 顯示名稱（例如 `GitHub`） |
| `host` | Git 主機網域（例如 `github.com`、`gitee.com`） |

---

## 完整範例

```json
{
  "lang": "zh-TW",
  "profiles": [
    {
      "id": "personal",
      "name": "個人專案(GitHub)",
      "git_name": "iol",
      "git_email": "iol@example.com",
      "provider_id": "github",
      "ssh_identity_file": "~/.ssh/id_ed25519_personal"
    },
    {
      "id": "work",
      "name": "公司專案(GitLab)",
      "git_name": "iol-work",
      "git_email": "iol@company.com",
      "provider_id": "gitlab",
      "ssh_identity_file": "~/.ssh/id_ed25519_company"
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

## SSH 設定

### 自動設定（推薦）
GitFace 透過 `git config core.sshCommand "ssh -i <金鑰路徑>"` 自動注入 SSH 金鑰，**無需手動設定 `~/.ssh/config`**。切換身分時自動注入，切走時自動清除。

### 手動 ~/.ssh/config（可選）
如果你希望用傳統方式，也可以在 `~/.ssh/config` 中設定：

```
Host github-personal
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_personal
```

但注意：GitFace 的提供者系統會改寫遠端位址主機為真實網域（如 `github.com`），而非別名。金鑰選擇由 `core.sshCommand` 自動處理，無需 SSH config 別名。

---

## 切換身分時執行的操作

切換到某個身分時，GitFace 自動執行：

1. `git config user.name` ← 身分的 `git_name`
2. `git config user.email` ← 身分的 `git_email`
3. 遠端位址主機名稱替換為提供者的主機網域
4. `git config core.sshCommand "ssh -i <ssh_identity_file>"` ← 金鑰自動注入

---

> 首頁：[README](../../README.md)
> 安裝指南：[安裝指南](install.md)
> 使用指南：[使用指南](usage.md)
