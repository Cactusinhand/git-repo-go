package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"code.alibaba-inc.com/force/git-repo/path"
	"github.com/jiangxin/goconfig"
	log "github.com/jiangxin/multi-log"
)

// Macros for git-extra-config
const (
	GitExtraConfigVersion = "3"
	GitExtraConfigFile    = "~/.git-repo/gitconfig"
	CfgRepoConfigVersion  = "repo.configversion"
)

var (
	gitConfigExtension = `
# This file is generated by git-repo.
# DO NOT edit this file! Any modification will be overwritten.
#
# Command alias
[alias]
	br = branch
	ci = commit -s
	co = checkout
	cp = cherry-pick
	logf = log --pretty=fuller
	logs = log --pretty=refs  --date=short
	pr = repo upload --single
	peer-review = repo upload --single
	review = repo upload --single
	st = status
[color]
	ui = auto
[core]
	# Do not quote path, show UTF-8 characters directly
	quotepath = false
[merge]
	# Add at most 20 commit logs in merge log message
	log = true
[pretty]
	refs = format:%h (%s, %ad)
[rebase]
	# Run git rebase with --autosquash option
	autosquash = true
[repo]
	# Version of this git config extension
	configversion = ` + GitExtraConfigVersion + `
`
)

// CheckGitAlias checks if any alias command has been overridden.
func CheckGitAlias() {
	var aliasCommands = []string{
		"git-review",
		"git-pr",
		"git-peer-review",
	}

	for _, cmd := range aliasCommands {
		p, err := exec.LookPath(cmd)
		if err == nil {
			log.Warnf("you cannot use the git-repo alias command '%s', it is overrided by '%s' installed", cmd, p)
		}
	}
}

func saveExtraGitConfig() error {
	var (
		err error
	)

	filename, _ := path.Abs(GitExtraConfigFile)
	dir := filepath.Dir(filename)

	if _, err := os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(gitConfigExtension)
	return err
}

// InstallExtraGitConfig if necessary
func InstallExtraGitConfig() error {
	var err error

	globalConfig, err := goconfig.GlobalConfig()
	version := globalConfig.Get(CfgRepoConfigVersion)
	if version == GitExtraConfigVersion {
		return nil
	}

	log.Debugf("unmatched git config version: %s != %s", version, GitExtraConfigVersion)
	found := false
	gitExtraConfigFile, _ := path.Abs(GitExtraConfigFile)
	for _, p := range globalConfig.GetAll("include.path") {
		p, _ = path.Abs(p)
		if p == gitExtraConfigFile {
			found = true
			break
		}
	}
	if !found {
		cmds := []string{"git",
			"config",
			"--global",
			"--add",
			"include.path",
			GitExtraConfigFile,
		}
		err = exec.Command(cmds[0], cmds[1:]...).Run()
		if err != nil {
			return err
		}
	}

	err = saveExtraGitConfig()
	if err != nil {
		return err
	}
	return nil
}
