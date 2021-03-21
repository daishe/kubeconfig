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

const switchLong = `Switches the current kubectl config file to the requested one, by overriding the
'${HOME}/.kube/config'.

If the current kubectl config file is not known to kubeconfig and overriding it
will loose some inforatmion, the command will fail.

If the kubectl config file is not specified, the command presents an interactive
list of all files in the registry with an option to select one.`

// newSwitchCmd generates a new switch command
func newSwitchCmd(global *rootOpts) *cobra.Command {
	o := &switchOpts{}

	cmd := &cobra.Command{
		Use:     "switch [config name]",
		Short:   "Switch between available kubectl config files",
		Long:    switchLong,
		Aliases: []string{"swch", "sw", "s", "select", "sel"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			o.interactive = true
			if len(args) > 0 {
				o.name = args[0]
				o.interactive = false
			}
			switchRun(global, o)
		},
	}

	cmd.Flags().BoolVarP(&o.force, "force", "f", false, "force switching, ignore overriding not known kubectl config file")

	return cmd
}

type switchOpts struct {
	name        string
	force       bool
	interactive bool
}

func switchRun(g *rootOpts, o *switchOpts) {
	reg, err := g.regPath()
	ui.DisplayAndExitOnError(err)

	kcCfgPath, err := registry.KubectlConfigPath()
	ui.DisplayAndExitOnError(err)

	kcCfgHash, err := registry.Hash(kcCfgPath)
	ui.DisplayAndExitOnError(err)

	current, found, err := registry.Find(reg, kcCfgHash)
	ui.DisplayAndExitOnError(err)

	if !found && !o.force {
		ui.DisplayAndExitOnError(fmt.Errorf("the current kubectl config do not exists in the registry and switching config files will override it; If that is intended force with '--force' flag"))
	}

	path := ""
	if o.interactive {
		path, err = ui.SelectPrompt(reg, "Which kubectl config file to switch to", kcCfgHash)
		ui.DisplayAndExitOnError(err)
	} else {
		path = registry.NameToPath(reg, o.name)
	}

	cfg, err := registry.Read(path)
	ui.DisplayAndExitOnError(err)
	defer cfg.Close()

	err = registry.ForceWrite(kcCfgPath, cfg)
	ui.DisplayAndExitOnError(err)

	if !found {
		fmt.Printf("Successfully switched from unknown kubectl config file to %q.\n", registry.PathToName(reg, path))
	} else {
		fmt.Printf("Successfully switched from %q to %q.\n", registry.PathToName(reg, current), registry.PathToName(reg, path))
	}
}
