package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateNormal state = iota
	stateCommitSelect
	stateCommitInput
	stateCommitOutput
	stateAccountManage
	stateAccountForm
	stateProviderManage
	stateProviderForm
	stateProviderPicker
	stateSettings
	stateSSHKeyPicker
	stateTagPrompt
	stateTagInput
	stateRepoConfig
	stateRepoConfigEdit
	statePRInput
	statePROutput
	stateBranchSwitch
)

type refreshMsg struct {
	isRepo     bool
	branch     string
	remoteURL  string
	localName  string
	localEmail string
	isDirty    bool
	dirtyCount int
}

type switchDoneMsg struct {
	log string
	err error
}

type commitResultMsg struct {
	output string
	err    error
}

type editDoneMsg struct {
	err error
}

type accountSavedMsg struct {
	err error
	name string
}

type accountDeletedMsg struct {
	err  error
	name string
}

type providerSavedMsg struct {
	err error
}

type providerDeletedMsg struct {
	err      error
	provider string
}

type copiedMsg struct {
	err error
}

type tagResultMsg struct {
	output string
	err    error
}

type prResultMsg struct {
	output string
	err    error
}

type settingsSavedMsg struct {
	err error
}

type model struct {
	cfg     *Config
	cfgPath string
	tr      *Translator

	isRepo        bool
	repoChecked   bool
	branch        string
	remoteURL     string
	localName     string
	localEmail    string
	isDirty       bool
	dirtyCount    int

	state state
	errMsg string
	infoMsg string

	commitType  string
	commitInput textarea.Model
	commitOutput string
	commitErr   error
	prOutput    string
	prErr       error

	tagInput     textinput.Model
	prTitleInput textinput.Model

	accountEditIdx int
	accountInputs  []textinput.Model
	accountFocused int

	providerCursor   int
	providerEditIdx  int
	providerInputs   []textinput.Model
	providerFocused  int

	cursor            int
	commitCursor      int
	accountCursor     int
	delMode           bool
	providerDelMode   bool

	settingsCursor int

	sshKeys           []string
	sshKeyCursor      int
	providerPickList  []Provider
	providerPickCursor int

	repoConfig       *RepoConfig
	repoConfigCursor int
	repoConfigKeys   []string
	repoConfigProviders []Provider
	repoConfigInput  textinput.Model

	width  int
	height int

	branches     []string
	branchCursor int
}

var commitTypes = []struct {
	key   string
	label string
	desc  string
}{
	{"feat", "feat", "new feature"},
	{"fix", "fix", "bug fix"},
	{"docs", "docs", "documentation"},
	{"refactor", "refactor", "code refactor"},
}

func NewModel(cfg *Config, cfgPath string, tr *Translator) model {
	ti := textarea.New()
	ti.Placeholder = tr.Tr("prompt_enter_desc_ph")
	ti.ShowLineNumbers = false
	ti.SetHeight(8)
	ti.SetWidth(80)
	ti.CharLimit = 1000
	ti.Prompt = ""
	ti.Focus()

	tgi := textinput.New()
	tgi.Placeholder = "v1.0.0"
	tgi.Focus()
	tgi.CharLimit = 50

	pti := textinput.New()
	pti.Placeholder = tr.Tr("pr_input_ph")
	pti.Focus()
	pti.CharLimit = 200

	return model{
		cfg:         cfg,
		cfgPath:     cfgPath,
		tr:          tr,
		state:       stateNormal,
		commitInput:  ti,
		tagInput:     tgi,
		prTitleInput: pti,
		cursor:      0,
		width:       80,
		height:      24,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.commitInput.Focus(), m.refreshCmd())
}

func (m model) refreshCmd() tea.Cmd {
	return func() tea.Msg {
		isRepo := IsInsideWorkTree() == nil
		msg := refreshMsg{isRepo: isRepo}
		if isRepo {
			msg.branch = GetCurrentBranch()
			msg.remoteURL = GetRemoteURL()
			msg.localName = GetLocalUserName()
			msg.localEmail = GetLocalUserEmail()
			msg.isDirty, msg.dirtyCount = GetWorkTreeStatus()
		}
		return msg
	}
}

func (m model) saveConfigCmd(profileName string) tea.Cmd {
	return func() tea.Msg {
		err := SaveConfig(m.cfg, m.cfgPath)
		return accountSavedMsg{err: err, name: profileName}
	}
}

func (m model) copyCmd(text string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(text)
		return copiedMsg{err: err}
	}
}

func (m model) createTagCmd(version string) tea.Cmd {
	return func() tea.Msg {
		output, err := CreateTag(version, m.tr)
		return tagResultMsg{output: output, err: err}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		w := msg.Width - 12
		if w < 40 {
			w = 40
		}
		m.commitInput.SetWidth(w)
		h := msg.Height - 14
		if h < 3 {
			h = 3
		}
		if h > 30 {
			h = 30
		}
		m.commitInput.SetHeight(h)
		return m, nil

	case tea.MouseMsg:
		return m.handleMouseEvent(msg), nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case refreshMsg:
		m.repoChecked = true
		m.isRepo = msg.isRepo
		if msg.isRepo {
			m.branch = msg.branch
			m.remoteURL = msg.remoteURL
			m.localName = msg.localName
			m.localEmail = msg.localEmail
			m.isDirty = msg.isDirty
			m.dirtyCount = msg.dirtyCount
		}
		return m, nil

	case switchDoneMsg:
		m.errMsg = ""
		m.infoMsg = ""
		if msg.err != nil {
			m.errMsg = m.tr.Tr("msg_switch_fail", msg.err)
		} else {
			parts := []string{m.tr.Tr("msg_switch_ok")}
			for _, l := range strings.Split(msg.log, "\n") {
				if l = strings.TrimSpace(l); l != "" {
					parts = append(parts, l)
				}
			}
			m.infoMsg = strings.Join(parts, "\n")
		}
		return m, m.refreshCmd()

	case commitResultMsg:
		m.commitOutput = msg.output
		m.commitErr = msg.err
		m.state = stateCommitOutput
		return m, nil

	case tagResultMsg:
		m.commitOutput = msg.output
		m.commitErr = msg.err
		m.state = stateCommitOutput
		return m, nil

	case prResultMsg:
		m.commitOutput = msg.output
		m.commitErr = msg.err
		m.state = stateCommitOutput
		return m, nil

	case editDoneMsg:
		if msg.err != nil {
			m.errMsg = m.tr.Tr("msg_edit_fail", msg.err)
			return m, nil
		}
		cfg, err := LoadConfig(m.cfgPath)
		if err != nil {
			m.errMsg = fmt.Sprintf("reload config failed: %v", err)
		} else {
			detected := DetectLang()
			lang := detected
			if cfg.Lang != "" {
				lang = Lang(cfg.Lang)
			}
			m.cfg = cfg
			m.tr = NewTranslator(lang)
			m.commitInput.Placeholder = m.tr.Tr("prompt_enter_desc_ph")
			m.infoMsg = m.tr.Tr("msg_reload_ok")
		}
		return m, m.refreshCmd()

	case accountSavedMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("save failed: %v", msg.err)
		} else {
			m.infoMsg = m.tr.Tr("account_saved", msg.name)
		}
		m.state = stateAccountManage
		return m, nil

	case accountDeletedMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("delete failed: %v", msg.err)
		} else {
			m.infoMsg = m.tr.Tr("account_deleted", msg.name)
		}
		m.state = stateAccountManage
		if m.accountCursor >= len(m.cfg.Profiles) {
			m.accountCursor = len(m.cfg.Profiles) - 1
		}
		if m.accountCursor < 0 {
			m.accountCursor = 0
		}
		if m.cursor >= len(m.cfg.Profiles) {
			m.cursor = len(m.cfg.Profiles) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		return m, nil

	case providerSavedMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("save failed: %v", msg.err)
		} else {
			m.infoMsg = m.tr.Tr("provider_saved")
		}
		m.state = stateProviderManage
		return m, nil

	case providerDeletedMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("delete failed: %v", msg.err)
		} else {
			m.infoMsg = m.tr.Tr("provider_deleted", msg.provider)
		}
		m.state = stateProviderManage
		if m.providerCursor >= len(m.cfg.Providers) {
			m.providerCursor = len(m.cfg.Providers) - 1
		}
		return m, nil

	case copiedMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("copy failed: %v", msg.err)
		} else {
			m.infoMsg = m.tr.Tr("msg_copied")
		}
		return m, nil

	case settingsSavedMsg:
		if msg.err != nil {
			m.errMsg = m.tr.Tr("settings_save_fail", msg.err)
		} else {
			m.infoMsg = m.tr.Tr("settings_saved")
		}
		return m, nil
	}

	return m, nil
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	ks := msg.String()
	if ks == "q" || ks == "Q" {
		if m.state == stateNormal || m.state == stateAccountManage || m.state == stateProviderManage || m.state == stateSettings {
			return m, tea.Quit
		}
	}

	m.infoMsg = ""
	m.errMsg = ""

	switch m.state {
	case stateNormal:
		return m.handleNormalKey(msg)
	case stateCommitSelect:
		return m.handleCommitSelectKey(msg)
	case stateCommitInput:
		return m.handleCommitInputKey(msg)
	case statePRInput:
		return m.handlePRInputKey(msg)
	case stateCommitOutput:
		switch ks {
		case "y", "Y":
			var text string
			if m.commitErr != nil {
				text = m.commitOutput + "\n" + m.tr.Tr("cli_error", m.commitErr)
			} else {
				text = m.commitOutput
			}
			if text == "" {
				m.errMsg = m.tr.Tr("msg_cannot_copy")
			} else {
				m.infoMsg = m.tr.Tr("msg_copied")
				return m, m.copyCmd(text)
			}
			return m, nil
		case "r", "R":
			m.state = statePRInput
			m.prTitleInput.SetValue(GetLastCommitMessage())
			return m, m.prTitleInput.Focus()
		default:
			if m.commitOutput == "..." {
				return m, nil
			}
			m.state = stateNormal
			m.commitOutput = ""
			m.commitErr = nil
			m.infoMsg = ""
			m.errMsg = ""
			return m, nil
		}
	case stateTagInput:
		return m.handleTagInputKey(msg)
	case stateAccountManage:
		return m.handleAccountManageKey(msg)
	case stateAccountForm:
		return m.handleAccountFormKey(msg)
	case stateProviderManage:
		return m.handleProviderManageKey(msg)
	case stateProviderForm:
		return m.handleProviderFormKey(msg)
	case stateProviderPicker:
		return m.handleProviderPickerKey(msg)
	case stateSettings:
		return m.handleSettingsKey(msg)
	case stateSSHKeyPicker:
		return m.handleSSHKeyPickerKey(msg)
	case stateRepoConfig:
		return m.handleRepoConfigKey(msg)
	case stateRepoConfigEdit:
		return m.handleRepoConfigEditKey(msg)
	case statePROutput:
		m.state = stateNormal
		m.commitOutput = ""
		m.commitErr = nil
		m.prOutput = ""
		m.prErr = nil
		return m, nil
	case stateBranchSwitch:
		return m.handleBranchSwitchKey(msg)
	}
	return m, nil
}

