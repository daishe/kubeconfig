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
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/daishe/kubeconfig/cmd/ui"
	"github.com/daishe/kubeconfig/registry"
)

const showLong = `Displays the content of the requested file from the kubeconfig registry.

If the kubectl config file is not specified, the command presents an interactive
list of all files in the registry with an option to select one.`

// newShowCmd generates a new show command
func newShowCmd(global *rootOpts) *cobra.Command {
	o := &showOpts{}

	cmd := &cobra.Command{
		Use:     "show [config name]",
		Short:   "Show the content of the requested file from the registry",
		Long:    showLong,
		Aliases: []string{"sho", "display", "disp", "d"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			o.interactive = true
			if len(args) > 0 {
				o.name = args[0]
				o.interactive = false
			}
			showRun(global, o)
		},
	}

	return cmd
}

type showOpts struct {
	name        string
	interactive bool
}

func showRun(g *rootOpts, o *showOpts) {
	reg, err := g.regPath()
	ui.DisplayAndExitOnError(err)

	path := ""
	if o.interactive {
		kcCfgPath, err := registry.KubectlConfigPath()
		ui.DisplayAndExitOnError(err)

		kcCfgHash, err := registry.Hash(kcCfgPath)
		ui.DisplayAndExitOnError(err)

		path, err = ui.SelectPrompt(reg, "Which kubectl config file to show", kcCfgHash)
		ui.DisplayAndExitOnError(err)
	} else {
		path = registry.NameToPath(reg, o.name)
	}

	toShow, err := registry.Read(path)
	ui.DisplayAndExitOnError(err)
	defer toShow.Close()

	_, err = io.Copy(os.Stdout, toShow)
	ui.DisplayAndExitOnError(err)
}
