# 使用指南

## CLI 命令

| 命令 | 说明 |
|---|---|
| `gitf` | 启动 TUI 交互界面 |
| `gitf tui` | 同上 |
| `gitf switch <id>` | 按 ID 切换身份 |
| `gitf status` | 查看当前 Git 身份与远程仓库 |
| `gitf tag <version>` | 创建并推送发布标签 |
| `gitf edit` | 用编辑器打开配置文件 |
| `gitf help` | 显示帮助信息 |
| `gitf version` | 显示版本号 |

### CLI 示例

```bash
# 查看状态
gitf status

# 快速切换身份
gitf switch personal
gitf switch company

# 编辑配置
gitf edit

# 创建发布标签
gitf tag v1.0.0
```

---

## TUI 交互界面

运行 `gitf`（不加参数）进入 TUI。

### 界面布局

```
┌──────────────────────────────────────────────┐
│ GitFace v1.0  │  main [5]                    │
├──────────────────────────────────────────────┤

  【当前状态】
  • 分支:  main (2 文件未暂存)
  • 远程仓库:  git@github.com:user/repo.git
  • 身份:  iol <iol@personal.com> [Personal]

  ──────────────────────────────────────────────

  【快捷操作】
  [1] Personal (GitHub)
  [2] Company (GitLab)
  [C] 提交
  [A] 账户管理
  [P] 提供商
  [S] 设置

  [Y]复制  [R]刷新  [Q]退出
```

### 主界面快捷键

| 按键 | 功能 |
|---|---|
| `[1]` `[2]` ... | 切换对应编号的身份 |
| `[C]` | 规范化提交 |
| `[A]` | 账户管理 |
| `[P]` | 提供商管理 |
| `[S]` | 设置（切换语言 + 编辑配置） |
| `[Y]` | 复制身份信息到剪贴板 |
| `[R]` | 刷新状态 |
| `[Q]` | 退出 |
| ↑ ↓ / 滚轮 | 导航 |

### 账户管理

按 `[A]` 进入账户管理：

| 按键 | 功能 |
|---|---|
| `[N]` | 新建身份 |
| `[D]` | 删除模式（再按 [Enter] 确认删除） |
| `[Enter]` | 编辑选中身份 |
| `[Esc]` | 返回 |
| ↑ ↓ / 滚轮 | 导航 |

### 账户表单（6 个字段）

| 字段 | 说明 |
|---|---|
| ID | 唯一标识符，用于 `switch <id>` |
| Name | 显示名称 |
| Git Name | `git config user.name` |
| Git Email | `git config user.email` |
| 提供商 | 托管提供商 — 直接输入名称或按 `[Ctrl+P]` 选择 |
| SSH 密钥 | 私钥路径 — 直接输入路径或按 `[Ctrl+O]` 扫描 |

### 提供商管理

按 `[P]` 进入提供商管理：

| 按键 | 功能 |
|---|---|
| `[N]` | 新建提供商 |
| `[D]` | 删除模式 |
| `[Enter]` | 编辑选中提供商 |
| `[Esc]` | 返回 |

内置提供商：GitHub / Gitee / GitLab / Bitbucket。可自由增删改。

### 设置

按 `[S]` 进入设置：

- **语言** — ← → / Tab 切换，Enter 确认
- **编辑配置** — 用系统编辑器打开配置文件（保存后自动重载）

### SSH 密钥扫描

1. 在账户表单中聚焦 SSH Key 字段
2. 按 `[Ctrl+O]`
3. 用 ↑↓ 或滚轮选择密钥
4. 按 `[Enter]` 填入路径

### 提供商选择器

1. 在账户表单中聚焦 Provider 字段
2. 按 `[Ctrl+P]`
3. 用 ↑↓ 或滚轮选择提供商
4. 按 `[Enter]` 填入名称

### 非 Git 仓库

在非 Git 目录中启动时，仅可管理提供商、账户、设置和退出。

---

## 规范化提交流程

1. 按 `[C]` 进入提交类型选择
2. 选择类型：
   - `f` → `feat`（新功能）
   - `x` → `fix`（修复）
   - `d` → `docs`（文档）
   - `r` → `refactor`（重构）
3. 输入提交描述，按 `Enter` 确认
4. 自动执行 `git add .` → `git commit -m "type: desc"` → `git push`
5. 在「提交结果」界面查看输出，按任意键返回
6. 提交成功后，系统会询问是否创建发布标签（是/否）
7. 选择「是」，输入版本号如 `v1.0.0`，标签将被创建并推送

---

## 多语言

自动检测 `LANG` 环境变量切换语言：

```bash
export LANG=zh_CN.UTF-8   # 简体中文
export LANG=zh_TW.UTF-8   # 繁體中文
export LANG=en_US.UTF-8   # English
export LANG=ja_JP.UTF-8   # 日本語
export LANG=ko_KR.UTF-8   # 한국어
```

也可在配置中强制指定：

```json
{ "lang": "zh-CN" }   // 简体中文
{ "lang": "zh-TW" }   // 繁體中文
{ "lang": "en" }      // English
{ "lang": "ja" }      // 日本語
{ "lang": "ko" }      // 한국어
```

优先级：配置文件 > `LANG` > `LC_ALL` > `LANGUAGE` > 英文

---

> 首页：[首页](../../README.md)  
> 安装指南：[安装指南](install.md)  
> 配置参考：[配置参考](config.md)
