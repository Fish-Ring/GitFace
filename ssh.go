package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var defaultKeyNames = []string{
	"id_rsa",
	"id_ed25519",
	"id_ecdsa",
	"id_ecdsa_sk",
	"id_ed25519_sk",
	"id_dsa",
	"id_xmss",
}

var sshExcludedFiles = map[string]bool{
	"config":          true,
	"known_hosts":     true,
	"authorized_keys": true,
	"authorized_keys2": true,
	"known_hosts2":    true,
	"environment":     true,
	"rc":              true,
}

func scanSSHKeys() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	sshDir := filepath.Join(home, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return nil
	}

	var known, others []string
	knownSet := make(map[string]bool)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".pub") {
			continue
		}
		if sshExcludedFiles[name] {
			continue
		}

		fullPath := filepath.Join(sshDir, name)
		if _, err := os.Stat(fullPath); err != nil {
			continue
		}

		isKnown := false
		for _, kn := range defaultKeyNames {
			if name == kn {
				isKnown = true
				break
			}
		}

		if isKnown && !knownSet[fullPath] {
			knownSet[fullPath] = true
			known = append(known, fullPath)
		} else if !isKnown {
			others = append(others, fullPath)
		}
	}

	sort.Strings(known)
	sort.Strings(others)
	return append(known, others...)
}
