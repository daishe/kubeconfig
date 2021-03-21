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
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/daishe/kubeconfig/cmd/ui"
	"github.com/daishe/kubeconfig/registry"
)

const editLong = `Enables editition of a new kubectl config file to kubeconfig registry, by
creating a new empty file with the provided name and opening editor described by
the ${EDITOR} environment variable.

If the kubectl config file is not specified, the command presents an interactive
list of all files in the registry with an option to select one.`

// newEditCmd generates a new edit command.
func newEditCmd(global *rootOpts) *cobra.Command {
	o := &editOpts{}

	cmd := &cobra.Command{
		Use:     "edit [name of new config]",
		Short:   "Edit a new kubectl config file",
		Long:    editLong,
		Aliases: []string{"new", "create"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			o.interactive = true
			if len(args) > 0 {
				o.name = args[0]
				o.interactive = false
			}
			editRun(cmd.Context(), global, o)
		},
	}

	cmd.Flags().StringVarP(&o.editor, "editor", "e", "", "sets the editor used directly, instead of using the ${EDITOR} environment variable")

	return cmd
}

type editOpts struct {
	name        string
	editor      string
	interactive bool
}

func editRun(ctx context.Context, g *rootOpts, o *editOpts) {
	editor := o.editor
	if editor == "" {
		e, ok := os.LookupEnv("EDITOR")
		if !ok {
			ui.DisplayAndExitOnError(fmt.Errorf("$EDITOR environment vatiabble not set"))
		}
		editor = e
	}

	reg, err := g.regPath()
	ui.DisplayAndExitOnError(err)

	path := ""
	if o.interactive {
		kcCfgPath, err := registry.KubectlConfigPath()
		ui.DisplayAndExitOnError(err)

		kcCfgHash, err := registry.Hash(kcCfgPath)
		ui.DisplayAndExitOnError(err)

		path, err = ui.SelectPrompt(reg, "Which kubectl config file to edit", kcCfgHash)
		ui.DisplayAndExitOnError(err)
	} else {
		path = registry.NameToPath(reg, o.name)
	}

	cmd := exec.CommandContext(ctx, editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		ui.DisplayAndExitOnError(fmt.Errorf("starting editor %q failed: %w", editor, err))
	}
	if err := cmd.Wait(); err != nil {
		ui.DisplayAndExitOnError(fmt.Errorf("editor %q failed: %w", editor, err))
	}
}
