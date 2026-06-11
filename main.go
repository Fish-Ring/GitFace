package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfgPath := ConfigPath()
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		cfg = DefaultConfig()
	}

	detected := DetectLang()
	lang := detected
	if cfg.Lang != "" {
		lang = Lang(cfg.Lang)
	}
	tr := NewTranslator(lang)

	if len(os.Args) > 1 {
		exitCode := runCLI(os.Args[1:], cfg, cfgPath, tr)
		if exitCode != 0 || os.Args[1] != "tui" {
			os.Exit(exitCode)
		}
	}

	m := NewModel(cfg, cfgPath, tr)
	if err != nil {
		m.errMsg = fmt.Sprintf("Config error: %v", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