func (m model) totalItems() int {
	return len(m.cfg.Profiles) + 8
}

func (m model) dispatchAction() (tea.Model, tea.Cmd) {
	np := len(m.cfg.Profiles)
	switch {
	case m.cursor < np:
		return m.switchToCursor()
	case m.cursor == np:
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		m.state = stateCommitSelect
		m.commitOutput = ""
		m.commitErr = nil
		m.commitCursor = 0
		return m, nil
	case m.cursor == np+1:
		m.errMsg = ""
		m.infoMsg = ""
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		m.state = stateTagInput
		m.tagInput.SetValue("")
		return m, textinput.Blink
	case m.cursor == np+2:
		m.errMsg = ""
		m.infoMsg = ""
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		m.state = statePRInput
		m.prTitleInput.SetValue(GetLastCommitMessage())
		return m, m.prTitleInput.Focus()
	case m.cursor == np+3:
		m.errMsg = ""
		m.infoMsg = ""
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		rc := LoadRepoConfig()
		m.repoConfig = rc
		m.repoConfigProviders = m.cfg.Providers
		m.repoConfigKeys = make([]string, 0, len(m.cfg.Providers))
		for _, p := range m.cfg.Providers {
			m.repoConfigKeys = append(m.repoConfigKeys, p.ID)
		}
		m.repoConfigCursor = 0
		m.state = stateRepoConfig
		return m, nil
	case m.cursor == np+4:
		m.errMsg = ""
		m.infoMsg = ""
		m.state = stateAccountManage
		m.accountCursor = 0
		m.delMode = false
		return m, nil
	case m.cursor == np+5:
		m.errMsg = ""
		m.infoMsg = ""
		m.state = stateProviderManage
		m.providerCursor = 0
		m.providerDelMode = false
		return m, nil
	case m.cursor == np+6:
		m.errMsg = ""
		m.infoMsg = ""
		m.state = stateSettings
		m.settingsCursor = 0
		return m, nil
	case m.cursor == np+7:
		m.errMsg = ""
		m.infoMsg = ""
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		bd := ListBranches()
		if len(bd) == 0 {
			m.errMsg = m.tr.Tr("msg_no_branches")
			return m, nil
		}
		m.branches = bd
		m.branchCursor = 0
		for i, b := range bd {
			if b == m.branch {
				m.branchCursor = i
				break
			}
		}
		m.state = stateBranchSwitch
		return m, nil
	}
	return m, nil
}

func (m model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	total := m.totalItems()
	switch msg.String() {
	case "up", "k":
		m.cursor--
		if m.cursor < 0 {
			m.cursor = total - 1
		}
		return m, nil

	case "down", "j":
		m.cursor++
		if m.cursor >= total {
			m.cursor = 0
		}
		return m, nil

	case "enter":
		return m.dispatchAction()

	case "r", "R":
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_no_repo")
			return m, nil
		}
		m.errMsg = ""
		m.infoMsg = ""
		m.state = statePRInput
		m.prTitleInput.SetValue(GetLastCommitMessage())
		return m, m.prTitleInput.Focus()

	case "f5":
		m.errMsg = ""
		m.infoMsg = ""
		return m, m.refreshCmd()

	case "t", "T":
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		m.errMsg = ""
		m.infoMsg = ""
		m.state = stateTagInput
		m.tagInput.SetValue("")
		return m, textinput.Blink

	case "e", "E":
		m.errMsg = ""
		m.infoMsg = ""
		cmd := BuildEditCmd(m.cfgPath)
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			return editDoneMsg{err: err}
		})

	case "c", "C":
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		m.state = stateCommitSelect
		m.commitOutput = ""
		m.commitErr = nil
		m.commitCursor = 0
		return m, nil

	case "p", "P":
		m.errMsg = ""
		m.infoMsg = ""
		m.state = stateProviderManage
		m.providerCursor = 0
		m.providerDelMode = false
		return m, nil

	case "a", "A":
		m.errMsg = ""
		m.infoMsg = ""
		m.state = stateAccountManage
		m.accountCursor = 0
		m.delMode = false
		return m, nil

	case "s", "S":
		m.errMsg = ""
		m.infoMsg = ""
		m.state = stateSettings
		m.settingsCursor = 0
		return m, nil

	case "g", "G":
		m.errMsg = ""
		m.infoMsg = ""
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		rc := LoadRepoConfig()
		m.repoConfig = rc
		m.repoConfigProviders = m.cfg.Providers
		m.repoConfigKeys = make([]string, 0, len(m.cfg.Providers))
		for _, p := range m.cfg.Providers {
			m.repoConfigKeys = append(m.repoConfigKeys, p.ID)
		}
		m.repoConfigCursor = 0
		m.state = stateRepoConfig
		return m, nil

	case "y", "Y":
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_copy")
			return m, nil
		}
		text := fmt.Sprintf("%s <%s>", m.localName, m.localEmail)
		if m.remoteURL != "" {
			text += "\n" + m.remoteURL
		}
		return m, m.copyCmd(text)

	case "b", "B":
		if !m.isRepo {
			m.errMsg = m.tr.Tr("msg_cannot_commit")
			return m, nil
		}
		bd := ListBranches()
		if len(bd) == 0 {
			m.errMsg = m.tr.Tr("msg_no_branches")
			return m, nil
		}
		m.branches = bd
		m.branchCursor = 0
		for i, b := range bd {
			if b == m.branch {
				m.branchCursor = i
				break
			}
		}
		m.state = stateBranchSwitch
		return m, nil

	default:
		for i := range m.cfg.Profiles {
			key := fmt.Sprintf("%d", i+1)
			if msg.String() == key {
				m.cursor = i
				return m.switchToCursor()
			}
		}
	}
	return m, nil
}

