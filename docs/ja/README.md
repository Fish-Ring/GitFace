> [简体中文](../../README.md) · [English](../en/README.md) · [繁體中文](../zh-TW/README.md) · [한국어](../ko/README.md)

# GitFace

Git 多ID・リモート管理ツール — ID、ホスティングプロバイダ、SSH鍵をワンキーで切替。TUI + CLI デュアルインターフェース。

[最新版をダウンロード](https://github.com/Fish-Ring/GitFace/releases/latest)

---

## 機能

- **プロフィール切替** — `user.name`/`user.email` + リモートリポジトリURL + SSH鍵をワンキーで切替
- **ホスティングプロバイダ** — GitHub / Gitee / GitLab / Bitbucket を内蔵、自由に追加可能
- **SSH鍵自動注入** — 切替時に`core.sshCommand`を自動設定、`~/.ssh/config`不要
- **SSH鍵スキャン** — `[Ctrl+O]`で`~/.ssh/`ディレクトリをスキャンし、鍵パスを即入力
- **リポジトリ設定** — リポジトリごとのリモートパスマッピング、TUIで`[G]`で編集
- **定型コミット** — インタラクティブな Conventional Commit（feat/fix/docs/refactor）、自動 add/commit/push
- **リリースタグ** — TUIで `[T]` を押してリリースタグを作成・プッシュ、またはCLIで操作
- **ブランチ切替** — TUIで `[B]` を押してローカルブランチを切替
- **デュアルインターフェース** — TUI インタラクティブ + CLI サブコマンド
- **マウススクロール** — 全リストでスクロールホイール対応
- **多言語** — 简体中文 / 繁體中文 / English / 日本語 / 한국어 自動切替
- **ゼロ設定** — 初回実行で設定ファイルを自動生成

---

## クイックインストール

```powershell
# ソースからビルド（Goが必要）
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

> プラットフォーム別詳細：[インストールガイド](install.md)

---

## クイックスタート

```bash
# 任意のGitリポジトリで実行
gitf                # TUIを起動
gitf status         # Gitステータスを表示
gitf switch personal  # プロフィールを素早く切替
```

---

## ドキュメント

[インストール](install.md) · [使い方](usage.md) · [設定](config.md)

---

## ライセンス

MIT
