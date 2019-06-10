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

package cmd

import (
	"fmt"

	"code.alibaba-inc.com/force/git-repo/config"
	"code.alibaba-inc.com/force/git-repo/versions"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of git-repo",
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckGitAlias()
		showVersion()
	},
}

func showVersion() {
	fmt.Printf("git-repo version %s\n", versions.GetVersion())
	fmt.Printf("git version %s\n", versions.GetGitVersion())
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