func (m model) switchToCursor() (tea.Model, tea.Cmd) {
	if !m.isRepo {
		m.errMsg = m.tr.Tr("msg_cannot_switch")
		return m, nil
	}
	if len(m.cfg.Profiles) == 0 || m.cursor < 0 || m.cursor >= len(m.cfg.Profiles) {
		m.errMsg = m.tr.Tr("msg_no_profile")
		return m, nil
	}
	m.errMsg = ""
	m.infoMsg = ""
	profile := m.cfg.Profiles[m.cursor]
	return m, func() tea.Msg {
		log, err := SwitchProfile(&profile, m.cfg.Providers, m.tr)
		return switchDoneMsg{log: log, err: err}
	}
}

func (m model) handleCommitSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateNormal
		m.infoMsg = ""
		m.errMsg = ""
		return m, nil

	case "up", "k":
		m.commitCursor--
		if m.commitCursor < 0 {
			m.commitCursor = len(commitTypes) - 1
		}
		return m, nil

	case "down", "j":
		m.commitCursor++
		if m.commitCursor >= len(commitTypes) {
			m.commitCursor = 0
		}
		return m, nil

	case "enter":
		m.commitType = commitTypes[m.commitCursor].key
		m.state = stateCommitInput
		m.commitInput.SetValue("")
		return m, m.commitInput.Focus()

	case "f":
		m.commitType = "feat"
	case "x":
		m.commitType = "fix"
	case "d":
		m.commitType = "docs"
	case "r":
		m.commitType = "refactor"
	default:
		return m, nil
	}
	m.state = stateCommitInput
	m.commitInput.SetValue("")
	return m, m.commitInput.Focus()
}

func (m model) handleCommitInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.state = stateNormal
		m.commitInput.Blur()
		m.infoMsg = ""
		m.errMsg = ""
		return m, nil
	case msg.Type == tea.KeyF2:
		desc := strings.TrimSpace(m.commitInput.Value())
		if desc == "" {
			m.errMsg = m.tr.Tr("msg_desc_empty")
			return m, nil
		}
		m.state = stateCommitOutput
		m.commitInput.Blur()
		m.infoMsg = ""
		m.errMsg = ""
		m.commitOutput = "..."
		ct := m.commitType
		tr := m.tr
		return m, func() tea.Msg {
			output, err := RunConventionalCommit(ct, desc, tr)
			return commitResultMsg{output: output, err: err}
		}
	default:
		var cmd tea.Cmd
		m.commitInput, cmd = m.commitInput.Update(msg)
		return m, cmd
	}
}

func (m model) handlePRInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		title := strings.TrimSpace(m.prTitleInput.Value())
		if title == "" {
			m.errMsg = m.tr.Tr("pr_title_empty")
			return m, nil
		}
		m.state = statePROutput
		m.prOutput = "..."
		tr := m.tr
		return m, func() tea.Msg {
			output, err := CreatePR(title, tr)
			return prResultMsg{output: output, err: err}
		}
	case "esc":
		m.state = stateNormal
		return m, nil
	default:
		var cmd tea.Cmd
		m.prTitleInput, cmd = m.prTitleInput.Update(msg)
		return m, cmd
	}
}

func (m model) handleMouseEvent(msg tea.MouseMsg) model {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		switch m.state {
		case stateNormal:
			if m.cursor > 0 {
				m.cursor--
			}
		case stateCommitSelect:
			if m.commitCursor > 0 {
				m.commitCursor--
			}
		case stateAccountManage:
			if m.accountCursor > 0 {
				m.accountCursor--
			}
		case stateAccountForm:
			if m.accountFocused > 0 {
				m.accountInputs[m.accountFocused].Blur()
				m.accountFocused--
				m.accountInputs[m.accountFocused].Focus()
			}
		case stateProviderManage:
			if m.providerCursor > 0 {
				m.providerCursor--
			}
		case stateProviderForm:
			if m.providerFocused > 0 {
				m.providerInputs[m.providerFocused].Blur()
				m.providerFocused--
				m.providerInputs[m.providerFocused].Focus()
			}
		case stateProviderPicker:
			if m.providerPickCursor > 0 {
				m.providerPickCursor--
			}
		case stateSSHKeyPicker:
			if m.sshKeyCursor > 0 {
				m.sshKeyCursor--
			}
		case stateSettings:
			if m.settingsCursor > 0 {
				m.settingsCursor--
			}
		case stateRepoConfig:
			if m.repoConfigCursor > 0 {
				m.repoConfigCursor--
			}
		case stateBranchSwitch:
			if m.branchCursor > 0 {
				m.branchCursor--
			}
		}
	case tea.MouseButtonWheelDown:
		switch m.state {
		case stateNormal:
			max := len(m.cfg.Profiles) + 7
			if m.cursor < max-1 {
				m.cursor++
			}
		case stateCommitSelect:
			if m.commitCursor < len(commitTypes)-1 {
				m.commitCursor++
			}
		case stateAccountManage:
			if m.accountCursor < len(m.cfg.Profiles)-1 {
				m.accountCursor++
			}
		case stateAccountForm:
			if m.accountFocused < len(m.accountInputs)-1 {
				m.accountInputs[m.accountFocused].Blur()
				m.accountFocused++
				m.accountInputs[m.accountFocused].Focus()
			}
		case stateProviderManage:
			if m.providerCursor < len(m.cfg.Providers)-1 {
				m.providerCursor++
			}
		case stateProviderForm:
			if m.providerFocused < len(m.providerInputs)-1 {
				m.providerInputs[m.providerFocused].Blur()
				m.providerFocused++
				m.providerInputs[m.providerFocused].Focus()
			}
		case stateProviderPicker:
			if m.providerPickCursor < len(m.providerPickList)-1 {
				m.providerPickCursor++
			}
		case stateSSHKeyPicker:
			if m.sshKeyCursor < len(m.sshKeys)-1 {
				m.sshKeyCursor++
			}
		case stateSettings:
			if m.settingsCursor < len(m.settingsItems())-1 {
				m.settingsCursor++
			}
		case stateRepoConfig:
			if m.repoConfigCursor < len(m.repoConfigKeys)-1 {
				m.repoConfigCursor++
			}
		case stateBranchSwitch:
			if m.branchCursor < len(m.branches)-1 {
				m.branchCursor++
			}
		}
	}
	return m
}

func (m model) handleAccountManageKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateNormal
		m.infoMsg = ""
		m.errMsg = ""
		return m, nil

	case "up", "k":
		if len(m.cfg.Profiles) == 0 {
			return m, nil
		}
		m.accountCursor--
		if m.accountCursor < 0 {
			m.accountCursor = len(m.cfg.Profiles) - 1
		}
		return m, nil

	case "down", "j":
		if len(m.cfg.Profiles) == 0 {
			return m, nil
		}
		m.accountCursor++
		if m.accountCursor >= len(m.cfg.Profiles) {
			m.accountCursor = 0
		}
		return m, nil

	case "enter":
		if m.delMode {
			if m.accountCursor >= 0 && m.accountCursor < len(m.cfg.Profiles) {
				return m.confirmDeleteProfile(m.accountCursor)
			}
			return m, nil
		}
		if m.accountCursor >= 0 && m.accountCursor < len(m.cfg.Profiles) {
			return m.startAccountForm(m.accountCursor)
		}
		return m, nil

	case "n", "N":
		return m.startAccountForm(-1)

	case "d", "D":
		m.delMode = !m.delMode
		return m, nil

	default:
		for i := range m.cfg.Profiles {
			key := fmt.Sprintf("%d", i+1)
			if msg.String() == key {
				m.accountCursor = i
				if m.delMode {
					return m.confirmDeleteProfile(i)
				}
				return m.startAccountForm(i)
			}
		}
	}
	return m, nil
}

