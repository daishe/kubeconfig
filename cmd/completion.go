// Copyright 2020 Marek Dalewski
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
	"os"
	"strings"

	"github.com/daishe/kubeconfig/cmd/ui"

	"github.com/spf13/cobra"
)

const completionLong = `To load completion run

  . <(kubeconfig completion <shell>)

Supported shells

  bash
  powershell
  zsh

To configure your bash shell to load completions for each session add to your bashrc

  # ~/.bashrc or ~/.profile
  . <(kubeconfig completion bash)
`

// newCompletionCmd generates a new completion command
func newCompletionCmd(rootCmd *cobra.Command, global *rootOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [shell]",
		Short: "Generates shell completion scripts",
		Long:  completionLong,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			completionRun(rootCmd, global, args[0])
		},
	}

	return cmd
}

func completionRun(rootCmd *cobra.Command, global *rootOpts, shell string) {
	var err error
	switch strings.ToLower(shell) {
	case "bash":
		err = rootCmd.GenBashCompletion(os.Stdout)
	case "powershell":
		err = rootCmd.GenPowerShellCompletion(os.Stdout)
	case "zsh":
		err = rootCmd.GenZshCompletion(os.Stdout)
	default:
		err = fmt.Errorf("unknown shell %q", shell)
	}
	ui.DisplayAndExitOnError(err)
}
