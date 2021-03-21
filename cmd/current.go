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
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/daishe/kubeconfig/cmd/ui"
	"github.com/daishe/kubeconfig/registry"
)

// newCurrentCmd generates a new current command
func newCurrentCmd(global *rootOpts) *cobra.Command {
	o := &currentOpts{}

	cmd := &cobra.Command{
		Use:     "current",
		Short:   "Show the current kubectl config file",
		Long:    `Displays under what name the current kubectl config file is known to kubeconfig.`,
		Aliases: []string{"curr", "cur", "c"},
		Run: func(cmd *cobra.Command, args []string) {
			currentRun(global, o)
		},
	}

	cmd.Flags().BoolVarP(&o.dumpConfig, "dump", "d", false, "dump the content of the kubectl config file instead of reporting kubeconfig name")

	return cmd
}

type currentOpts struct {
	dumpConfig bool
}

func currentRun(g *rootOpts, o *currentOpts) {
	if o.dumpConfig {
		kcCfgPath, err := registry.KubectlConfigPath()
		ui.DisplayAndExitOnError(err)

		toShow, err := registry.Read(kcCfgPath)
		ui.DisplayAndExitOnError(err)
		defer toShow.Close()

		_, err = io.Copy(os.Stdout, toShow)
		ui.DisplayAndExitOnError(err)

		return
	}

	kcCfgPath, err := registry.KubectlConfigPath()
	ui.DisplayAndExitOnError(err)

	reg, err := g.regPath()
	ui.DisplayAndExitOnError(err)

	kcCfgHash, err := registry.Hash(kcCfgPath)
	ui.DisplayAndExitOnError(err)

	currentPath, found, err := registry.Find(reg, kcCfgHash)
	ui.DisplayAndExitOnError(err)

	if found {
		fmt.Println(registry.PathToName(reg, currentPath))
	} else {
		fmt.Println("Current kubectl config file is not in the registry.")
	}
}