func (m model) confirmDeleteProfile(idx int) (tea.Model, tea.Cmd) {
	name := m.cfg.Profiles[idx].Name
	m.cfg.Profiles = append(m.cfg.Profiles[:idx], m.cfg.Profiles[idx+1:]...)
	m.delMode = false
	if m.accountCursor >= len(m.cfg.Profiles) {
		m.accountCursor = len(m.cfg.Profiles) - 1
	}
	if m.accountCursor < 0 {
		m.accountCursor = 0
	}
	if m.cursor >= len(m.cfg.Profiles) {
		m.cursor = len(m.cfg.Profiles) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	return m, func() tea.Msg {
		err := SaveConfig(m.cfg, m.cfgPath)
		return accountDeletedMsg{err: err, name: name}
	}
}

func (m model) startAccountForm(idx int) (tea.Model, tea.Cmd) {
	m.accountEditIdx = idx
	m.delMode = false
	inputs := make([]textinput.Model, 6)
	placeholders := []string{
		m.tr.Tr("account_id_ph"),
		m.tr.Tr("account_name_ph"),
		m.tr.Tr("account_gitname_ph"),
		m.tr.Tr("account_gitemail_ph"),
		m.tr.Tr("account_provider_ph"),
		m.tr.Tr("account_ssh_identity_ph"),
	}
	for i := range inputs {
		ti := textinput.New()
		ti.Placeholder = placeholders[i]
		ti.CharLimit = 100
		inputs[i] = ti
	}

	if idx >= 0 && idx < len(m.cfg.Profiles) {
		p := m.cfg.Profiles[idx]
		inputs[0].SetValue(p.ID)
		inputs[1].SetValue(p.Name)
		inputs[2].SetValue(p.GitName)
		inputs[3].SetValue(p.GitEmail)
		if p.ProviderID != "" {
			for _, prov := range m.cfg.Providers {
				if prov.ID == p.ProviderID {
					inputs[4].SetValue(prov.Name)
					break
				}
			}
		}
		inputs[5].SetValue(p.SSHIdentityFile)
	}

	m.accountFocused = 0
	inputs[0].Focus()
	m.accountInputs = inputs
	m.state = stateAccountForm
	return m, nil
}

func (m model) handleAccountFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateAccountManage
		return m, nil

	case "tab", "down", "j":
		m.accountInputs[m.accountFocused].Blur()
		m.accountFocused = (m.accountFocused + 1) % len(m.accountInputs)
		m.accountInputs[m.accountFocused].Focus()
		return m, nil

	case "shift+tab", "up", "k":
		m.accountInputs[m.accountFocused].Blur()
		m.accountFocused--
		if m.accountFocused < 0 {
			m.accountFocused = len(m.accountInputs) - 1
		}
		m.accountInputs[m.accountFocused].Focus()
		return m, nil

	case "ctrl+o":
		if m.accountFocused == 5 {
			return m.startSSHKeyScan()
		}
		return m, nil

	case "ctrl+p":
		if m.accountFocused == 4 {
			return m.startProviderPicker()
		}
		return m, nil

	case "enter":
		values := make([]string, 6)
		for i, input := range m.accountInputs {
			values[i] = strings.TrimSpace(input.Value())
		}
		if values[0] == "" || values[2] == "" || values[3] == "" {
			m.errMsg = m.tr.Tr("msg_fields_required")
			return m, nil
		}

		for i, p := range m.cfg.Profiles {
			if p.ID == values[0] && i != m.accountEditIdx {
				m.errMsg = m.tr.Tr("msg_duplicate_id", values[0])
				return m, nil
			}
		}

		providerID := ""
		for _, prov := range m.cfg.Providers {
			if prov.Name == values[4] {
				providerID = prov.ID
				break
			}
		}

		profile := Profile{
			ID:              values[0],
			Name:            values[1],
			GitName:         values[2],
			GitEmail:        values[3],
			ProviderID:      providerID,
			SSHIdentityFile: values[5],
		}

		if m.accountEditIdx >= 0 && m.accountEditIdx < len(m.cfg.Profiles) {
			m.cfg.Profiles[m.accountEditIdx] = profile
			m.accountCursor = m.accountEditIdx
		} else {
			m.cfg.Profiles = append(m.cfg.Profiles, profile)
			m.accountCursor = len(m.cfg.Profiles) - 1
		}
		m.cursor = m.accountCursor

		m.state = stateAccountManage
		return m, m.saveConfigCmd(profile.Name)

	default:
		if m.accountFocused >= 0 && m.accountFocused < len(m.accountInputs) {
			var cmd tea.Cmd
			m.accountInputs[m.accountFocused], cmd = m.accountInputs[m.accountFocused].Update(msg)
			return m, cmd
		}
		return m, nil
	}
}

func (m model) startSSHKeyScan() (tea.Model, tea.Cmd) {
	keys := scanSSHKeys()
	if len(keys) == 0 {
		m.errMsg = m.tr.Tr("ssh_scan_none")
		return m, nil
	}
	m.sshKeys = keys
	m.sshKeyCursor = 0
	m.state = stateSSHKeyPicker
	return m, nil
}

func (m model) handleSSHKeyPickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateAccountForm
		m.sshKeys = nil
		return m, nil
	case "up", "k":
		if m.sshKeyCursor > 0 {
			m.sshKeyCursor--
		}
		return m, nil
	case "down", "j":
		if m.sshKeyCursor < len(m.sshKeys)-1 {
			m.sshKeyCursor++
		}
		return m, nil
	case "enter":
		m.accountInputs[5].SetValue(m.sshKeys[m.sshKeyCursor])
		m.state = stateAccountForm
		m.sshKeys = nil
		return m, nil
	}
	return m, nil
}

func (m model) handleTagInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		version := m.tagInput.Value()
		if version == "" {
			m.errMsg = m.tr.Tr("tag_version_empty")
			return m, nil
		}
		m.state = stateNormal
		return m, m.createTagCmd(version)
	case "esc":
		m.state = stateNormal
		m.commitOutput = ""
		m.commitErr = nil
		return m, nil
	default:
		var cmd tea.Cmd
		m.tagInput, cmd = m.tagInput.Update(msg)
		return m, cmd
	}
}

func (m model) handleBranchSwitchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateNormal
		m.infoMsg = ""
		m.errMsg = ""
		return m, nil
	case "up", "k":
		if m.branchCursor > 0 {
			m.branchCursor--
		}
		return m, nil
	case "down", "j":
		if m.branchCursor < len(m.branches)-1 {
			m.branchCursor++
		}
		return m, nil
	case "enter":
		selected := m.branches[m.branchCursor]
		if selected == m.branch {
			m.state = stateNormal
			return m, nil
		}
		out, err := SwitchBranch(selected)
		if err != nil {
			m.errMsg = m.tr.Tr("msg_branch_switch_fail", strings.TrimSpace(out))
			m.state = stateNormal
			return m, nil
		}
		m.infoMsg = m.tr.Tr("msg_branch_switched", selected)
		m.state = stateNormal
		m.branch = GetCurrentBranch()
		return m, m.refreshCmd()
	}
	return m, nil
}

func (m model) handleProviderManageKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	providers := m.cfg.Providers
	count := len(providers)
	switch msg.String() {
	case "esc":
		m.state = stateNormal
		m.infoMsg = ""
		m.errMsg = ""
		return m, nil

	case "up", "k":
		if count == 0 {
			return m, nil
		}
		m.providerCursor--
		if m.providerCursor < 0 {
			m.providerCursor = count - 1
		}
		return m, nil

	case "down", "j":
		if count == 0 {
			return m, nil
		}
		m.providerCursor++
		if m.providerCursor >= count {
			m.providerCursor = 0
		}
		return m, nil

	case "enter":
		if m.providerDelMode {
			return m.confirmDeleteProvider()
		}
		return m.startProviderForm(m.providerCursor)

	case "n", "N":
		return m.startProviderForm(-1)

	case "d", "D":
		m.providerDelMode = !m.providerDelMode
		return m, nil

	default:
		for i := range providers {
			key := fmt.Sprintf("%d", i+1)
			if msg.String() == key {
				m.providerCursor = i
				if m.providerDelMode {
					return m.confirmDeleteProvider()
				}
				return m.startProviderForm(i)
			}
		}
	}
	return m, nil
}

