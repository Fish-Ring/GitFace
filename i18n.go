package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Lang string

const (
	EN   Lang = "en"
	ZH   Lang = "zh-CN"
	ZH_TW Lang = "zh-TW"
	JA   Lang = "ja"
	KO   Lang = "ko"
)

var LangList = []Lang{EN, ZH, ZH_TW, JA, KO}

var langKeys = map[Lang]string{
	EN: "en", ZH: "zh", ZH_TW: "zh_TW", JA: "ja", KO: "ko",
}

var supportedLangs = map[string]Lang{
	"en":    EN, "en-US": EN, "en_US": EN,
	"zh":    ZH, "zh-CN": ZH, "zh_CN": ZH, "zh-cn": ZH,
	"zh-TW": ZH_TW, "zh_TW": ZH_TW, "zh-tw": ZH_TW,
	"ja":    JA, "ja-JP": JA, "ja_JP": JA,
	"ko":    KO, "ko-KR": KO, "ko_KR": KO,
}

func DetectLang() Lang {
	for _, env := range []string{"LANG", "LC_ALL", "LANGUAGE"} {
		if v := os.Getenv(env); v != "" {
			parts := strings.SplitN(v, ".", 2)
			if lang, ok := supportedLangs[parts[0]]; ok {
				return lang
			}
		}
	}
	if runtime.GOOS == "windows" {
		if lang := detectWindowsLang(); lang != "" {
			return lang
		}
	}
	return EN
}

func detectWindowsLang() Lang {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"[System.Globalization.CultureInfo]::CurrentCulture.TwoLetterISOLanguageName").Output()
	if err != nil {
		return ""
	}
	code := strings.TrimSpace(string(out))
	if lang, ok := supportedLangs[code]; ok {
		return lang
	}
	return ""
}

type Translator struct {
	lang Lang
}

func NewTranslator(lang Lang) *Translator {
	if lang == "" {
		lang = EN
	}
	return &Translator{lang: lang}
}

func (t *Translator) Tr(key string, args ...interface{}) string {
	text := texts[t.lang][key]
	if text == "" {
		text = texts[EN][key]
	}
	if text == "" {
		return key
	}
	if len(args) > 0 {
		return fmt.Sprintf(text, args...)
	}
	return text
}

