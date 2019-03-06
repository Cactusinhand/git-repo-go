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
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Define macros for config
const (
	DefaultConfigPath = ".git-repo"
	DefaultLogRotate  = 20 * 1024 * 1024

	RepoDefaultManifestKey = "repo.manifestDefault"
	RepoDefaultManifestXML = "default.xml"
)

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
	if logfile != "" && !path.IsAbs(logfile) {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		logfile = path.Join(home, DefaultConfigPath, logfile)
	}
	return logfile
}

// GetLogLevel gets --loglevel option
func GetLogLevel() string {
	return viper.GetString("loglevel")
}

// GetLogRotate gets logrotate size from config
func GetLogRotate() int64 {
	viper.SetDefault("logrotate", DefaultLogRotate)
	return viper.GetInt64("logrotate")
}