func (m model) confirmDeleteProvider() (tea.Model, tea.Cmd) {
	providers := m.cfg.Providers
	if m.providerCursor < 0 || m.providerCursor >= len(providers) {
		return m, nil
	}
	name := providers[m.providerCursor].Name
	m.cfg.Providers = append(providers[:m.providerCursor], providers[m.providerCursor+1:]...)
	m.providerDelMode = false
	if m.providerCursor >= len(m.cfg.Providers) {
		m.providerCursor = len(m.cfg.Providers) - 1
	}
	return m, func() tea.Msg {
		err := SaveConfig(m.cfg, m.cfgPath)
		return providerDeletedMsg{err: err, provider: name}
	}
}

func (m model) startProviderForm(idx int) (tea.Model, tea.Cmd) {
	m.providerEditIdx = idx
	m.providerDelMode = false
	inputs := make([]textinput.Model, 3)
	placeholders := []string{
		m.tr.Tr("provider_id_ph"),
		m.tr.Tr("provider_name_ph"),
		m.tr.Tr("provider_host_ph"),
	}
	for i := range inputs {
		ti := textinput.New()
		ti.Placeholder = placeholders[i]
		ti.CharLimit = 100
		inputs[i] = ti
	}

	if idx >= 0 && idx < len(m.cfg.Providers) {
		p := m.cfg.Providers[idx]
		inputs[0].SetValue(p.ID)
		inputs[1].SetValue(p.Name)
		inputs[2].SetValue(p.Host)
	}

	m.providerFocused = 0
	inputs[0].Focus()
	m.providerInputs = inputs
	m.state = stateProviderForm
	return m, nil
}

func (m model) handleProviderFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateProviderManage
		return m, nil

	case "tab", "down", "j":
		m.providerInputs[m.providerFocused].Blur()
		m.providerFocused = (m.providerFocused + 1) % len(m.providerInputs)
		m.providerInputs[m.providerFocused].Focus()
		return m, nil

	case "shift+tab", "up", "k":
		m.providerInputs[m.providerFocused].Blur()
		m.providerFocused--
		if m.providerFocused < 0 {
			m.providerFocused = len(m.providerInputs) - 1
		}
		m.providerInputs[m.providerFocused].Focus()
		return m, nil

	case "enter":
		id := strings.TrimSpace(m.providerInputs[0].Value())
		name := strings.TrimSpace(m.providerInputs[1].Value())
		host := strings.TrimSpace(m.providerInputs[2].Value())
		if id == "" || name == "" || host == "" {
			m.errMsg = m.tr.Tr("msg_provider_fields_required")
			return m, nil
		}

		for i, p := range m.cfg.Providers {
			if p.ID == id && i != m.providerEditIdx {
				m.errMsg = m.tr.Tr("msg_duplicate_id", id)
				return m, nil
			}
		}

		if m.providerEditIdx >= 0 && m.providerEditIdx < len(m.cfg.Providers) {
			m.cfg.Providers[m.providerEditIdx] = Provider{ID: id, Name: name, Host: host}
		} else {
			m.cfg.Providers = append(m.cfg.Providers, Provider{ID: id, Name: name, Host: host})
		}

		m.state = stateProviderManage
		return m, func() tea.Msg {
			err := SaveConfig(m.cfg, m.cfgPath)
			return providerSavedMsg{err: err}
		}

	default:
		if m.providerFocused >= 0 && m.providerFocused < len(m.providerInputs) {
			var cmd tea.Cmd
			m.providerInputs[m.providerFocused], cmd = m.providerInputs[m.providerFocused].Update(msg)
			return m, cmd
		}
		return m, nil
	}
}

func (m model) startProviderPicker() (tea.Model, tea.Cmd) {
	if len(m.cfg.Providers) == 0 {
		m.errMsg = m.tr.Tr("provider_pick_none")
		return m, nil
	}
	m.providerPickList = m.cfg.Providers
	m.providerPickCursor = 0
	m.state = stateProviderPicker
	return m, nil
}

func (m model) handleProviderPickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	list := m.providerPickList
	switch msg.String() {
	case "esc":
		m.state = stateAccountForm
		m.providerPickList = nil
		return m, nil
	case "up", "k":
		if m.providerPickCursor > 0 {
			m.providerPickCursor--
		}
		return m, nil
	case "down", "j":
		if m.providerPickCursor < len(list)-1 {
			m.providerPickCursor++
		}
		return m, nil
	case "enter":
		if m.providerPickCursor >= 0 && m.providerPickCursor < len(list) {
			selected := list[m.providerPickCursor]
			m.accountInputs[4].SetValue(selected.Name)
		}
		m.state = stateAccountForm
		m.providerPickList = nil
		return m, nil
	}
	return m, nil
}

type settingsItem struct {
	label string
	value string
}

func (m model) settingsItems() []settingsItem {
	return []settingsItem{
		{m.tr.Tr("settings_language"), string(m.cfg.Lang)},
		{m.tr.Tr("settings_edit_config"), ""},
	}
}

func (m model) handleSettingsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	items := m.settingsItems()
	switch msg.String() {
	case "up", "k":
		m.settingsCursor--
		if m.settingsCursor < 0 {
			m.settingsCursor = len(items) - 1
		}
		return m, nil

	case "down", "j":
		m.settingsCursor++
		if m.settingsCursor >= len(items) {
			m.settingsCursor = 0
		}
		return m, nil

	case "left", "right", "tab":
		if m.settingsCursor == 0 {
			curIdx := 0
			for i, l := range LangList {
				if Lang(m.cfg.Lang) == l {
					curIdx = i
					break
				}
			}
			var nextIdx int
			if msg.String() == "left" {
				nextIdx = (curIdx - 1 + len(LangList)) % len(LangList)
			} else {
				nextIdx = (curIdx + 1) % len(LangList)
			}
			m.cfg.Lang = string(LangList[nextIdx])
			lang := Lang(m.cfg.Lang)
			m.tr = NewTranslator(lang)
			m.commitInput.Placeholder = m.tr.Tr("prompt_enter_desc_ph")
			return m, func() tea.Msg {
				err := SaveConfig(m.cfg, m.cfgPath)
				return settingsSavedMsg{err: err}
			}
		}
	case "enter":
		if m.settingsCursor == 1 {
			cmd := BuildEditCmd(m.cfgPath)
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return editDoneMsg{err: err}
			})
		}
		return m, nil

	case "esc":
		m.state = stateNormal
		m.infoMsg = ""
		m.errMsg = ""
		return m, nil
	}
	return m, nil
}

func (m model) renderSettings(s *strings.Builder) {
	var buf strings.Builder
	buf.WriteString(labelStyle.Render(m.tr.Tr("label_settings")) + "\n\n")

	items := m.settingsItems()
	for i, item := range items {
		cursor := "  "
		if i == m.settingsCursor {
			cursor = cursorStyle.Render("▸ ")
		}
		var val string
		if item.label == m.tr.Tr("settings_language") {
			var parts []string
			for i, l := range LangList {
				name := m.tr.Tr("settings_lang_" + langKeys[l])
				if Lang(m.cfg.Lang) == l {
					parts = append(parts, "▸ "+name)
				} else {
					parts = append(parts, "  "+name)
				}
				_ = i
			}
			val = strings.Join(parts, "  ")
		} else if item.label == m.tr.Tr("settings_edit_config") {
			val = dimStyle.Render("[Enter]")
		} else {
			val = item.value
		}
		buf.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, labelStyle.Render(item.label), val))
	}

	s.WriteString(m.centerBlock(buf.String()) + "\n")
}

func (m model) buildSettingsBar() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[←→]"), m.tr.Tr("action_nav")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Esc]"), m.tr.Tr("action_back")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Q]"), m.tr.Tr("action_quit")))
	return b.String()
}

func (m model) buildSSHKeyPickerBar() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[↑↓]"), m.tr.Tr("action_nav")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Enter]"), m.tr.Tr("action_select")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Esc]"), m.tr.Tr("action_back")))
	return b.String()
}

func (m model) buildRepoConfigBar() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[↑↓]"), m.tr.Tr("action_nav")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Enter]"), m.tr.Tr("action_edit_btn")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Esc]"), m.tr.Tr("action_back")))
	return b.String()
}

func (m model) buildRepoConfigEditBar() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Enter]"), m.tr.Tr("action_save")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Esc]"), m.tr.Tr("action_cancel")))
	return b.String()
}

