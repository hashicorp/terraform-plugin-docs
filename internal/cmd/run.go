// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/cli"
	"github.com/mattn/go-colorable"

	"github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs/build"
)

type commonCmd struct {
	ui cli.Ui
}

func (cmd *commonCmd) run(r func() error) int {
	err := r()
	if err != nil {
		// TODO: unwraps? check for special exit code error?
		cmd.ui.Error(fmt.Sprintf("Error executing command: %s\n", err))
		os.Exit(1)
	}
	return 0
}

func initCommands(ui cli.Ui) map[string]cli.CommandFactory {

	generateFactory := func() (cli.Command, error) {
		return &generateCmd{
			commonCmd: commonCmd{
				ui: ui,
			},
		}, nil
	}

	defaultFactory := func() (cli.Command, error) {
		return &defaultCmd{
			synopsis: "the generate command is run by default",
			Command: &generateCmd{
				commonCmd: commonCmd{
					ui: ui,
				},
			},
		}, nil
	}

	validateFactory := func() (cli.Command, error) {
		return &validateCmd{
			commonCmd: commonCmd{
				ui: ui,
			},
		}, nil
	}

	migrateFactory := func() (cli.Command, error) {
		return &migrateCmd{
			commonCmd: commonCmd{
				ui: ui,
			},
		}, nil
	}

	return map[string]cli.CommandFactory{
		"":         defaultFactory,
		"generate": generateFactory,
		"validate": validateFactory,
		"migrate":  migrateFactory,
		//"serve": serveFactory,
	}
}

type defaultCmd struct {
	cli.Command
	synopsis string
}

func (cmd *defaultCmd) Synopsis() string {
	return cmd.synopsis
}

func Run(name, version string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var ui cli.Ui = &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,

		Ui: &cli.BasicUi{
			Reader:      stdin,
			Writer:      stdout,
			ErrorWriter: stderr,
		},
	}

	commands := initCommands(ui)

	cli := cli.CLI{
		Name:       name,
		Args:       args,
		Commands:   commands,
		HelpFunc:   cli.BasicHelpFunc(name),
		HelpWriter: stderr,
		Version:    version,
	}

	exitCode, err := cli.Run()
	if err != nil {
		return 1
	}
	return exitCode
}

func Main() int {
	return Run(
		"tfplugindocs",
		build.GetVersion(),
		os.Args[1:],
		os.Stdin,
		colorable.NewColorableStdout(),
		colorable.NewColorableStderr(),
	)
}

// TestScriptMain has the required function signature for use with testscript
func TestScriptMain() {
	Main() // Exit code is no longer used by testscript
}
