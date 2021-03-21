package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"

	"github.com/daishe/kubeconfig/registry"
)

// DisplayAndExitOnError (if the provided error is not nil) prints the given error to stderr and exits (with exit code 1).
func DisplayAndExitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// ListToNames converts the provided list of files in registry and the given comparison results to list of names.
func ListToNames(reg string, files []string) []string {
	names := make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, registry.PathToName(reg, f))
	}
	return names
}

// AnnotateNamesWithCurrent annotates list of names with information whether the given entry is the current one.
func AnnotateNamesWithCurrent(names []string, cmp []bool) []string {
	anno := make([]string, 0, len(names))
	for i := 0; i < len(names); i++ {
		if cmp[i] {
			anno = append(anno, names[i]+"   <---- current -----")
		} else {
			anno = append(anno, names[i])
		}
	}
	return anno
}

// SelectPrompt will display the select prompt with the list of entries in the registry.
func SelectPrompt(reg string, msg string, hash []byte) (string, error) {
	ls, cmp, err := registry.ListWithCmp(reg, hash)
	if err != nil {
		return "", err
	}
	names := ListToNames(reg, ls)

	component := promptui.Select{
		Label:             msg,
		Items:             AnnotateNamesWithCurrent(names, cmp),
		Size:              20,
		HideHelp:          true,
		StartInSearchMode: true,
		Searcher: func(input string, i int) bool {
			return strings.Contains(names[i], input)
		},
		Keys: &promptui.SelectKeys{
			Prev:     promptui.Key{Code: readline.CharPrev, Display: "↑"},
			Next:     promptui.Key{Code: readline.CharNext, Display: "↓"},
			PageUp:   promptui.Key{Code: readline.CharBackward, Display: "←"},
			PageDown: promptui.Key{Code: readline.CharForward, Display: "→"},
			Search:   promptui.Key{Code: readline.CharCtrlW, Display: "^W"},
		},
	}

	idx, _, err := component.Run()
	if err != nil {
		return "", err
	}

	return ls[idx], nil
}