func (m model) handleRepoConfigKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.repoConfigCursor > 0 {
			m.repoConfigCursor--
		}
		return m, nil
	case "down", "j":
		if m.repoConfigCursor < len(m.repoConfigKeys)-1 {
			m.repoConfigCursor++
		}
		return m, nil
	case "enter":
		if len(m.repoConfigKeys) > 0 {
			key := m.repoConfigKeys[m.repoConfigCursor]
			ti := textinput.New()
			ti.Placeholder = m.tr.Tr("repo_config_path_ph")
			ti.SetValue(m.repoConfig.Paths[key])
			ti.CharLimit = 200
			ti.Focus()
			m.repoConfigInput = ti
			m.state = stateRepoConfigEdit
		}
		return m, nil
	case "esc":
		m.state = stateNormal
		m.infoMsg = ""
		m.errMsg = ""
		return m, nil
	}
	return m, nil
}

func (m model) handleRepoConfigEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		key := m.repoConfigKeys[m.repoConfigCursor]
		val := strings.TrimSpace(m.repoConfigInput.Value())
		if val == "" {
			delete(m.repoConfig.Paths, key)
		} else {
			m.repoConfig.Paths[key] = val
		}
		if err := SaveRepoConfig(m.repoConfig); err != nil {
			m.errMsg = fmt.Sprintf(m.tr.Tr("repo_config_save_fail"), err)
		} else {
			m.infoMsg = m.tr.Tr("repo_config_saved")
		}
		m.state = stateRepoConfig
		return m, nil
	case "esc":
		m.state = stateRepoConfig
		return m, nil
	case "tab", "down":
		var cmd tea.Cmd
		m.repoConfigInput, cmd = m.repoConfigInput.Update(msg)
		return m, cmd
	case "shift+tab", "up":
		var cmd tea.Cmd
		m.repoConfigInput, cmd = m.repoConfigInput.Update(msg)
		return m, cmd
	default:
		var cmd tea.Cmd
		m.repoConfigInput, cmd = m.repoConfigInput.Update(msg)
		return m, cmd
	}
}

func (m model) renderRepoConfig(s *strings.Builder) {
	var buf strings.Builder
	buf.WriteString(labelStyle.Render(m.tr.Tr("label_repo_config")) + "\n\n")

	if len(m.repoConfigKeys) == 0 {
		buf.WriteString("  " + dimStyle.Render(m.tr.Tr("repo_config_empty")) + "\n")
	} else {
		for i, key := range m.repoConfigKeys {
			cursor := "  "
			if i == m.repoConfigCursor {
				cursor = cursorStyle.Render("▸ ")
			}
			var provName string
			for _, p := range m.cfg.Providers {
				if p.ID == key {
					provName = p.Name + " (" + p.Host + ")"
					break
				}
			}
			if provName == "" {
				provName = key
			}

			if m.state == stateRepoConfigEdit && i == m.repoConfigCursor {
				buf.WriteString(fmt.Sprintf("  %s%s\n", cursor, labelStyle.Render(provName)))
				buf.WriteString(fmt.Sprintf("      %s %s\n", dimStyle.Render(m.tr.Tr("repo_config_path")+":"), m.repoConfigInput.View()))
			} else {
				path := m.repoConfig.Paths[key]
				if path == "" {
					path = dimStyle.Render(m.tr.Tr("repo_config_not_set"))
				}
				buf.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, labelStyle.Render(provName), valueStyle.Render(path)))
			}
		}
	}

	s.WriteString(m.centerBlock(buf.String()) + "\n")
}

