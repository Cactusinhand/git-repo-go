// Copyright © 2019 Alibaba Co. Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"code.alibaba-inc.com/force/git-repo/path"
	"github.com/jiangxin/goconfig"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	// DotRepo is '.repo', admin directory for git-repo
	DotRepo = path.DotRepo

	// CommitIDPattern indicates raw commit ID
	CommitIDPattern = regexp.MustCompile(`^[0-9a-f]{40}([0-9a-f]{24})?$`)
)

// Exported macros
const (
	DefaultConfigPath = ".git-repo"
	DefaultLogRotate  = 20 * 1024 * 1024
	DefaultLogLevel   = "warn"

	CfgRepoArchive           = "repo.archive"
	CfgRepoDepth             = "repo.depth"
	CfgRepoDissociate        = "repo.dissociate"
	CfgRepoMirror            = "repo.mirror"
	CfgRepoReference         = "repo.reference"
	CfgRepoSubmodules        = "repo.submodules"
	CfgManifestGroups        = "manifest.groups"
	CfgManifestName          = "manifest.name"
	CfgRemoteOriginURL       = "remote.origin.url"
	CfgBranchDefaultMerge    = "branch.default.merge"
	CfgManifestRemoteType    = "manifest.remote.%s.type"
	CfgManifestRemoteSSHInfo = "manifest.remote.%s.sshinfo"

	ManifestsDotGit  = "manifests.git"
	Manifests        = "manifests"
	DefaultXML       = "default.xml"
	ManifestXML      = "manifest.xml"
	LocalManifestXML = "local_manifest.xml"
	LocalManifests   = "local_manifests"
	ProjectObjects   = "project-objects"
	Projects         = "projects"

	RefsChanges = "refs/changes/"
	RefsMr      = "refs/merge-requests/"
	RefsHeads   = "refs/heads/"
	RefsTags    = "refs/tags/"
	RefsPub     = "refs/published/"
	RefsM       = "refs/remotes/m/"
	Refs        = "refs/"
	RefsRemotes = "refs/remotes/"

	RemoteTypeGerrit  = "gerrit"
	RemoteTypeAGit    = "agit"
	RemoteTypeUnknown = "unknown"

	MaxJobs = 32

	ViperEnvPrefix = "GIT_REPO"
)

// AssumeNo checks --asume-no option
func AssumeNo() bool {
	return viper.GetBool("assume-no")
}

// AssumeYes checks --asume-yes option
func AssumeYes() bool {
	return viper.GetBool("assume-yes")
}

// GetVerbose gets --verbose option
func GetVerbose() int {
	return viper.GetInt("verbose")
}

// GetQuiet gets --quiet option
func GetQuiet() bool {
	return viper.GetBool("quiet")
}

// IsSingleMode checks --single option
func IsSingleMode() bool {
	return viper.GetBool("single")
}

// GetLogFile gets --logfile option
func GetLogFile() string {
	logfile := viper.GetString("logfile")
	if logfile != "" && !filepath.IsAbs(logfile) {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		logfile = filepath.Join(home, DefaultConfigPath, logfile)
	}
	return logfile
}

// GetLogLevel gets --loglevel option
func GetLogLevel() string {
	return viper.GetString("loglevel")
}

// GetLogRotateSize gets logrotate size from config
func GetLogRotateSize() int64 {
	return viper.GetInt64("logrotate")
}

// NoCertChecks indicates whether ignore ssl cert checks
func NoCertChecks() bool {
	var verify bool

	if viper.GetBool("no-cert-checks") {
		return true
	}

	cfg, err := goconfig.LoadAll("")
	if err != nil {
		return false
	}
	verify = cfg.GetBool("http.sslverify", true)
	return !verify
}

// GetMockSSHInfoStatus gets --mock-ssh-info-status option
func GetMockSSHInfoStatus() int {
	return viper.GetInt("mock-ssh-info-status")
}

// GetMockSSHInfoResponse gets --mock-ssh-info-status option
func GetMockSSHInfoResponse() string {
	return viper.GetString("mock-ssh-info-response")
}

// MockGitPush checks --mock-git-push option
func MockGitPush() bool {
	return viper.GetBool("mock-git-push")
}

// MockNoSymlink checks --mock-no-symlink option
func MockNoSymlink() bool {
	return viper.GetBool("mock-no-symlink")
}

// MockNoTTY checks --mock-no-tty option
func MockNoTTY() bool {
	return viper.GetBool("mock-no-tty")
}

// MockEditScript checks --mock-edit-script option
func MockEditScript() string {
	return viper.GetString("mock-edit-script")
}

// IsDryRun gets --dryrun option
func IsDryRun() bool {
	return viper.GetBool("dryrun")
}

func init() {
	viper.SetDefault("logrotate", DefaultLogRotate)
	viper.SetDefault("loglevel", DefaultLogLevel)

	viper.SetEnvPrefix(ViperEnvPrefix)
	viper.AutomaticEnv()
}
