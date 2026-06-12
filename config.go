package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Provider struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
}

type Profile struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	GitName         string            `json:"git_name"`
	GitEmail        string            `json:"git_email"`
	ProviderID      string            `json:"provider_id,omitempty"`
	SSHIdentityFile string            `json:"ssh_identity_file"`
	RemotePaths     map[string]string `json:"remote_paths,omitempty"`
}

type Config struct {
	Lang      string     `json:"lang"`
	Profiles  []Profile  `json:"profiles"`
	Providers []Provider `json:"providers"`
}

func DefaultConfig() *Config {
	detectedLang := DetectLang()
	return &Config{
		Lang:      string(detectedLang),
		Profiles:  []Profile{},
		Providers: defaultProviders(),
	}
}

func defaultProviders() []Provider {
	return []Provider{
		{ID: "github", Name: "GitHub", Host: "github.com"},
		{ID: "gitee", Name: "Gitee", Host: "gitee.com"},
		{ID: "gitlab", Name: "GitLab", Host: "gitlab.com"},
		{ID: "bitbucket", Name: "Bitbucket", Host: "bitbucket.org"},
	}
}

func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".config", "gitface", "config.json")
	}
	return filepath.Join(home, ".config", "gitface", "config.json")
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("create config dir failed: %w", err)
			}
			cfg := DefaultConfig()
			if err := SaveConfig(cfg, path); err != nil {
				return nil, fmt.Errorf("write default config failed: %w", err)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}
	if cfg.Providers == nil {
		cfg.Providers = defaultProviders()
	}
	if cfg.Lang == "" {
		cfg.Lang = string(DetectLang())
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize config failed: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func BuildEditCmd(path string) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vim"
		}
	}
	editor = strings.TrimSpace(editor)

	if info, err := os.Stat(editor); err == nil && !info.IsDir() {
		return exec.Command(editor, path)
	}

	parts := strings.Fields(editor)
	if len(parts) == 0 {
		if runtime.GOOS == "windows" {
			parts = []string{"notepad"}
		} else {
			parts = []string{"vim"}
		}
	}
	args := append(parts[1:], path)
	return exec.Command(parts[0], args...)
}