var (
	cyan   = lipgloss.Color("#00FFFF")
	yellow = lipgloss.Color("#FFD700")
	red    = lipgloss.Color("#FF4444")
	green  = lipgloss.Color("#00FF7F")
	gray   = lipgloss.Color("#888888")
	white  = lipgloss.Color("#FFFFFF")
	orange = lipgloss.Color("#FFA500")

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(lipgloss.Color("#1a1b2e")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2d2d44")).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(white)

	keyStyle = lipgloss.NewStyle().
			Foreground(orange).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(green)

	dimStyle = lipgloss.NewStyle().
			Foreground(gray)

	faintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	dimValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#999999"))
)

func (m model) boxW() int {
	bw := m.width - 8
	if bw < 40 {
		bw = 40
	}
	if bw > 80 {
		bw = 80
	}
	return bw
}

func (m model) padLines(content string) string {
	bw := m.boxW()
	lines := strings.Split(content, "\n")
	halfPad := (m.width - bw) / 2
	if halfPad < 0 {
		halfPad = 0
	}
	leftPad := strings.Repeat(" ", halfPad)
	for i, line := range lines {
		w := lipgloss.Width(line)
		if w < bw {
			line = line + strings.Repeat(" ", bw-w)
		}
		lines[i] = leftPad + line
	}
	return strings.Join(lines, "\n")
}

func (m model) centerBlock(content string) string {
	return m.padLines(content)
}

func (m model) separator() string {
	bw := m.boxW()
	line := dimStyle.Render(strings.Repeat("─", bw))
	return "\n" + m.padLines(line) + "\n"
}

func (m model) statusLine(label, value string) string {
	return fmt.Sprintf(" %s %s %s", labelStyle.Render("•"), dimStyle.Render(label), valueStyle.Render(value))
}

func (m model) View() string {
	if m.state != stateCommitOutput {
		m.commitOutput = ""
		m.commitErr = nil
	}
	var content strings.Builder
	m.renderHeader(&content)

	if !m.repoChecked {
		content.WriteString("\n\n" + m.centerBlock(dimStyle.Render(m.tr.Tr("msg_checking"))))
		return content.String()
	}

	content.WriteString("\n\n")

	switch m.state {
	case stateCommitSelect, stateCommitInput, stateCommitOutput, stateTagInput, statePRInput, statePROutput:
		m.renderCommitFlow(&content)
	case stateAccountManage:
		m.renderAccountManage(&content)
	case stateAccountForm:
		m.renderAccountForm(&content)
	case stateProviderManage:
		m.renderProviderManage(&content)
	case stateProviderForm:
		m.renderProviderForm(&content)
	case stateProviderPicker:
		m.renderProviderPicker(&content)
	case stateSSHKeyPicker:
		m.renderSSHKeyPicker(&content)
	case stateSettings:
		m.renderSettings(&content)
	case stateRepoConfig, stateRepoConfigEdit:
		m.renderRepoConfig(&content)
	case stateBranchSwitch:
		m.renderBranchSwitch(&content)
	default:
		m.renderNormalContent(&content)
	}

	// Build bottom bar (no messages — they go in a dedicated zone above)
	var bottom strings.Builder
	m.renderBottomBar(&bottom)
	bottomStr := bottom.String()

	// Build message zone (max 3 lines, right-aligned)
	var msgLines []string
	if m.errMsg != "" {
		for _, l := range strings.Split(m.errMsg, "\n") {
			if l = strings.TrimSpace(l); l != "" {
				msgLines = append(msgLines, l)
			}
		}
	} else if m.infoMsg != "" {
		for _, l := range strings.Split(m.infoMsg, "\n") {
			if l = strings.TrimSpace(l); l != "" {
				msgLines = append(msgLines, l)
			}
		}
	}
	if len(msgLines) > 4 {
		msgLines = msgLines[len(msgLines)-4:]
	}
	msgHeight := len(msgLines)

	// Pad content to push bottom bar to window bottom
	if m.height > 0 {
		contentLines := lipgloss.Height(content.String())
		bottomLines := lipgloss.Height(bottomStr)
		remaining := m.height - contentLines - bottomLines - msgHeight
		if remaining < 0 {
			remaining = 0
		}
		for i := 0; i < remaining; i++ {
			content.WriteString("\n")
		}
	}

	// Render message zone (right-aligned, above bottom bar)
	for _, l := range msgLines {
		styled := successStyle.Render(l)
		if m.errMsg != "" {
			styled = errorStyle.Render(l)
		}
		pad := m.width - lipgloss.Width(styled)
		if pad < 0 {
			pad = 0
		}
		content.WriteString(strings.Repeat(" ", pad) + styled + "\n")
	}

	content.WriteString(bottomStr)
	return content.String()
}

func (m model) renderHeader(s *strings.Builder) {
	sep := faintStyle.Render(" │ ")

	// Left: app title + branch
	var leftParts []string
	leftParts = append(leftParts, m.tr.Tr("app_title"))
	if m.repoChecked && m.isRepo && m.branch != "" {
		b := m.branch
		if m.isDirty {
			b += fmt.Sprintf(" [%d]", m.dirtyCount)
		}
		leftParts = append(leftParts, b)
	} else if m.repoChecked && !m.isRepo {
		leftParts = append(leftParts, m.tr.Tr("label_no_repo"))
	} else {
		leftParts = append(leftParts, "...")
	}
	leftStr := strings.Join(leftParts, sep)

	// Right: remote URL + identity
	var rightParts []string
	if m.repoChecked && m.isRepo {
		if m.localName != "" {
			rightParts = append(rightParts, m.localName)
		}
	}
	rightStr := strings.Join(rightParts, sep)

	// Spacer to push left and right to edges
	leftW := lipgloss.Width(leftStr)
	rightW := lipgloss.Width(rightStr)
	contentW := m.width - 2 // minus padding left+right
	spacerW := contentW - leftW - rightW
	if spacerW < 1 {
		spacerW = 1
	}
	spacer := strings.Repeat(" ", spacerW)

	line := leftStr + spacer + rightStr
	s.WriteString(headerStyle.Width(m.width).Render(line) + "\n")
}

func (m model) renderNormalContent(s *strings.Builder) {
	if !m.isRepo {
		s.WriteString(m.centerBlock(
			errorStyle.Render(m.tr.Tr("msg_not_repo")) + "\n\n" +
				dimStyle.Render(m.tr.Tr("msg_press_e")),
		) + "\n")
		return
	}

	var buf strings.Builder
	buf.WriteString(labelStyle.Render(m.tr.Tr("label_status")) + "\n")
	buf.WriteString(m.statusLine(m.tr.Tr("label_branch"), m.branchDisplay()) + "\n")
	buf.WriteString(m.statusLine(m.tr.Tr("label_route"), m.routeDisplay()) + "\n")
	buf.WriteString(m.statusLine(m.tr.Tr("label_identity"), m.identityDisplay()))
	s.WriteString(m.centerBlock(buf.String()) + "\n")

	s.WriteString(m.separator())

	var ab strings.Builder
	ab.WriteString(labelStyle.Render(m.tr.Tr("label_actions")) + "\n" + m.actionsView())
	s.WriteString(m.centerBlock(ab.String()) + "\n")
}

func (m model) renderCommitFlow(s *strings.Builder) {
	switch m.state {
	case stateCommitSelect:
		var buf strings.Builder
		buf.WriteString(labelStyle.Render(" "+m.tr.Tr("prompt_select_type")) + "\n\n")
		for i, ct := range commitTypes {
			cursor := "  "
			if i == m.commitCursor {
				cursor = cursorStyle.Render("▸ ")
			}
			buf.WriteString(fmt.Sprintf("  %s%s  %s\n", cursor, keyStyle.Render(ct.key), dimValueStyle.Render(ct.desc)))
		}
		buf.WriteString("\n" + dimStyle.Render("  "+m.tr.Tr("prompt_cancel")))
		s.WriteString(m.centerBlock(buf.String()))

	case stateCommitInput:
		s.WriteString(m.centerBlock(
			labelStyle.Render(" "+m.tr.Tr("commit_header", m.commitType))+"\n\n"+
				"  "+m.commitInput.View()+"\n\n"+
				dimStyle.Render("  "+m.tr.Tr("commit_confirm_hint")),
		))

	case statePRInput:
		s.WriteString(m.centerBlock(
			labelStyle.Render(" "+m.tr.Tr("pr_result"))+"\n\n"+
				"  "+m.prTitleInput.View()+"\n\n"+
				dimStyle.Render("  "+m.tr.Tr("prompt_confirm_cancel")),
		))

	case stateCommitOutput:
		out := m.commitOutput
		if m.commitErr != nil {
			out += "\n" + errorStyle.Render(m.tr.Tr("cli_error", m.commitErr))
		}
		s.WriteString(m.centerBlock(
			labelStyle.Render(" "+m.tr.Tr("commit_result"))+"\n\n"+
				"  "+out+"\n\n"+
				dimStyle.Render("  "+m.tr.Tr("prompt_any_key")+"  |  "+
					keyStyle.Render("[R]")+" "+m.tr.Tr("pr_hint")),
		))

	case statePROutput:
		out := m.prOutput
		if m.prErr != nil {
			out += "\n" + errorStyle.Render(m.tr.Tr("cli_error", m.prErr))
		}
		s.WriteString(m.centerBlock(
			labelStyle.Render(" "+m.tr.Tr("pr_result"))+"\n\n"+
				"  "+out+"\n\n"+
				dimStyle.Render("  "+m.tr.Tr("prompt_any_key")),
		))

	case stateTagInput:
		var buf strings.Builder
		buf.WriteString(labelStyle.Render(" "+m.tr.Tr("tag_result"))+"\n\n")
		if m.commitOutput != "" {
			buf.WriteString("  " + m.commitOutput + "\n\n")
		}
		buf.WriteString(labelStyle.Render(" "+m.tr.Tr("tag_input_ph"))+"\n\n"+
			"  "+m.tagInput.View()+"\n\n"+
			dimStyle.Render("  "+m.tr.Tr("prompt_confirm_cancel")))
		s.WriteString(m.centerBlock(buf.String()))
	}
}

func (m model) renderAccountManage(s *strings.Builder) {
	var buf strings.Builder
	buf.WriteString(labelStyle.Render(" " + m.tr.Tr("label_accounts")) + "\n\n")

	if len(m.cfg.Profiles) == 0 {
		buf.WriteString("  " + dimStyle.Render(m.tr.Tr("label_no_accounts")) + "\n")
	} else {
		for i, p := range m.cfg.Profiles {
			cursor := "  "
			if i == m.accountCursor {
				cursor = cursorStyle.Render("▸ ")
			}
			delMark := ""
			if m.delMode {
				delMark = " " + errorStyle.Render("[X]")
			}
			idPart := ""
			if p.ID != "" {
				idPart = " (" + p.ID + ")"
			}
			providerName := ""
			if p.ProviderID != "" {
				for _, prov := range m.cfg.Providers {
					if prov.ID == p.ProviderID {
						providerName = prov.Name
						break
					}
				}
			}
			buf.WriteString(fmt.Sprintf("  %s%s%s %s — %s <%s> [%s]\n",
				cursor, keyStyle.Render(p.Name), delMark, idPart, p.GitName, p.GitEmail, providerName))
		}
	}

	s.WriteString(m.centerBlock(strings.TrimRight(buf.String(), "\n")))
}

func (m model) renderAccountForm(s *strings.Builder) {
	title := m.tr.Tr("account_new")
	if m.accountEditIdx >= 0 && m.accountEditIdx < len(m.cfg.Profiles) {
		title = m.tr.Tr("account_edit")
	}

	var buf strings.Builder
	buf.WriteString(labelStyle.Render(" "+title) + "\n\n")

	labels := []string{"ID:", "Name:", "Git Name:", "Git Email:", "Provider:", "SSH Key:"}
	descs := []string{"account_id_desc", "account_name_desc", "account_gitname_desc", "account_gitemail_desc", "account_provider_desc", "account_ssh_identity_desc"}
	for i, input := range m.accountInputs {
		lbl := labels[i]
		if i == m.accountFocused {
			lbl = keyStyle.Render("> " + lbl)
		}
		buf.WriteString(fmt.Sprintf("  %s %s\n", lbl, input.View()))
		buf.WriteString(dimStyle.Render("   " + m.tr.Tr(descs[i])) + "\n")
		if i < len(m.accountInputs)-1 {
			buf.WriteString("\n")
		}
	}

	buf.WriteString("\n\n  " + dimStyle.Render(m.tr.Tr("prompt_tab_next")))
	if m.accountFocused == 5 {
		buf.WriteString("\n  " + dimStyle.Render(m.tr.Tr("ssh_scan_hint")))
	}
	if m.accountFocused == 4 {
		buf.WriteString("\n  " + dimStyle.Render(m.tr.Tr("provider_pick_hint")))
	}

	s.WriteString(m.centerBlock(buf.String()))
}

func (m model) renderSSHKeyPicker(s *strings.Builder) {
	var buf strings.Builder
	buf.WriteString(labelStyle.Render(" " + m.tr.Tr("ssh_scan_title")) + "\n\n")
	for i, key := range m.sshKeys {
		cursor := "  "
		if i == m.sshKeyCursor {
			cursor = cursorStyle.Render("▸ ")
		}
		buf.WriteString(fmt.Sprintf("  %s%s\n", cursor, key))
	}
	buf.WriteString("\n  " + dimStyle.Render(m.tr.Tr("prompt_confirm_cancel")))
	s.WriteString(m.centerBlock(buf.String()))
}

func (m model) renderProviderManage(s *strings.Builder) {
	var buf strings.Builder
	buf.WriteString(labelStyle.Render(" " + m.tr.Tr("label_providers")) + "\n\n")

	providers := m.cfg.Providers
	if len(providers) == 0 {
		buf.WriteString("  " + dimStyle.Render(m.tr.Tr("label_no_providers")) + "\n")
	} else {
		for i, p := range providers {
			cursor := "  "
			if i == m.providerCursor {
				cursor = cursorStyle.Render("▸ ")
			}
			delMark := ""
			if m.providerDelMode {
				delMark = " " + errorStyle.Render("[X]")
			}
			num := fmt.Sprintf("[%d]", i+1)
			buf.WriteString(fmt.Sprintf("  %s%s%s %s (%s) — %s\n",
				cursor, keyStyle.Render(num), delMark,
				p.Name, p.ID, p.Host))
		}
	}

	s.WriteString(m.centerBlock(strings.TrimRight(buf.String(), "\n")))
}

func (m model) renderProviderForm(s *strings.Builder) {
	title := m.tr.Tr("provider_new")
	if m.providerEditIdx >= 0 && m.providerEditIdx < len(m.cfg.Providers) {
		title = m.tr.Tr("provider_edit")
	}

	var buf strings.Builder
	buf.WriteString(labelStyle.Render(" "+title) + "\n\n")

	labels := []string{m.tr.Tr("provider_id_ph") + ":", m.tr.Tr("provider_name_ph") + ":", m.tr.Tr("provider_host_ph") + ":"}
	descs := []string{"provider_id_desc", "provider_name_desc", "provider_host_desc"}
	for i, input := range m.providerInputs {
		lbl := labels[i]
		if i == m.providerFocused {
			lbl = keyStyle.Render("> " + lbl)
		}
		buf.WriteString(fmt.Sprintf("  %s %s\n", lbl, input.View()))
		buf.WriteString(dimStyle.Render("   " + m.tr.Tr(descs[i])) + "\n")
		if i < len(m.providerInputs)-1 {
			buf.WriteString("\n")
		}
	}

	buf.WriteString("\n\n  " + dimStyle.Render(m.tr.Tr("prompt_tab_next")))

	s.WriteString(m.centerBlock(buf.String()))
}

func (m model) renderProviderPicker(s *strings.Builder) {
	var buf strings.Builder
	buf.WriteString(labelStyle.Render(" " + m.tr.Tr("provider_pick_title")) + "\n\n")
	for i, p := range m.providerPickList {
		cursor := "  "
		if i == m.providerPickCursor {
			cursor = cursorStyle.Render("▸ ")
		}
		buf.WriteString(fmt.Sprintf("  %s%s (%s) — %s\n", cursor, p.Name, p.ID, p.Host))
	}
	buf.WriteString("\n  " + dimStyle.Render(m.tr.Tr("prompt_confirm_cancel")))
	s.WriteString(m.centerBlock(buf.String()))
}

func (m model) renderBranchSwitch(s *strings.Builder) {
	var buf strings.Builder
	buf.WriteString(labelStyle.Render(" " + m.tr.Tr("branch_switch_title")) + "\n\n")
	for i, b := range m.branches {
		cursor := "  "
		if i == m.branchCursor {
			cursor = cursorStyle.Render("▸ ")
		}
		mark := ""
		if b == m.branch {
			mark = fmt.Sprintf(" %s", keyStyle.Render(m.tr.Tr("branch_current")))
		}
		buf.WriteString(fmt.Sprintf("  %s%s%s\n", cursor, b, mark))
	}
	buf.WriteString("\n" + dimStyle.Render("  "+m.tr.Tr("prompt_cancel")))
	s.WriteString(m.centerBlock(buf.String()))
}

func (m model) renderBottomBar(s *strings.Builder) {
	var hint string
	switch m.state {
	case stateAccountManage:
		hint = m.buildAccountBar()
	case stateAccountForm:
		hint = ""
	case stateProviderManage:
		hint = m.buildProviderBar()
	case stateProviderForm:
		hint = ""
	case stateProviderPicker:
		hint = m.buildProviderPickerBar()
	case stateSettings:
		hint = m.buildSettingsBar()
	case stateSSHKeyPicker:
		hint = m.buildSSHKeyPickerBar()
	case stateRepoConfig:
		hint = m.buildRepoConfigBar()
	case stateRepoConfigEdit:
		hint = m.buildRepoConfigEditBar()
	case statePROutput:
		hint = fmt.Sprintf("  %s %s", keyStyle.Render("[Esc]"), m.tr.Tr("action_back"))
	default:
		hint = m.buildActionBar()
	}
	s.WriteString(statusBarStyle.Width(m.width).Render(hint))
}

func (m model) buildActionBar() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Y]"), m.tr.Tr("action_copy")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[F5]"), m.tr.Tr("action_refresh")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Q]"), m.tr.Tr("action_quit")))
	return b.String()
}

