# 使用指南

## CLI 命令

| 命令 | 說明 |
|---|---|
| `gitf` | 啟動 TUI 互動介面 |
| `gitf tui` | 同上 |
| `gitf switch <id>` | 按 ID 切換身分 |
| `gitf status` | 查看目前 Git 身分與遠端 |
| `gitf tag <version>` | 建立並推送發布標籤 |
| `gitf edit` | 用編輯器開啟配置檔案 |
| `gitf help` | 顯示說明資訊 |
| `gitf version` | 顯示版本號 |

### CLI 範例

```bash
gitf status
gitf switch personal
gitf edit
gitf tag v1.0.0
```

---

## TUI 互動介面

執行 `gitf`（不加參數）進入 TUI。

### 介面佈局

```
┌──────────────────────────────────────────────┐
│ GitFace v1.0  │  main [5]                    │
├──────────────────────────────────────────────┤

  【目前狀態】
  • 分支:  main (2 檔案未暫存)
  • 遠端:  git@github.com:user/repo.git
  • 身分:  iol <iol@personal.com> [Personal]

  ──────────────────────────────────────────────

  【快捷操作】
  [1] Personal (GitHub)
  [2] Company (GitLab)
  [C] 提交
  [A] 帳戶管理
  [P] 提供者管理
  [S] 設定

  [Y]複製  [R]重新整理  [Q]離開
```

### 主介面快捷鍵

| 按鍵 | 功能 |
|---|---|
| `[1]` `[2]` ... | 切換對應編號的身分 |
| `[C]` | 規範化提交 |
| `[A]` | 帳戶管理 |
| `[P]` | 提供者管理 |
| `[S]` | 設定（切換語言 + 編輯設定） |
| `[Y]` | 複製身分資訊到剪貼簿 |
| `[R]` | 重新整理狀態 |
| `[Q]` | 離開 |
| ↑ ↓ / 滾輪 | 導覽 |

### 帳戶管理

按 `[A]` 進入帳戶管理：

| 按鍵 | 功能 |
|---|---|
| `[N]` | 新增身分 |
| `[D]` | 刪除模式（再按 [Enter] 確認刪除） |
| `[Enter]` | 編輯選中身分 |
| `[Esc]` | 返回 |
| ↑ ↓ / 滾輪 | 導覽 |

### 帳戶表單（6 個欄位）

| 欄位 | 說明 |
|---|---|
| ID | 唯一識別碼，用於 `switch <id>` |
| Name | 顯示名稱 |
| Git Name | `git config user.name` |
| Git Email | `git config user.email` |
| 提供者 | 托管提供者 — 直接輸入名稱或按 `[Ctrl+P]` 選擇 |
| SSH 金鑰 | 私鑰路徑 — 直接輸入路徑或按 `[Ctrl+O]` 掃描 |

### 提供者管理

按 `[P]` 進入提供者管理：

| 按鍵 | 功能 |
|---|---|
| `[N]` | 新增提供者 |
| `[D]` | 刪除模式 |
| `[Enter]` | 編輯選中提供者 |
| `[Esc]` | 返回 |

內建提供者：GitHub / Gitee / GitLab / Bitbucket。可自由增刪改。

### 設定

按 `[S]` 進入設定：

- **語言** — ← → / Tab 切換，Enter 確認
- **編輯設定** — 用系統編輯器開啟配置檔案（儲存後自動重新載入）

### SSH 金鑰掃描

1. 在帳戶表單中聚焦 SSH Key 欄位
2. 按 `[Ctrl+O]`
3. 用 ↑↓ 或滾輪選擇金鑰
4. 按 `[Enter]` 填入路徑

### 提供者選擇器

1. 在帳戶表單中聚焦 Provider 欄位
2. 按 `[Ctrl+P]`
3. 用 ↑↓ 或滾輪選擇提供者
4. 按 `[Enter]` 填入名稱

### 非 Git 儲存庫

在非 Git 目錄中啟動時，僅可管理提供者、帳戶、設定和離開。

---

## 規範化提交流程

1. 按 `[C]` 進入提交類型選擇
2. 選擇類型：
   - `f` → `feat`（新功能）
   - `x` → `fix`（修復）
   - `d` → `docs`（文件）
   - `r` → `refactor`（重構）
3. 輸入提交描述，按 `Enter` 確認
4. 自動執行 `git add .` → `git commit -m "type: desc"` → `git push`
5. 在「提交結果」介面檢視輸出，按任意鍵返回
6. 提交成功後，系統會詢問是否建立發布標籤（是/否）
7. 選擇「是」，輸入版本號如 `v1.0.0`，標籤將被建立並推送

---

## 多語言

自動檢測 `LANG` 環境變數切換語言：

```bash
export LANG=zh_TW.UTF-8   # 繁體中文
export LANG=zh_CN.UTF-8   # 简体中文
export LANG=en_US.UTF-8   # English
export LANG=ja_JP.UTF-8   # 日本語
export LANG=ko_KR.UTF-8   # 한국어
```

也可在配置中強制指定：

```json
{ "lang": "zh-TW" }   // 繁體中文
{ "lang": "zh-CN" }   // 简体中文
{ "lang": "en" }      // English
{ "lang": "ja" }      // 日本語
{ "lang": "ko" }      // 한국어
```

優先序：配置檔案 > `LANG` > `LC_ALL` > `LANGUAGE` > English

---

> 首頁：[README](../../README.md)
> 安裝指南：[安裝指南](install.md)
> 配置參考：[配置參考](config.md)