var texts = map[Lang]map[string]string{
	EN: {
		// App
		"app_title": "GitFace v1.0",

		// Status section
		"label_status":  "[Status]",
		"label_branch":  "Branch:",
		"label_route":   "Remote:",
		"label_identity": "Identity:",

		// Status values
		"status_dirty":    "%s (%d files unstaged)",
		"status_clean":    "(clean)",
		"status_unable":   "Unable to get",
		"status_no_remote": "No remote",
		"status_not_set":  "Not set",
		"status_custom":   "Custom/Unmatched",

		// Actions
		"label_actions": "[Actions]",

		// Messages
		"msg_not_repo":     "✗ Not in Git repository, actions disabled",
		"msg_checking":     "Checking repository status...",
		"msg_switch_fail":  "Switch failed: %v",
		"msg_switch_ok":    "✓ Profile switched successfully",
		"msg_edit_fail":    "Edit config failed: %v",
		"msg_reload_ok":    "✓ Config reloaded",
		"msg_cannot_commit": "Not in a Git repository, cannot commit",
		"msg_cannot_switch": "Not in a Git repository, cannot switch identity",
		"msg_desc_empty":   "Commit description cannot be empty",

		// SwitchProfile logs
		"sp_name_fail":   "Failed to set user.name: %s",
		"sp_email_fail":  "Failed to set user.email: %s",
		"sp_remote_same": "Remote URL unchanged",
		"sp_remote_fail": "Failed to set remote URL: %s",
		"sp_no_url":      "Provider %s has no remote URL",
		"sp_prov_not_found": "Provider ID '%s' not found, skipped",
		"sp_no_provider": "No provider configured, skipped",
		"sp_ssh_fail":    "Failed to set SSH command: %s",

		// Prompts
		"prompt_select_type":     "Select commit type:",
		"prompt_enter_desc_ph":   "Enter commit description...",
		"prompt_confirm_cancel":  "[Enter] Confirm  [Esc] Cancel",
		"prompt_cancel":          "[Esc] Cancel",
		"prompt_any_key":         "Press any key to return",

		// Commit
		"commit_result": "Commit Result",
		"commit_header": "%s: enter description",
		"cli_error":     "Error: %v",

		// Commit output
		"commit_git_add_fail":    "git add failed: %s",
		"commit_nothing_new":     "No new changes to commit",
		"commit_git_commit_fail": "git commit failed: %s",
		"commit_git_push_fail":   "git push failed: %s",
		"commit_push_done":       "Committed and pushed successfully",
		"commit_pushed_only":     "Pushed to remote",

		// Tag
		"tag_prompt":          "Create release tag? [y/N]",
		"tag_input_ph":        "Enter version (e.g. v1.0.0)...",
		"tag_result":          "Tag Result",
		"tag_exists":          "Tag '%s' already exists",
		"tag_success":         "Tag '%s' created and pushed",
		"tag_git_tag_fail":    "git tag failed: %s",
		"tag_git_push_fail":   "git push tag failed: %s",
		"tag_version_empty":   "Version cannot be empty",
		"tag_hint":            "[t] Tag",

		// Action names
		"action_provider_mgmt": "Manage",
		"action_commit":        "Commit",
		"action_accounts":      "Accounts",
		"action_settings":      "Settings",
		"action_copy":          "Copy",
		"action_refresh":       "Refresh",
		"action_quit":          "Quit",
		"action_new":           "New",
		"action_delete":        "Delete",
		"action_edit_btn":      "Edit",
		"action_select":        "Select",
		"action_back":          "Back",
		"action_confirm_del":   "Confirm Del",
		"action_nav":           "Navigate",

		// Settings
		"label_settings":       "Settings",
		"settings_language":    "Language",
		"settings_lang_en":     "English",
		"settings_lang_zh":     "简体中文",
		"settings_lang_ja":     "日本語",
		"settings_lang_ko":     "한국어",
		"settings_lang_zh_TW":  "繁體中文",
		"settings_edit_config": "Edit Config",
		"settings_saved":       "Settings saved",
		"settings_save_fail":   "Failed to save settings: %v",

		// Header
		"label_no_repo": "Not a repo",

		// Account management
		"label_accounts":        "Account Management",
		"label_no_accounts":     "No profiles configured",
		"account_new":           "New Account",
		"account_edit":          "Edit Account",
		"account_deleted":       "Deleted profile '%s'",
		"account_saved":         "Saved profile '%s'",
		"account_id_ph":         "Profile ID",
		"account_name_ph":       "Display Name",
		"account_gitname_ph":    "Git User Name",
		"account_gitemail_ph":   "Git Email",
		"account_provider_ph":   "Provider",
		"account_provider_desc": "Git hosting provider (e.g. GitHub, Gitee) — [Ctrl+P] to pick, type to search",
		"prompt_tab_next":       "[Tab/↑↓] Next  [Enter] Save  [Esc] Cancel",
		"account_id_desc":      "Unique ID for quick switching (e.g. work)",
		"account_name_desc":    "Display name (e.g. Work Account)",
		"account_gitname_desc": "git user.name value (e.g. Zhang San)",
		"account_gitemail_desc": "git user.email value (e.g. zhang@company.com)",
		"account_ssh_identity_ph": "SSH Identity File",
		"account_ssh_identity_desc": "Path to SSH private key (e.g. ~/.ssh/id_ed25519_work)",
		"ssh_scan_hint":  "[Ctrl+O] Scan SSH keys",
		"ssh_scan_title": "Select SSH Key",
		"ssh_scan_none":  "No SSH keys found in ~/.ssh/. Use ssh-keygen to generate one.",
		"provider_pick_hint": "[Ctrl+P] Pick provider",

		// Messages
		"msg_copied":         "Copied to clipboard",
		"msg_cannot_copy":    "Not in a Git repository, nothing to copy",
		"msg_fields_required": "ID, Git Name and Git Email are required",
		"msg_duplicate_id":   "ID '%s' is already in use",
		"msg_no_profile":    "No profile selected",

		// Provider management
		"label_providers":        "Providers",
		"label_no_providers":     "No providers configured",
		"provider_new":           "New Provider",
		"provider_edit":          "Edit Provider",
		"provider_saved":         "Provider saved",
		"provider_deleted":       "Deleted provider: %s",
		"provider_id_ph":         "Provider ID",
		"provider_name_ph":       "Provider Name",
		"provider_host_ph":       "Host Domain",
		"provider_id_desc":       "Unique ID (e.g. github)",
		"provider_name_desc":     "Display name (e.g. GitHub, Gitee)",
		"provider_host_desc":     "Git host domain (e.g. github.com, gitee.com)",
		"msg_provider_fields_required": "ID, Name and Host are required",

		// Provider picker
		"provider_pick_title":  "Select Provider",
		"provider_pick_none":   "No providers configured",

		// Not in repo message
		"msg_press_e": "Press E to edit config  |  Q to quit",

		// CLI help
		"help_title":    "GitFace - Git Multi-Identity & Remote Manager",
		"help_usage":    "Usage:",
		"help_cmd_tui":  "  gitf [tui]               Launch TUI (default)",
		"help_cmd_switch": "  gitf switch <id>         Switch to profile by ID",
		"help_cmd_status": "  gitf status              Show current git identity & remote",
		"help_cmd_tag":    "  gitf tag <version>       Create and push a release tag",
		"help_cmd_edit": "  gitf edit                Open config in editor",
		"help_cmd_help": "  gitf help                Show this help",

		// CLI output
		"cli_switch_ok":       "Switched to profile '%s'",
		"cli_switch_not_found": "Profile not found: %s",
		"cli_switch_detail":   "%s",
		"cli_status_header":   "GitFace Status",
		"cli_branch":          "Branch:  %s",
		"cli_remote":          "Remote:  %s",
		"cli_name":            "Name:    %s",
		"cli_email":           "Email:   %s",
		"cli_dirty":           "Status:  %d files dirty",
		"cli_clean":           "Status:  Clean",
		"cli_edit_opening":    "Opening config in editor...",
		"cli_no_repo":         "Not in a Git repository",
		"cli_version":         "GitFace v1.0",
	},

	ZH: {
		// App
		"app_title": "GitFace v1.0",

		// Status section
		"label_status":  "【当前状态】",
		"label_branch":  "分支:",
		"label_route":   "远程仓库:",
		"label_identity": "身份:",

		// Status values
		"status_dirty":    "%s (%d 文件未暂存)",
		"status_clean":    "(clean)",
		"status_unable":   "无法获取",
		"status_no_remote": "无远程仓库",
		"status_not_set":  "未设置",
		"status_custom":   "自定义/未匹配",

		// Actions
		"label_actions": "【快捷操作】",

		// Messages
		"msg_not_repo":     "✗ 当前目录不是 Git 仓库，操作已禁用",
		"msg_checking":     "检查仓库状态...",
		"msg_switch_fail":  "切换失败: %v",
		"msg_switch_ok":    "✓ 身份切换成功",
		"msg_edit_fail":    "编辑配置失败: %v",
		"msg_reload_ok":    "✓ 配置已重载",
		"msg_cannot_commit": "不在 Git 仓库中，无法提交",
		"msg_cannot_switch": "不在 Git 仓库中，无法切换身份",
		"msg_desc_empty":   "提交描述不能为空",

		// SwitchProfile logs
		"sp_name_fail":      "设置 user.name 失败: %s",
		"sp_email_fail":     "设置 user.email 失败: %s",
		"sp_remote_same":    "远程地址未变化",
		"sp_remote_fail":    "设置远程仓库地址失败: %s",
		"sp_no_url":         "提供商 %s 已关联，无 remote URL",
		"sp_prov_not_found": "未找到提供商 ID '%s'，跳过远程地址设置",
		"sp_no_provider":    "未关联提供商，跳过远程地址设置",
		"sp_ssh_fail":       "设置 SSH 命令失败: %s",

		// Prompts
		"prompt_select_type":     "选择提交类型:",
		"prompt_enter_desc_ph":   "输入提交描述...",
		"prompt_confirm_cancel":  "[Enter] 确认  [Esc] 取消",
		"prompt_cancel":          "[Esc] 取消",
		"prompt_any_key":         "按任意键返回",

		// Commit
		"commit_result": "提交结果",
		"commit_header": "%s: 输入提交描述",
		"cli_error":     "错误: %v",

		// Commit output
		"commit_git_add_fail":    "git add 失败: %s",
		"commit_nothing_new":     "没有新的更改需要提交",
		"commit_git_commit_fail": "git commit 失败: %s",
		"commit_git_push_fail":   "git push 失败: %s",
		"commit_push_done":       "提交并推送完成",
		"commit_pushed_only":     "已推送到远程",

		// Tag
		"tag_prompt":          "创建发布标签？[y/N]",
		"tag_input_ph":        "输入版本号（如 v1.0.0）...",
		"tag_result":          "标签结果",
		"tag_exists":          "标签 '%s' 已存在",
		"tag_success":         "标签 '%s' 已创建并推送",
		"tag_git_tag_fail":    "git tag 失败: %s",
		"tag_git_push_fail":   "git push tag 失败: %s",
		"tag_version_empty":   "版本号不能为空",
		"tag_hint":            "[t] 标签",

		// Duplicate ID
		"msg_duplicate_id": "ID '%s' 已被使用",

		// Action names
		"action_provider_mgmt": "提供商",
		"action_commit":        "提交",
		"action_accounts":      "账户管理",
		"action_settings":      "设置",
		"action_copy":          "复制",
		"action_refresh":       "刷新",
		"action_quit":          "退出",
		"action_new":           "新建",
		"action_delete":        "删除",
		"action_edit_btn":      "编辑",
		"action_select":        "选择",
		"action_back":          "返回",
		"action_confirm_del":   "确认删除",
		"action_nav":           "导航",

		// Settings
		"label_settings":       "设置",
		"settings_language":    "语言",
		"settings_lang_en":     "English",
		"settings_lang_zh":     "简体中文",
		"settings_lang_ja":     "日本語",
		"settings_lang_ko":     "한국어",
		"settings_lang_zh_TW":  "繁體中文",
		"settings_edit_config": "编辑配置",
		"settings_saved":       "设置已保存",
		"settings_save_fail":   "保存设置失败: %v",

		// Header
		"label_no_repo": "非仓库目录",

		// Account management
		"label_accounts":        "账户管理",
		"label_no_accounts":     "未配置任何身份",
		"account_new":           "新建身份",
		"account_edit":          "编辑身份",
		"account_deleted":       "已删除 '%s'",
		"account_saved":         "已保存 '%s'",
		"account_id_ph":         "身份 ID",
		"account_name_ph":       "显示名称",
		"account_gitname_ph":    "Git 用户名",
		"account_gitemail_ph":   "Git 邮箱",
		"account_provider_ph":   "Git 提供商",
		"account_provider_desc": "Git 托管提供商（例: GitHub, Gitee）— [Ctrl+P] 选择，或直接输入",
		"prompt_tab_next":       "[Tab/↑↓] 下一项  [Enter] 保存  [Esc] 取消",
		"account_id_desc":      "用于快速切换身份的唯一标识符（例: work）",
		"account_name_desc":    "显示名称（例: 工作账号）",
		"account_gitname_desc": "git user.name 的值（例: 张三）",
		"account_gitemail_desc": "git user.email 的值（例: zhang@company.com）",
		"account_ssh_identity_ph": "SSH 密钥文件",
		"account_ssh_identity_desc": "SSH 私钥路径（例: ~/.ssh/id_ed25519_work）",
		"ssh_scan_hint":  "[Ctrl+O] 扫描 SSH 密钥",
		"ssh_scan_title": "选择 SSH 密钥",
		"ssh_scan_none":  "~/.ssh/ 中未找到 SSH 密钥。可用 ssh-keygen 命令生成。",
		"provider_pick_hint": "[Ctrl+P] 选择提供商",

		// Messages
		"msg_copied":         "已复制到剪贴板",
		"msg_cannot_copy":    "不在 Git 仓库中，无内容可复制",
		"msg_fields_required": "ID、Git 用户名和 Git 邮箱为必填项",
		"msg_no_profile":    "未选择任何身份",

		// Provider management
		"label_providers":        "提供商管理",
		"label_no_providers":     "未配置任何提供商",
		"provider_new":           "新建提供商",
		"provider_edit":          "编辑提供商",
		"provider_saved":         "提供商已保存",
		"provider_deleted":       "已删除提供商: %s",
		"provider_id_ph":         "提供商 ID",
		"provider_name_ph":       "提供商名称",
		"provider_host_ph":       "主机域名",
		"provider_id_desc":       "唯一标识符（例: github）",
		"provider_name_desc":     "显示名称（例: GitHub, Gitee）",
		"provider_host_desc":     "Git 主机域名（例: github.com, gitee.com）",
		"msg_provider_fields_required": "ID、名称和主机域名为必填项",

		// Provider picker
		"provider_pick_title":  "选择提供商",
		"provider_pick_none":   "未配置任何提供商",

		// Not in repo message
		"msg_press_e": "按 E 编辑配置  |  Q 退出",

		// CLI help
		"help_title":    "GitFace - Git 多身份与企业远程管理器",
		"help_usage":    "用法:",
		"help_cmd_tui":  "  gitf [tui]               启动 TUI（默认）",
		"help_cmd_switch": "  gitf switch <id>         按 ID 切换身份",
		"help_cmd_status": "  gitf status              查看当前 Git 身份与远程仓库",
		"help_cmd_tag":    "  gitf tag <version>       创建并推送发布标签",
		"help_cmd_edit": "  gitf edit                用编辑器打开配置",
		"help_cmd_help": "  gitf help                显示帮助信息",

		// CLI output
		"cli_switch_ok":       "已切换至身份 '%s'",
		"cli_switch_not_found": "未找到身份: %s",
		"cli_switch_detail":   "%s",
		"cli_status_header":   "GitFace 状态",
		"cli_branch":          "分支:  %s",
		"cli_remote":          "远程仓库:  %s",
		"cli_name":            "姓名:  %s",
		"cli_email":           "邮箱:  %s",
		"cli_dirty":           "状态:  %d 个文件未暂存",
		"cli_clean":           "状态:  干净",
		"cli_edit_opening":    "正在打开编辑器...",
		"cli_no_repo":         "不在 Git 仓库中",
		"cli_version":         "GitFace v1.0",
	},

	JA: {
		"app_title":           "GitFace v1.0",
		"label_status":        "【ステータス】",
		"label_branch":        "ブランチ:",
		"label_route":         "リモート:",
		"label_identity":      "身份:",
		"status_dirty":        "%s (%d ファイル未ステージ)",
		"status_clean":        "(クリーン)",
		"status_unable":       "取得不可",
		"status_no_remote":    "リモートなし",
		"status_not_set":      "未設定",
		"status_custom":       "カスタム/未一致",
		"label_actions":       "【操作】",
		"msg_not_repo":        "✗ Gitリポジトリではありません。操作は無効です",
		"msg_checking":        "リポジトリ状態を確認中...",
		"msg_switch_fail":     "切替失敗: %v",
		"msg_switch_ok":       "✓ 身分の切替に成功しました",
		"msg_edit_fail":       "設定の編集に失敗: %v",
		"msg_reload_ok":       "✓ 設定をリロードしました",
		"msg_cannot_commit":   "Gitリポジトリではありません。コミットできません",
		"msg_cannot_switch":   "Gitリポジトリではありません。身分を切替できません",
		"msg_desc_empty":      "コミット説明は空にできません",

		// SwitchProfile logs
		"sp_name_fail":      "user.name の設定に失敗: %s",
		"sp_email_fail":     "user.email の設定に失敗: %s",
		"sp_remote_same":    "リモートURLは変更ありません",
		"sp_remote_fail":    "リモートURLの設定に失敗: %s",
		"sp_no_url":         "プロバイダ %s に関連URLがありません",
		"sp_prov_not_found": "プロバイダ ID '%s' が見つかりません、スキップ",
		"sp_no_provider":    "プロバイダが設定されていません、スキップ",
		"sp_ssh_fail":       "SSH コマンドの設定に失敗: %s",
		"prompt_select_type":  "コミットタイプを選択:",
		"prompt_enter_desc_ph": "コミット説明を入力...",
		"prompt_confirm_cancel": "[Enter] 確認  [Esc] キャンセル",
		"prompt_cancel":       "[Esc] キャンセル",
		"prompt_any_key":      "任意キーで戻る",
		"commit_result":       "コミット結果",
		"commit_header":       "%s: 説明を入力",
		"cli_error":           "エラー: %v",

		// Commit output
		"commit_git_add_fail":    "git add に失敗: %s",
		"commit_nothing_new":     "新しい変更がありません",
		"commit_git_commit_fail": "git commit に失敗: %s",
		"commit_git_push_fail":   "git push に失敗: %s",
		"commit_push_done":       "コミットしてプッシュしました",
		"commit_pushed_only":     "リモートにプッシュしました",

		// Tag
		"tag_prompt":          "リリースタグを作成しますか？[y/N]",
		"tag_input_ph":        "バージョンを入力（例: v1.0.0）...",
		"tag_result":          "タグ結果",
		"tag_exists":          "タグ '%s' は既に存在します",
		"tag_success":         "タグ '%s' を作成してプッシュしました",
		"tag_git_tag_fail":    "git tag に失敗: %s",
		"tag_git_push_fail":   "git push tag に失敗: %s",
		"tag_version_empty":   "バージョンを入力してください",
		"tag_hint":            "[t] タグ",

		// Duplicate ID
		"msg_duplicate_id": "ID '%s' はすでに使用されています",
		"action_provider_mgmt": "管理",
		"action_commit":       "コミット",
		"action_accounts":     "アカウント",
		"action_settings":     "設定",
		"action_copy":         "コピー",
		"action_refresh":      "更新",
		"action_quit":         "終了",
		"action_new":          "新規",
		"action_delete":       "削除",
		"action_edit_btn":     "編集",
		"action_select":       "選択",
		"action_back":         "戻る",
		"action_confirm_del":  "削除確認",
		"action_nav":          "ナビゲート",
		"label_settings":      "設定",
		"settings_language":   "言語",
		"settings_lang_en":    "English",
		"settings_lang_zh":    "简体中文",
		"settings_lang_ja":    "日本語",
		"settings_lang_ko":    "한국어",
		"settings_lang_zh_TW": "繁體中文",
		"settings_edit_config": "設定を編集",
		"settings_saved":      "設定を保存しました",
		"settings_save_fail":  "設定の保存に失敗: %v",
		"label_no_repo":       "リポジトリではない",
		"label_accounts":      "アカウント管理",
		"label_no_accounts":   "プロファイルが未設定です",
		"account_new":         "新規アカウント",
		"account_edit":        "アカウントを編集",
		"account_deleted":     "プロファイル '%s' を削除しました",
		"account_saved":       "プロファイル '%s' を保存しました",
		"account_id_ph":       "プロファイル ID",
		"account_name_ph":     "表示名",
		"account_gitname_ph":  "Git ユーザー名",
		"account_gitemail_ph": "Git メール",
		"account_provider_ph": "プロバイダー",
		"account_provider_desc": "Gitホスティングプロバイダー（例: GitHub, Gitee）— [Ctrl+P] で選択、または入力",
		"prompt_tab_next":     "[Tab/↑↓] 次へ  [Enter] 保存  [Esc] キャンセル",
		"account_id_desc":     "クイック切替用のユニークID（例: work）",
		"account_name_desc":   "表示名（例: 仕事アカウント）",
		"account_gitname_desc": "git user.name の値（例: 山田太郎）",
		"account_gitemail_desc": "git user.email の値（例: yamada@company.com）",
		"account_ssh_identity_ph": "SSH 身分ファイル",
		"account_ssh_identity_desc": "SSH秘密鍵のパス（例: ~/.ssh/id_ed25519_work）",
		"ssh_scan_hint":       "[Ctrl+O] SSH鍵をスキャン",
		"ssh_scan_title":      "SSH鍵を選択",
		"ssh_scan_none":       "~/.ssh/ にSSH鍵が見つかりません。ssh-keygen で生成してください。",
		"provider_pick_hint":  "[Ctrl+P] プロバイダーを選択",
		"msg_copied":          "クリップボードにコピーしました",
		"msg_cannot_copy":     "Gitリポジトリではありません。コピーするものはありません",
		"msg_fields_required": "ID、Gitユーザー名、Gitメールは必須です",
		"msg_no_profile":      "プロファイルが選択されていません",
		"label_providers":     "プロバイダー管理",
		"label_no_providers":  "プロバイダーが未設定です",
		"provider_new":        "新規プロバイダー",
		"provider_edit":       "プロバイダーを編集",
		"provider_saved":      "プロバイダーを保存しました",
		"provider_deleted":    "プロバイダーを削除しました: %s",
		"provider_id_ph":      "プロバイダー ID",
		"provider_name_ph":    "プロバイダー名",
		"provider_host_ph":    "ホストドメイン",
		"provider_id_desc":    "ユニークID（例: github）",
		"provider_name_desc":  "表示名（例: GitHub, Gitee）",
		"provider_host_desc":  "Gitホストドメイン（例: github.com, gitee.com）",
		"msg_provider_fields_required": "ID、名前、ホストは必須です",
		"provider_pick_title": "プロバイダーを選択",
		"provider_pick_none":  "プロバイダーが未設定です",
		"msg_press_e":         "E で設定を編集  |  Q で終了",
		"help_title":          "GitFace - Git 多身分・リモート管理",
		"help_usage":          "使い方:",
		"help_cmd_tui":        "  gitf [tui]               TUIを起動（デフォルト）",
		"help_cmd_switch":     "  gitf switch <id>         IDで身分を切替",
		"help_cmd_status":     "  gitf status              現在のGit身分とリモートを表示",
		"help_cmd_tag":        "  gitf tag <version>       リリースタグを作成してプッシュ",
		"help_cmd_edit":       "  gitf edit                エディタで設定を開く",
		"help_cmd_help":       "  gitf help                このヘルプを表示",
		"cli_switch_ok":       "プロファイル '%s' に切替ました",
		"cli_switch_not_found": "プロファイルが見つかりません: %s",
		"cli_switch_detail":   "%s",
		"cli_status_header":   "GitFace ステータス",
		"cli_branch":          "ブランチ:  %s",
		"cli_remote":          "リモート:  %s",
		"cli_name":            "名前:  %s",
		"cli_email":           "メール:  %s",
		"cli_dirty":           "状態:  %d ファイル未ステージ",
		"cli_clean":           "状態:  クリーン",
		"cli_edit_opening":    "エディタを開いています...",
		"cli_no_repo":         "Gitリポジトリではありません",
		"cli_version":         "GitFace v1.0",
	},

	KO: {
		"app_title":           "GitFace v1.0",
		"label_status":        "【상태】",
		"label_branch":        "브랜치:",
		"label_route":         "원격:",
		"label_identity":      "신분:",
		"status_dirty":        "%s (%d 파일 언스테이지)",
		"status_clean":        "(클린)",
		"status_unable":       "획득 불가",
		"status_no_remote":    "원격 없음",
		"status_not_set":      "미설정",
		"status_custom":       "사용자 정의/미일치",
		"label_actions":       "【작업】",
		"msg_not_repo":        "✗ Git 저장소가 아닙니다. 작업이 비활성화됩니다",
		"msg_checking":        "저장소 상태 확인 중...",
		"msg_switch_fail":     "전환 실패: %v",
		"msg_switch_ok":       "✓ 신분 전환에 성공했습니다",
		"msg_edit_fail":       "설정 편집 실패: %v",
		"msg_reload_ok":       "✓ 설정을 다시 로드했습니다",
		"msg_cannot_commit":   "Git 저장소가 아닙니다. 커밋할 수 없습니다",
		"msg_cannot_switch":   "Git 저장소가 아닙니다. 신분을 전환할 수 없습니다",
		"msg_desc_empty":      "커밋 설명은 비워둘 수 없습니다",

		// SwitchProfile logs
		"sp_name_fail":      "user.name 설정 실패: %s",
		"sp_email_fail":     "user.email 설정 실패: %s",
		"sp_remote_same":    "원격 URL 변경 없음",
		"sp_remote_fail":    "원격 URL 설정 실패: %s",
		"sp_no_url":         "제공자 %s에 연결된 URL이 없습니다",
		"sp_prov_not_found": "제공자 ID '%s'를 찾을 수 없습니다. 건너뜀",
		"sp_no_provider":    "제공자가 설정되지 않았습니다. 건너뜀",
		"sp_ssh_fail":       "SSH 명령 설정 실패: %s",
		"prompt_select_type":  "커밋 유형 선택:",
		"prompt_enter_desc_ph": "커밋 설명을 입력하세요...",
		"prompt_confirm_cancel": "[Enter] 확인  [Esc] 취소",
		"prompt_cancel":       "[Esc] 취소",
		"prompt_any_key":      "아무 키나 눌러 돌아가기",
		"commit_result":       "커밋 결과",
		"commit_header":       "%s: 설명을 입력하세요",
		"cli_error":           "오류: %v",

		// Commit output
		"commit_git_add_fail":    "git add 실패: %s",
		"commit_nothing_new":     "새로운 변경 사항이 없습니다",
		"commit_git_commit_fail": "git commit 실패: %s",
		"commit_git_push_fail":   "git push 실패: %s",
		"commit_push_done":       "커밋 및 푸시 완료",
		"commit_pushed_only":     "원격으로 푸시했습니다",

		// Tag
		"tag_prompt":          "릴리스 태그를 생성하시겠습니까？[y/N]",
		"tag_input_ph":        "버전을 입력하세요（예: v1.0.0）...",
		"tag_result":          "태그 결과",
		"tag_exists":          "태그 '%s'가 이미 존재합니다",
		"tag_success":         "태그 '%s'를 생성하고 푸시했습니다",
		"tag_git_tag_fail":    "git tag 실패: %s",
		"tag_git_push_fail":   "git push tag 실패: %s",
		"tag_version_empty":   "버전을 입력하세요",
		"tag_hint":            "[t] 태그",

		// Duplicate ID
		"msg_duplicate_id": "ID '%s'는 이미 사용 중입니다",
		"action_provider_mgmt": "관리",
		"action_commit":       "커밋",
		"action_accounts":     "계정",
		"action_settings":     "설정",
		"action_copy":         "복사",
		"action_refresh":      "새로고침",
		"action_quit":         "종료",
		"action_new":          "새로 만들기",
		"action_delete":       "삭제",
		"action_edit_btn":     "편집",
		"action_select":       "선택",
		"action_back":         "뒤로",
		"action_confirm_del":  "삭제 확인",
		"action_nav":          "탐색",
		"label_settings":      "설정",
		"settings_language":   "언어",
		"settings_lang_en":    "English",
		"settings_lang_zh":    "简体中文",
		"settings_lang_ja":    "日本語",
		"settings_lang_ko":    "한국어",
		"settings_lang_zh_TW": "繁體中文",
		"settings_edit_config": "설정 편집",
		"settings_saved":      "설정을 저장했습니다",
		"settings_save_fail":  "설정 저장 실패: %v",
		"label_no_repo":       "저장소 아님",
		"label_accounts":      "계정 관리",
		"label_no_accounts":   "프로필이 설정되지 않았습니다",
		"account_new":         "새 계정",
		"account_edit":        "계정 편집",
		"account_deleted":     "프로필 '%s'을(를) 삭제했습니다",
		"account_saved":       "프로필 '%s'을(를) 저장했습니다",
		"account_id_ph":       "프로필 ID",
		"account_name_ph":     "표시 이름",
		"account_gitname_ph":  "Git 사용자 이름",
		"account_gitemail_ph": "Git 이메일",
		"account_provider_ph": "프로바이더",
		"account_provider_desc": "Git 호스팅 프로바이더 (예: GitHub, Gitee) — [Ctrl+P]로 선택하거나 입력",
		"prompt_tab_next":     "[Tab/↑↓] 다음  [Enter] 저장  [Esc] 취소",
		"account_id_desc":     "빠른 전환용 고유 ID (예: work)",
		"account_name_desc":   "표시 이름 (예: 업무 계정)",
		"account_gitname_desc": "git user.name 값 (예: 김철수)",
		"account_gitemail_desc": "git user.email 값 (예: kim@company.com)",
		"account_ssh_identity_ph": "SSH 신분 파일",
		"account_ssh_identity_desc": "SSH 개인키 경로 (예: ~/.ssh/id_ed25519_work)",
		"ssh_scan_hint":       "[Ctrl+O] SSH 키 스캔",
		"ssh_scan_title":      "SSH 키 선택",
		"ssh_scan_none":       "~/.ssh/에 SSH 키가 없습니다. ssh-keygen으로 생성하세요.",
		"provider_pick_hint":  "[Ctrl+P] 프로바이더 선택",
		"msg_copied":          "클립보드에 복사했습니다",
		"msg_cannot_copy":     "Git 저장소가 아닙니다. 복사할 것이 없습니다",
		"msg_fields_required": "ID, Git 사용자 이름, Git 이메일은 필수입니다",
		"msg_no_profile":      "프로필이 선택되지 않았습니다",
		"label_providers":     "프로바이더 관리",
		"label_no_providers":  "프로바이더가 설정되지 않았습니다",
		"provider_new":        "새 프로바이더",
		"provider_edit":       "프로바이더 편집",
		"provider_saved":      "프로바이더를 저장했습니다",
		"provider_deleted":    "프로바이더를 삭제했습니다: %s",
		"provider_id_ph":      "프로바이더 ID",
		"provider_name_ph":    "프로바이더 이름",
		"provider_host_ph":    "호스트 도메인",
		"provider_id_desc":    "고유 ID (예: github)",
		"provider_name_desc":  "표시 이름 (예: GitHub, Gitee)",
		"provider_host_desc":  "Git 호스트 도메인 (예: github.com, gitee.com)",
		"msg_provider_fields_required": "ID, 이름, 호스트는 필수입니다",
		"provider_pick_title": "프로바이더 선택",
		"provider_pick_none":  "프로바이더가 설정되지 않았습니다",
		"msg_press_e":         "E로 설정 편집  |  Q로 종료",
		"help_title":          "GitFace - Git 다중 신분 & 원격 관리자",
		"help_usage":          "사용법:",
		"help_cmd_tui":        "  gitf [tui]               TUI 실행 (기본값)",
		"help_cmd_switch":     "  gitf switch <id>         ID로 신분 전환",
		"help_cmd_status":     "  gitf status              현재 Git 신분과 원격 표시",
		"help_cmd_tag":        "  gitf tag <version>       릴리스 태그를 생성하고 푸시",
		"help_cmd_edit":       "  gitf edit                에디터로 설정 열기",
		"help_cmd_help":       "  gitf help                도움말 표시",
		"cli_switch_ok":       "프로필 '%s'로 전환했습니다",
		"cli_switch_not_found": "프로필을 찾을 수 없습니다: %s",
		"cli_switch_detail":   "%s",
		"cli_status_header":   "GitFace 상태",
		"cli_branch":          "브랜치:  %s",
		"cli_remote":          "원격:  %s",
		"cli_name":            "이름:  %s",
		"cli_email":           "이메일:  %s",
		"cli_dirty":           "상태:  %d 파일 언스테이지",
		"cli_clean":           "상태:  클린",
		"cli_edit_opening":    "에디터를 열고 있습니다...",
		"cli_no_repo":         "Git 저장소가 아닙니다",
		"cli_version":         "GitFace v1.0",
	},

	ZH_TW: {
		"app_title":           "GitFace v1.0",
		"label_status":        "【目前狀態】",
		"label_branch":        "分支:",
		"label_route":         "遠端:",
		"label_identity":      "身分:",
		"status_dirty":        "%s (%d 檔案未暫存)",
		"status_clean":        "(clean)",
		"status_unable":       "無法取得",
		"status_no_remote":    "無遠端",
		"status_not_set":      "未設定",
		"status_custom":       "自訂/未相符",
		"label_actions":       "【快捷操作】",
		"msg_not_repo":        "✗ 目前目錄不是 Git 儲存庫，操作已停用",
		"msg_checking":        "檢查儲存庫狀態中...",
		"msg_switch_fail":     "切換失敗: %v",
		"msg_switch_ok":       "✓ 身分切換成功",
		"msg_edit_fail":       "編輯設定失敗: %v",
		"msg_reload_ok":       "✓ 設定已重新載入",
		"msg_cannot_commit":   "不在 Git 儲存庫中，無法提交",
		"msg_cannot_switch":   "不在 Git 儲存庫中，無法切換身分",
		"msg_desc_empty":      "提交描述不能為空",

		// SwitchProfile logs
		"sp_name_fail":      "設定 user.name 失敗: %s",
		"sp_email_fail":     "設定 user.email 失敗: %s",
		"sp_remote_same":    "遠端位址未變化",
		"sp_remote_fail":    "設定遠端倉庫位址失敗: %s",
		"sp_no_url":         "提供商 %s 已關聯，無遠端 URL",
		"sp_prov_not_found": "未找到提供商 ID '%s'，跳過遠端位址設定",
		"sp_no_provider":    "未關聯提供商，跳過遠端位址設定",
		"sp_ssh_fail":       "設定 SSH 命令失敗: %s",
		"prompt_select_type":  "選擇提交類型:",
		"prompt_enter_desc_ph": "輸入提交描述...",
		"prompt_confirm_cancel": "[Enter] 確認  [Esc] 取消",
		"prompt_cancel":       "[Esc] 取消",
		"prompt_any_key":      "按任意鍵返回",
		"commit_result":       "提交結果",
		"commit_header":       "%s: 輸入提交描述",
		"cli_error":           "錯誤: %v",

		// Commit output
		"commit_git_add_fail":    "git add 失敗: %s",
		"commit_nothing_new":     "沒有新的更改需要提交",
		"commit_git_commit_fail": "git commit 失敗: %s",
		"commit_git_push_fail":   "git push 失敗: %s",
		"commit_push_done":       "提交並推送完成",
		"commit_pushed_only":     "已推送到遠端",

		// Tag
		"tag_prompt":          "建立發布標籤？[y/N]",
		"tag_input_ph":        "輸入版本號（如 v1.0.0）...",
		"tag_result":          "標籤結果",
		"tag_exists":          "標籤 '%s' 已存在",
		"tag_success":         "標籤 '%s' 已建立並推送",
		"tag_git_tag_fail":    "git tag 失敗: %s",
		"tag_git_push_fail":   "git push tag 失敗: %s",
		"tag_version_empty":   "版本號不能為空",
		"tag_hint":            "[t] 標籤",

		// Duplicate ID
		"msg_duplicate_id": "ID '%s' 已被使用",
		"action_provider_mgmt": "管理",
		"action_commit":       "提交",
		"action_accounts":     "帳戶管理",
		"action_settings":     "設定",
		"action_copy":         "複製",
		"action_refresh":      "重新整理",
		"action_quit":         "離開",
		"action_new":          "新增",
		"action_delete":       "刪除",
		"action_edit_btn":     "編輯",
		"action_select":       "選擇",
		"action_back":         "返回",
		"action_confirm_del":  "確認刪除",
		"action_nav":          "導覽",
		"label_settings":      "設定",
		"settings_language":   "語言",
		"settings_lang_en":    "English",
		"settings_lang_zh":    "简体中文",
		"settings_lang_zh_TW": "繁體中文",
		"settings_lang_ja":    "日本語",
		"settings_lang_ko":    "한국어",
		"settings_edit_config": "編輯設定",
		"settings_saved":      "設定已儲存",
		"settings_save_fail":  "儲存設定失敗: %v",
		"label_no_repo":       "非儲存庫目錄",
		"label_accounts":      "帳戶管理",
		"label_no_accounts":   "未設定任何身分",
		"account_new":         "新增身分",
		"account_edit":        "編輯身分",
		"account_deleted":     "已刪除 '%s'",
		"account_saved":       "已儲存 '%s'",
		"account_id_ph":       "身分 ID",
		"account_name_ph":     "顯示名稱",
		"account_gitname_ph":  "Git 使用者名稱",
		"account_gitemail_ph": "Git 信箱",
		"account_provider_ph": "Git 提供者",
		"account_provider_desc": "Git 托管提供者（例: GitHub, Gitee）— [Ctrl+P] 選擇，或直接輸入",
		"prompt_tab_next":     "[Tab/↑↓] 下一項  [Enter] 儲存  [Esc] 取消",
		"account_id_desc":     "用於快速切換身分的唯一識別碼（例: work）",
		"account_name_desc":   "顯示名稱（例: 工作帳戶）",
		"account_gitname_desc": "git user.name 的值（例: 王小明）",
		"account_gitemail_desc": "git user.email 的值（例: wang@company.com）",
		"account_ssh_identity_ph": "SSH 身分檔案",
		"account_ssh_identity_desc": "SSH 私鑰路徑（例: ~/.ssh/id_ed25519_work）",
		"ssh_scan_hint":       "[Ctrl+O] 掃描 SSH 金鑰",
		"ssh_scan_title":      "選擇 SSH 金鑰",
		"ssh_scan_none":       "~/.ssh/ 中未找到 SSH 金鑰。可用 ssh-keygen 命令產生。",
		"provider_pick_hint":  "[Ctrl+P] 選擇提供者",
		"msg_copied":          "已複製到剪貼簿",
		"msg_cannot_copy":     "不在 Git 儲存庫中，無內容可複製",
		"msg_fields_required": "ID、Git 使用者名稱和 Git 信箱為必填項",
		"msg_no_profile":      "未選擇任何身分",
		"label_providers":     "提供者管理",
		"label_no_providers":  "未設定任何提供者",
		"provider_new":        "新增提供者",
		"provider_edit":       "編輯提供者",
		"provider_saved":      "提供者已儲存",
		"provider_deleted":    "已刪除提供者: %s",
		"provider_id_ph":      "提供者 ID",
		"provider_name_ph":    "提供者名稱",
		"provider_host_ph":    "主機網域",
		"provider_id_desc":    "唯一識別碼（例: github）",
		"provider_name_desc":  "顯示名稱（例: GitHub, Gitee）",
		"provider_host_desc":  "Git 主機網域（例: github.com, gitee.com）",
		"msg_provider_fields_required": "ID、名稱和主機網域為必填項",
		"provider_pick_title": "選擇提供者",
		"provider_pick_none":  "未設定任何提供者",
		"msg_press_e":         "按 E 編輯設定  |  Q 離開",
		"help_title":          "GitFace - Git 多身分與遠端管理器",
		"help_usage":          "用法:",
		"help_cmd_tui":        "  gitf [tui]               啟動 TUI（預設）",
		"help_cmd_switch":     "  gitf switch <id>         按 ID 切換身分",
		"help_cmd_status":     "  gitf status              查看目前 Git 身分與遠端",
		"help_cmd_tag":        "  gitf tag <version>       建立並推送發布標籤",
		"help_cmd_edit":       "  gitf edit                用編輯器開啟設定",
		"help_cmd_help":       "  gitf help                顯示說明資訊",
		"cli_switch_ok":       "已切換至身分 '%s'",
		"cli_switch_not_found": "未找到身分: %s",
		"cli_switch_detail":   "%s",
		"cli_status_header":   "GitFace 狀態",
		"cli_branch":          "分支:  %s",
		"cli_remote":          "遠端:  %s",
		"cli_name":            "姓名:  %s",
		"cli_email":           "信箱:  %s",
		"cli_dirty":           "狀態:  %d 個檔案未暫存",
		"cli_clean":           "狀態:  乾淨",
		"cli_edit_opening":    "正在開啟編輯器...",
		"cli_no_repo":         "不在 Git 儲存庫中",
		"cli_version":         "GitFace v1.0",
	},
}
