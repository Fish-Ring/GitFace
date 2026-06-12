package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	if err != nil {
		return output, fmt.Errorf("%s", output)
	}
	return output, nil
}

func IsInsideWorkTree() error {
	_, err := runGit("rev-parse", "--is-inside-work-tree")
	return err
}

func GetCurrentBranch() string {
	out, _ := runGit("branch", "--show-current")
	return out
}

func GetRemoteURL() string {
	out, _ := runGit("config", "--get", "remote.origin.url")
	return out
}

func GetLocalUserName() string {
	out, _ := runGit("config", "--local", "--get", "user.name")
	return out
}

func GetLocalUserEmail() string {
	out, _ := runGit("config", "--local", "--get", "user.email")
	return out
}

func GetWorkTreeStatus() (dirty bool, count int) {
	out, err := runGit("status", "--porcelain")
	if err != nil || out == "" {
		return false, 0
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	return true, len(lines)
}

var sshPattern = regexp.MustCompile(`^git@([^:]+):(.*)$`)
var httpsPattern = regexp.MustCompile(`^https://([^/]+)/([^/]+)/(.*)$`)

func rewriteRemoteHost(currentURL, host string) string {
	if m := sshPattern.FindStringSubmatch(currentURL); len(m) == 3 {
		return fmt.Sprintf("git@%s:%s", host, m[2])
	}
	if m := httpsPattern.FindStringSubmatch(currentURL); len(m) == 4 {
		return fmt.Sprintf("git@%s:%s/%s", host, m[2], m[3])
	}
	return ""
}

func detectProviderFromURL(currentURL string, providers []Provider) *Provider {
	var host string
	if m := sshPattern.FindStringSubmatch(currentURL); len(m) == 3 {
		host = m[1]
	} else if m := httpsPattern.FindStringSubmatch(currentURL); len(m) == 4 {
		host = m[1]
	}
	if host == "" {
		return nil
	}
	for i := range providers {
		if providers[i].Host == host {
			return &providers[i]
		}
	}
	return nil
}

func extractPathFromURL(currentURL string) string {
	var repoPath string
	if m := sshPattern.FindStringSubmatch(currentURL); len(m) == 3 {
		repoPath = m[2]
	} else if m := httpsPattern.FindStringSubmatch(currentURL); len(m) == 4 {
		repoPath = m[2] + "/" + m[3]
	}
	return strings.TrimSuffix(repoPath, ".git")
}

func SwitchProfile(p *Profile, providers []Provider, tr *Translator) (string, error) {
	var logs []string

	if _, err := runGit("config", "--local", "user.name", p.GitName); err != nil {
		return "", errors.New(tr.Tr("sp_name_fail", err))
	}
	logs = append(logs, fmt.Sprintf("user.name -> %s", p.GitName))

	if _, err := runGit("config", "--local", "user.email", p.GitEmail); err != nil {
		return "", errors.New(tr.Tr("sp_email_fail", err))
	}
	logs = append(logs, fmt.Sprintf("user.email -> %s", p.GitEmail))

	if p.ProviderID != "" {
		var provider *Provider
		for i := range providers {
			if providers[i].ID == p.ProviderID {
				provider = &providers[i]
				break
			}
		}
		if provider != nil {
			repoCfg := LoadRepoConfig()
			remotesOut, _ := runGit("remote")
			remotes := strings.Fields(remotesOut)
			if len(remotes) == 0 {
				logs = append(logs, tr.Tr("sp_no_url", provider.Name))
			} else {
				for _, remote := range remotes {
					currentURL, _ := runGit("config", "--get", "remote."+remote+".url")
					if currentURL == "" {
						continue
					}
					if srcProvider := detectProviderFromURL(currentURL, providers); srcProvider != nil {
						if srcPath := extractPathFromURL(currentURL); srcPath != "" {
							repoCfg.Paths[srcProvider.ID] = srcPath
						}
					}
				}
				changed := false
				for _, remote := range remotes {
					currentURL, _ := runGit("config", "--get", "remote."+remote+".url")
					if currentURL == "" {
						continue
					}
					newURL := rewriteRemoteHost(currentURL, provider.Host)
					if newURL == "" {
						continue
					}
					hostChanged := newURL != currentURL
					if hostChanged {
						if repoPath, ok := repoCfg.Paths[provider.ID]; ok {
							path := strings.TrimSuffix(repoPath, ".git")
							newURL = fmt.Sprintf("git@%s:%s.git", provider.Host, path)
						}
					}
					if newURL == currentURL {
						continue
					}
					if _, err := runGit("remote", "set-url", remote, newURL); err != nil {
						return "", errors.New(tr.Tr("sp_remote_fail", err))
					}
					logs = append(logs, fmt.Sprintf("Remote (%s): %s -> %s (%s)", remote, currentURL, newURL, provider.Name))
					changed = true
				}
				SaveRepoConfig(repoCfg)
				if !changed {
					logs = append(logs, tr.Tr("sp_remote_same"))
				}
			}
		} else {
			logs = append(logs, fmt.Sprintf(tr.Tr("sp_prov_not_found"), p.ProviderID))
		}
	} else {
		logs = append(logs, tr.Tr("sp_no_provider"))
	}

	if p.SSHIdentityFile != "" {
		sshPath := strings.ReplaceAll(p.SSHIdentityFile, "\\", "/")
		if _, err := runGit("config", "core.sshCommand", fmt.Sprintf("ssh -i \"%s\"", sshPath)); err != nil {
			return "", errors.New(tr.Tr("sp_ssh_fail", err))
		}
		logs = append(logs, fmt.Sprintf("SSH key -> %s", p.SSHIdentityFile))
	} else {
		runGit("config", "--unset", "core.sshCommand")
	}

	return strings.Join(logs, "\n"), nil
}

func findMatchingProfile(cfg *Config, localName, localEmail string) *Profile {
	if localName == "" || localEmail == "" {
		return nil
	}
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		if p.GitName == localName && p.GitEmail == localEmail {
			return p
		}
	}
	return nil
}

