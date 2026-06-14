package main

import (
	"fmt"
	"os"
	"strings"
)

func runCLI(args []string, cfg *Config, cfgPath string, tr *Translator) int {
	if len(args) == 0 {
		return 0
	}

	cmd := strings.TrimLeft(args[0], "-")
	switch cmd {
	case "tui":
		return 0

	case "switch":
		return cmdSwitch(args[1:], cfg, tr)

	case "status":
		return cmdStatus(tr)

	case "edit":
		return cmdEdit(cfgPath, tr)

	case "pr":
		return cmdPR(args[1:], tr)

	case "tag":
		return cmdTag(args[1:], tr)

	case "help", "--help", "-h", "-help":
		printHelp(tr)
		return 0

	case "version", "--version", "-v":
		fmt.Println(tr.Tr("cli_version"))
		return 0

	default:
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_error", fmt.Sprintf("unknown command: %s", cmd)))
		printHelp(tr)
		return 1
	}
}

func cmdSwitch(args []string, cfg *Config, tr *Translator) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_error", "usage: gitf switch <id>"))
		return 1
	}

	id := args[0]
	var profile *Profile
	for i := range cfg.Profiles {
		if cfg.Profiles[i].ID == id {
			profile = &cfg.Profiles[i]
			break
		}
	}
	if profile == nil {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_switch_not_found", id))
		return 1
	}

	if err := IsInsideWorkTree(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_no_repo"))
		return 1
	}

	log, err := SwitchProfile(profile, cfg.Providers, tr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_error", err))
		return 1
	}
	fmt.Printf("%s\n", tr.Tr("cli_switch_ok", profile.Name))
	if log != "" {
		fmt.Println(tr.Tr("cli_switch_detail", log))
	}
	return 0
}

func cmdStatus(tr *Translator) int {
	if err := IsInsideWorkTree(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_no_repo"))
		return 1
	}

	fmt.Printf("%s\n", tr.Tr("cli_status_header"))
	fmt.Println(strings.Repeat("-", 30))

	branch := GetCurrentBranch()
	if branch == "" {
		branch = tr.Tr("status_unable")
	}
	fmt.Println(tr.Tr("cli_branch", branch))

	remote := GetRemoteURL()
	if remote == "" {
		remote = tr.Tr("status_no_remote")
	}
	fmt.Println(tr.Tr("cli_remote", remote))

	name := GetLocalUserName()
	email := GetLocalUserEmail()
	if name == "" {
		name = tr.Tr("status_not_set")
	}
	if email == "" {
		email = tr.Tr("status_not_set")
	}
	fmt.Println(tr.Tr("cli_name", name))
	fmt.Println(tr.Tr("cli_email", email))

	_, count := GetWorkTreeStatus()
	if count > 0 {
		fmt.Println(tr.Tr("cli_dirty", count))
	} else {
		fmt.Println(tr.Tr("cli_clean"))
	}

	return 0
}

func cmdEdit(cfgPath string, tr *Translator) int {
	fmt.Println(tr.Tr("cli_edit_opening"))
	cmd := BuildEditCmd(cfgPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_error", err))
		return 1
	}
	return 0
}

func cmdPR(args []string, tr *Translator) int {
	title := strings.Join(args, " ")
	if title == "" {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_error", "usage: gitf pr <title>"))
		return 1
	}

	if IsInsideWorkTree() != nil {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_no_repo"))
		return 1
	}

	output, err := CreatePR(title, tr)
	fmt.Print(output)
	if err != nil {
		return 1
	}
	return 0
}

func cmdTag(args []string, tr *Translator) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_error", "usage: gitf tag <version>"))
		return 1
	}

	version := args[0]
	if IsInsideWorkTree() != nil {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_no_repo"))
		return 1
	}

	if TagExists(version) {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_error", tr.Tr("tag_exists", version)))
		return 1
	}

	output, err := CreateTag(version, tr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", tr.Tr("cli_error", err))
		fmt.Print(output)
		return 1
	}
	fmt.Print(output)
	return 0
}

func printHelp(tr *Translator) {
	fmt.Printf("%s\n\n", tr.Tr("help_title"))
	fmt.Printf("%s\n", tr.Tr("help_usage"))
	fmt.Println(tr.Tr("help_cmd_tui"))
	fmt.Println(tr.Tr("help_cmd_switch"))
	fmt.Println(tr.Tr("help_cmd_status"))
	fmt.Println(tr.Tr("help_cmd_tag"))
	fmt.Println(tr.Tr("help_cmd_pr"))
	fmt.Println(tr.Tr("help_cmd_edit"))
	fmt.Println(tr.Tr("help_cmd_help"))
}
