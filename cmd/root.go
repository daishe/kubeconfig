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

	"github.com/spf13/cobra"

	"github.com/daishe/kubeconfig/cmd/ui"
	"github.com/daishe/kubeconfig/registry"
)

const rootLong = `A utility tool to manage, swap currently used, store, etc different kubectl
config files.`

// newRootCmd generates a new base command with all its subcommands.
func newRootCmd() *cobra.Command {
	o := &rootOpts{}

	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "kubeconfig - a utility to swap kubectl config files",
		Long:  rootLong,
	}

	cmd.PersistentFlags().StringVar(&o.altRegistryPath, "registry", "", "override the default registry path")

	cmd.AddCommand(newAddCmd(o))
	cmd.AddCommand(newCompletionCmd(cmd, o))
	cmd.AddCommand(newCurrentCmd(o))
	cmd.AddCommand(newEditCmd(o))
	cmd.AddCommand(newListCmd(o))
	cmd.AddCommand(newSaveCmd(o))
	cmd.AddCommand(newShowCmd(o))
	cmd.AddCommand(newSwitchCmd(o))

	return cmd
}

type rootOpts struct {
	altRegistryPath string
}

func (o *rootOpts) regPath() (string, error) {
	if o.altRegistryPath != "" {
		return registry.OverrodePath(o.altRegistryPath)
	}
	return registry.Path()
}

// Execute runs the application. It uses the os.Args[1:] and runs through the commands tree finding appropriate matches for commands and then corresponding flags.
func Execute(ctx context.Context) {
	err := newRootCmd().ExecuteContext(ctx)
	ui.DisplayAndExitOnError(err)
}

// ExecuteWith runs the application. It uses the provided args and runs through the commands tree finding appropriate matches for commands and then corresponding flags.
func ExecuteWith(ctx context.Context, args []string) {
	cmd := newRootCmd()
	cmd.SetArgs(args)
	err := cmd.ExecuteContext(ctx)
	ui.DisplayAndExitOnError(err)
}
