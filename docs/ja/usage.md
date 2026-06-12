# 使い方ガイド

## CLI コマンド

| コマンド | 説明 |
|---|---|
| `gitf` | TUIを起動 |
| `gitf tui` | TUIを起動（明示的） |
| `gitf switch <id>` | IDでプロファイルを切替 |
| `gitf status` | 現在の身分とリモートを表示 |
| `gitf tag <version>` | リリースタグを作成してプッシュ |
| `gitf edit` | エディタで設定を開く |
| `gitf help` | ヘルプを表示 |
| `gitf version` | バージョンを表示 |

### CLI サンプル

```bash
gitf status
gitf switch personal
gitf edit
gitf tag v1.0.0
```

---

## TUI インターフェース

`gitf`（引数なし）でTUIを起動。

### レイアウト

```
┌──────────────────────────────────────────────┐
│ GitFace v1.0  │  main [5]                    │
├──────────────────────────────────────────────┤

  【ステータス】
  • ブランチ:  main (2 ファイル未ステージ)
  • リモート:  git@github.com:user/repo.git
  • 身份:  iol <iol@personal.com> [Personal]

  ──────────────────────────────────────────────

  【操作】
  [1] Personal (GitHub)
  [2] Company (GitLab)
  [C] コミット
  [A] アカウント管理
  [P] プロバイダー管理
  [S] 設定

  [Y]コピー  [R]更新  [Q]終了
```

### ホットキー（メイン画面）

| キー | アクション |
|---|---|
| `[1]` `[2]` ... | 番号で身分を切替 |
| `[C]` | 規範コミット |
| `[A]` | アカウント管理 |
| `[P]` | プロバイダー管理 |
| `[S]` | 設定（言語切替 + 設定編集） |
| `[Y]` | 身分情報をクリップボードにコピー |
| `[R]` | ステータス更新 |
| `[Q]` | 終了 |
| ↑ ↓ / スクロール | ナビゲート |

### アカウント管理

`[A]` でアカウント管理：

| キー | アクション |
|---|---|
| `[N]` | 新規アカウント |
| `[D]` | 削除モード（[Enter] で確認） |
| `[Enter]` | 選択したアカウントを編集 |
| `[Esc]` | 戻る |
| ↑ ↓ / スクロール | ナビゲート |

### アカウントフォーム（7 フィールド）

| フィールド | 説明 |
|---|---|
| ID | `switch <id>` 用のユニークID |
| Name | 表示名 |
| Git Name | `git config user.name` |
| Git Email | `git config user.email` |
| Provider | ホスティングプロバイダー — 入力または `[Ctrl+P]` で選択 |
| SSH Key | 私密鍵パス — 入力または `[Ctrl+O]` でスキャン |
| リモートパス | パス上書き: `host:path=newpath`（カンマ区切り） |

### プロバイダー管理

`[P]` でプロバイダー管理：

| キー | アクション |
|---|---|
| `[N]` | 新規プロバイダー |
| `[D]` | 削除モード |
| `[Enter]` | 選択したプロバイダーを編集 |
| `[Esc]` | 戻る |

デフォルトプロバイダー：GitHub / Gitee / GitLab / Bitbucket。

### 設定

`[S]` で設定：

- **言語** — ← → / Tab で切替、Enter で確定
- **設定編集** — エディタで設定ファイルを開く（保存後自動リロード）

### SSH鍵スキャン

1. アカウントフォームのSSH Key フィールドにフォーカス
2. `[Ctrl+O]` を押す
3. ↑↓ またはスクロールで鍵を選択
4. `[Enter]` でパスを入力

### プロバイダー選択器

1. アカウントフォームのProvider フィールドにフォーカス
2. `[Ctrl+P]` を押す
3. ↑↓ またはスクロールでプロバイダーを選択
4. `[Enter]` で名前を入力

---

## 規範コミットフロー

1. `[C]` でコミットタイプを選択
2. タイプを選択：
   - `f` → `feat`（新機能）
   - `x` → `fix`（バグ修正）
   - `d` → `docs`（ドキュメント）
   - `r` → `refactor`（リファクタリング）
3. 説明を入力、`Enter` で確定
4. 自動 `git add .` → `git commit -m "type: desc"` → `git push`
5. 結果を確認、任意キーで戻る
6. コミット成功後、リリースタグの作成確認が出ます（はい/いいえ）
7. 「はい」を選択し、`v1.0.0` のようにバージョンを入力すると、タグが作成されてプッシュされます

---

## 多言語

`LANG` 環境変数から自動検出：

```bash
export LANG=ja_JP.UTF-8   # 日本語
export LANG=en_US.UTF-8   # English
export LANG=zh_CN.UTF-8   # 中文
export LANG=ko_KR.UTF-8   # 한국어
```

設定で強制指定：

```json
{ "lang": "ja" }      // 日本語
{ "lang": "en" }      // English
{ "lang": "zh-CN" }   // 中文
{ "lang": "ko" }      // 한국어
```

優先順位：設定ファイル > `LANG` > `LC_ALL` > `LANGUAGE` > English

---

> ホーム: [README](README.md)
> インストール: [Install Guide](install.md)
> 設定: [Config Reference](config.md)