func (m model) buildAccountBar() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[N]"), m.tr.Tr("action_new")))
	if m.delMode {
		b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Enter]"), m.tr.Tr("action_confirm_del")))
	} else {
		b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[D]"), m.tr.Tr("action_delete")))
		b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Enter]"), m.tr.Tr("action_edit_btn")))
	}
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[↑↓]"), m.tr.Tr("action_nav")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Esc]"), m.tr.Tr("action_back")))
	return b.String()
}

func (m model) buildProviderBar() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[N]"), m.tr.Tr("action_new")))
	if m.providerDelMode {
		b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Enter]"), m.tr.Tr("action_confirm_del")))
	} else {
		b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[D]"), m.tr.Tr("action_delete")))
		b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Enter]"), m.tr.Tr("action_edit_btn")))
	}
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[↑↓]"), m.tr.Tr("action_nav")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Esc]"), m.tr.Tr("action_back")))
	return b.String()
}

func (m model) buildProviderPickerBar() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[↑↓]"), m.tr.Tr("action_nav")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Enter]"), m.tr.Tr("action_select")))
	b.WriteString(fmt.Sprintf("  %s%s", keyStyle.Render("[Esc]"), m.tr.Tr("action_back")))
	return b.String()
}

func (m model) branchDisplay() string {
	if m.branch == "" {
		return m.tr.Tr("status_unable")
	}
	if m.isDirty {
		return m.tr.Tr("status_dirty", m.branch, m.dirtyCount)
	}
	return fmt.Sprintf("%s %s", m.branch, m.tr.Tr("status_clean"))
}

func (m model) routeDisplay() string {
	if m.remoteURL == "" {
		return m.tr.Tr("status_no_remote")
	}
	return m.remoteURL
}

func (m model) identityDisplay() string {
	if m.localName == "" || m.localEmail == "" {
		return m.tr.Tr("status_not_set")
	}
	match := ""
	if p := findMatchingProfile(m.cfg, m.localName, m.localEmail); p != nil {
		match = p.Name
	} else {
		match = m.tr.Tr("status_custom")
	}
	return fmt.Sprintf("%s <%s> [%s]", m.localName, m.localEmail, match)
}

func (m model) actionsView() string {
	var b strings.Builder
	for i, p := range m.cfg.Profiles {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▸ ")
		}
		num := fmt.Sprintf("[%d]", i+1)
		label := p.Name
		if p.ID != "" {
			label = p.Name + " (" + p.ID + ")"
		}
		b.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, keyStyle.Render(num), label))
	}
	np := len(m.cfg.Profiles)
	actionKeys := []string{"C", "T", "R", "G", "A", "P", "S", "B"}
	actionTrs := []string{"action_commit", "tag_hint", "pr_hint", "action_repo_config", "action_accounts", "action_provider_mgmt", "action_settings", "action_branch"}
	for i := 0; i < 8; i++ {
		cursor := "  "
		if np+i == m.cursor {
			cursor = cursorStyle.Render("▸ ")
		}
		b.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, keyStyle.Render("["+actionKeys[i]+"]"), m.tr.Tr(actionTrs[i])))
	}
	return b.String()
}