func TagExists(version string) bool {
	_, err := runGit("rev-parse", "refs/tags/"+version)
	return err == nil
}

func CreateTag(version string, tr *Translator) (string, error) {
	var output strings.Builder

	writeLine := func(format string, a ...any) {
		output.WriteString(fmt.Sprintf(format, a...))
	}

	writeLine("> git tag %s\n", version)
	if _, err := runGit("tag", version); err != nil {
		return output.String(), errors.New(tr.Tr("tag_git_tag_fail", err))
	}

	writeLine("> git push origin %s\n", version)
	out, err := runGit("push", "origin", version)
	output.WriteString(out)
	if err != nil {
		return output.String(), errors.New(tr.Tr("tag_git_push_fail", strings.TrimSpace(out)))
	}

	writeLine("\n%s\n", tr.Tr("tag_success", version))
	return output.String(), nil
}

func getGitIndexPath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	gitDir := strings.TrimSpace(string(out))
	return filepath.Join(gitDir, "index"), nil
}

func RunConventionalCommit(commitType, desc string, tr *Translator) (string, error) {
	var output strings.Builder
	msg := fmt.Sprintf("%s: %s", commitType, desc)

	writeLine := func(format string, a ...any) {
		output.WriteString(fmt.Sprintf(format, a...))
	}

	// Remove stale index.lock if present
	if lockPath, err := getGitIndexPath(); err == nil {
		lockFile := lockPath + ".lock"
		if _, statErr := os.Stat(lockFile); statErr == nil {
			os.Remove(lockFile)
			writeLine("Removed stale %s\n", filepath.Base(lockFile))
		}
	}

	writeLine("> git add .\n")
	cmd := exec.Command("git", "add", ".")
	out, err := cmd.CombinedOutput()
	output.Write(out)
	if err != nil {
		return output.String(), errors.New(tr.Tr("commit_git_add_fail", strings.TrimSpace(string(out))))
	}

	writeLine("> git commit -m \"%s\"\n", msg)
	cmd = exec.Command("git", "commit", "-m", msg)
	out, err = cmd.CombinedOutput()
	output.Write(out)
	hasNewCommit := err == nil
	if err != nil {
		commitOut := strings.TrimSpace(string(out))
		if strings.Contains(commitOut, "nothing to commit") {
			writeLine("\n%s\n", tr.Tr("commit_nothing_new"))
		} else {
			return output.String(), errors.New(tr.Tr("commit_git_commit_fail", commitOut))
		}
	}

	writeLine("> git push --set-upstream origin HEAD\n")
	cmd = exec.Command("git", "push", "--set-upstream", "origin", "HEAD")
	out, err = cmd.CombinedOutput()
	output.Write(out)
	if err != nil {
		return output.String(), errors.New(tr.Tr("commit_git_push_fail", strings.TrimSpace(string(out))))
	}

	if hasNewCommit {
		writeLine("\n%s\n", tr.Tr("commit_push_done"))
	} else {
		writeLine("\n%s\n", tr.Tr("commit_pushed_only"))
	}
	return output.String(), nil
}
