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
	"github.com/spf13/cobra"

	"github.com/daishe/kubeconfig/cmd/ui"
	"github.com/daishe/kubeconfig/registry"
)

const saveLong = `Saves the current kubectl config file under the provided name in the kubeconfig
registry.`

// newSaveCmd generates a new save command
func newSaveCmd(global *rootOpts) *cobra.Command {
	o := &saveOpts{}

	cmd := &cobra.Command{
		Use:     "save [name in registry]",
		Short:   "Save the current kubectl config file",
		Long:    saveLong,
		Aliases: []string{"sa"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			o.name = args[0]
			saveRun(global, o)
		},
	}

	cmd.Flags().BoolP("force", "f", false, "force override, if the entry with the provided name already exists in the registry")

	return cmd
}

type saveOpts struct {
	name  string
	force bool
}

func saveRun(g *rootOpts, o *saveOpts) {
	reg, err := g.regPath()
	ui.DisplayAndExitOnError(err)

	kcCfgPath, err := registry.KubectlConfigPath()
	ui.DisplayAndExitOnError(err)

	kcCgf, err := registry.Read(kcCfgPath)
	ui.DisplayAndExitOnError(err)
	defer kcCgf.Close()

	if o.force {
		err = registry.ForceWrite(registry.NameToPath(reg, o.name), kcCgf)
		ui.DisplayAndExitOnError(err)
	} else {
		err = registry.Write(registry.NameToPath(reg, o.name), kcCgf)
		ui.DisplayAndExitOnError(err)
	}
}
