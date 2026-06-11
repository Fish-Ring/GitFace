> [English](docs/en/README.md) · [繁體中文](docs/zh-TW/README.md) · [日本語](docs/ja/README.md) · [한국어](docs/ko/README.md)

# GitFace

Git 多身份与远程管理器 — 一键切换身份、托管提供商和 SSH 密钥，支持 TUI 和 CLI 双界面。

[下载最新版](https://github.com/Fish-Ring/GitFace/releases/latest)

---

## 特性

- **身份切换** — 一键切换 `user.name`/`user.email` + 远程仓库地址 + SSH 密钥
- **托管提供商** — 内置 GitHub / Gitee / GitLab / Bitbucket，自由添加
- **SSH 密钥自动注入** — 切身份时自动配置 `core.sshCommand`，无需手动配 `~/.ssh/config`
- **SSH 密钥扫描** — `[Ctrl+O]` 扫描 `~/.ssh/` 目录，一键填入密钥路径
- **规范化提交** — 交互式 Conventional Commit（feat/fix/docs/refactor），自动 add/commit/push
- **发布标签** — 提交后在 TUI 中创建并推送发布标签，或通过 CLI 操作
- **双界面** — TUI 交互界面 + CLI 子命令
- **鼠标滚轮** — 所有列表支持滚轮导航
- **多语言** — 简体中文 / 繁體中文 / English / 日本語 / 한국어 自动切换
- **零配置启动** — 首次运行自动生成配置文件

---

## 快速安装

```powershell
# 从源码构建（需要 Go）
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

> 各平台详细安装：[安装文档](docs/zh-CN/install.md)

---

## 快速开始

```bash
# 进入任意 Git 仓库
gitf                # 启动 TUI
gitf status         # 查看 Git 状态
gitf switch personal  # 快速切换身份
```

---

## 文档

[安装](docs/zh-CN/install.md) · [使用](docs/zh-CN/usage.md) · [配置](docs/zh-CN/config.md)

---

## 许可证

MIT
