# 配置参考

## 配置文件路径

- **默认位置**: `~/.config/gitface/config.json`
- 首次运行自动生成

---

## 结构

```json
{
  "lang": "zh-CN",
  "active_profile_id": "",
  "profiles": [
    {
      "id": "personal",
      "name": "个人项目(GitHub)",
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

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `lang` | string | 语言：`"zh-CN"`、`"zh-TW"`、`"en"`、`"ja"` 或 `"ko"` |
| `active_profile_id` | string | 预留 |
| `profiles` | array | 身份列表 |
| `providers` | array | 提供商列表（默认: GitHub, Gitee, GitLab, Bitbucket） |

### Profile 字段

| 字段 | 说明 |
|---|---|
| `id` | 唯一标识，CLI `switch <id>` 使用 |
| `name` | 显示名称（TUI 中显示） |
| `git_name` | 写入 `git config user.name` |
| `git_email` | 写入 `git config user.email` |
| `provider_id` | 关联的提供商 ID；切身份时会自动将远程地址主机替换为该提供商的主机域名 |
| `ssh_identity_file` | SSH 私钥路径；切身份时自动注入 `git config core.sshCommand` |

### Provider 字段

| 字段 | 说明 |
|---|---|
| `id` | 唯一标识（例如 `github`） |
| `name` | 显示名称（例如 `GitHub`） |
| `host` | Git 主机域名（例如 `github.com`、`gitee.com`） |

---

## 完整示例

```json
{
  "lang": "zh-CN",
  "profiles": [
    {
      "id": "personal",
      "name": "个人项目(GitHub)",
      "git_name": "iol",
      "git_email": "iol@example.com",
      "provider_id": "github",
      "ssh_identity_file": "~/.ssh/id_ed25519_personal"
    },
    {
      "id": "work",
      "name": "公司项目(GitLab)",
      "git_name": "iol-work",
      "git_email": "iol@company.com",
      "provider_id": "gitlab",
      "ssh_identity_file": "~/.ssh/id_ed25519_company"
    },
    {
      "id": "client",
      "name": "客户项目(Bitbucket)",
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

## SSH 配置

### 自动配置（推荐）
GitFace 通过 `git config core.sshCommand "ssh -i <密钥路径>"` 自动注入 SSH 密钥，**无需手动配置 `~/.ssh/config`**。切换身份时自动注入，切走时自动清理。

### 手动 ~/.ssh/config（可选）
如果你希望用传统方式，也可以在 `~/.ssh/config` 中配置：

```
Host github-personal
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_personal
```

但注意：GitFace 的提供商系统会改写远程地址主机为真实域名（如 `github.com`），而非别名。密钥选择由 `core.sshCommand` 自动处理，无需 SSH config 别名。

---

## 切换身份时执行的操作

切换到某个身份时，GitFace 自动执行：

1. `git config user.name` ← 身份的 `git_name`
2. `git config user.email` ← 身份的 `git_email`
3. 远程地址主机名替换为提供商的主机域名
4. `git config core.sshCommand "ssh -i <ssh_identity_file>"` ← 密钥自动注入

---

> 首页：[首页](../../README.md)  
> 安装指南：[安装指南](install.md)  
> 使用指南：[使用指南](usage.md)
