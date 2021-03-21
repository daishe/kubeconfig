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
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/daishe/kubeconfig/registry"
)

const addLong = `Enables addition of a new kubectl config file to kubeconfig registry, by
creating a new empty file with the provided name and opening editor described by
the ${EDITOR} environment variable.`

// newAddCmd generates a new add command.
func newAddCmd(global *rootOpts) *cobra.Command {
	o := &addOpts{}

	cmd := &cobra.Command{
		Use:     "add [name of new config]",
		Short:   "Add a new kubectl config file",
		Long:    addLong,
		Aliases: []string{"new", "create"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.name = args[0]
			return addRun(cmd.Context(), global, o)
		},
	}

	cmd.Flags().StringVarP(&o.editor, "editor", "e", "", "sets the editor used directly, instead of using the ${EDITOR} environment variable")
	cmd.Flags().BoolVarP(&o.force, "force", "f", false, "force override, if the entry with the provided name already exists in the registry")

	return cmd
}

type addOpts struct {
	name   string
	editor string
	force  bool
}

func addRun(ctx context.Context, g *rootOpts, o *addOpts) error {
	editor := o.editor
	if editor == "" {
		e, ok := os.LookupEnv("EDITOR")
		if !ok {
			return fmt.Errorf("$EDITOR environment vatiabble not set")
		}
		editor = e
	}

	reg, err := g.regPath()
	if err != nil {
		return err
	}

	tmp, err := ioutil.TempFile("", "*.config")
	if err != nil {
		return fmt.Errorf("cannot create temporary file: %w", err)
	}
	tmpPath := tmp.Name()
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing file %q filed: %w", tmpPath, err)
	}
	defer os.Remove(tmpPath)

	cmd := exec.CommandContext(ctx, editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting editor %q failed: %w", editor, err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("editor %q failed: %w", editor, err)
	}

	if tmp, err = os.Open(tmpPath); err != nil {
		return fmt.Errorf("cannot open temporary file %q: %w", tmpPath, err)
	}
	defer tmp.Close()
	stat, err := tmp.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat %q: %w", tmpPath, err)
	}

	if stat.Size() == 0 {
		fmt.Println("Skipping an empty config.")
		return nil // just delete temporary
	}

	path := registry.NameToPath(reg, o.name)
	if o.force {
		if err := registry.ForceWrite(path, tmp); err != nil {
			return err
		}
	} else {
		if err := registry.Write(path, tmp); err != nil {
			return err
		}
	}

	fmt.Printf("A new entry %q added to the registry.\n", o.name)
	return nil
}
