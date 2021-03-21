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

	"github.com/spf13/cobra"

	"github.com/daishe/kubeconfig/cmd/ui"
	"github.com/daishe/kubeconfig/registry"
)

// newListCmd generates a new list command
func newListCmd(global *rootOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Show all kubectl config files",
		Long:    `Shows names of all kubectl config files in the kubeconfig registry.`,
		Aliases: []string{"lst", "ls", "l", "li"},
		Run: func(cmd *cobra.Command, args []string) {
			listRun(global)
		},
	}

	return cmd
}

func listRun(g *rootOpts) {
	reg, err := g.regPath()
	ui.DisplayAndExitOnError(err)

	kcCfgPath, err := registry.KubectlConfigPath()
	ui.DisplayAndExitOnError(err)

	kcCfgHash, err := registry.Hash(kcCfgPath)
	ui.DisplayAndExitOnError(err)

	ls, cmp, err := registry.ListWithCmp(reg, kcCfgHash)
	ui.DisplayAndExitOnError(err)

	for _, name := range ui.AnnotateNamesWithCurrent(ui.ListToNames(reg, ls), cmp) {
		fmt.Println(name)
	}
}
