> [简体中文](../../README.md) · [English](../en/README.md) · [日本語](../ja/README.md) · [한국어](../ko/README.md)

# GitFace

Git 多身份與遠端管理器 — 一鍵切換身份、託管提供者和 SSH 金鑰，支援 TUI 和 CLI 雙介面。

[下載最新版](https://github.com/Fish-Ring/GitFace/releases/latest)

---

## 特色

- **身份切換** — 一鍵切換 `user.name`/`user.email` + 遠端倉庫位址 + SSH 金鑰
- **託管提供者** — 內建 GitHub / Gitee / GitLab / Bitbucket，自由新增
- **SSH 金鑰自動注入** — 切換身份時自動配置 `core.sshCommand`，無需手動配 `~/.ssh/config`
- **SSH 金鑰掃描** — `[Ctrl+O]` 掃描 `~/.ssh/` 目錄，一鍵填入金鑰路徑
- **規範化提交** — 互動式 Conventional Commit（feat/fix/docs/refactor），自動 add/commit/push
- **發布標籤** — 提交後在 TUI 中建立並推送發布標籤，或透過 CLI 操作
- **雙介面** — TUI 互動介面 + CLI 子命令
- **鼠標滾輪** — 所有列表支援滾輪導航
- **多語言** — 簡體中文 / 繁體中文 / English / 日本語 / 한국어 自動切換
- **零配置啟動** — 首次運行自動生成配置檔

---

## 快速安裝

```powershell
# 從原始碼構建（需要 Go）
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

> 各平台詳細安裝：[安裝文檔](install.md)

---

## 快速開始

```bash
# 進入任意 Git 倉庫
gitf                # 啟動 TUI
gitf status         # 查看 Git 狀態
gitf switch personal  # 快速切換身份
```

---

## 文檔

[安裝](install.md) · [使用](usage.md) · [配置](config.md)

---

## 許可證

MIT
