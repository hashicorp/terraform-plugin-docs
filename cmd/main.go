package cmd

import (
	"os"

	"github.com/mitchellh/cli"
)

const (
	Name = "tfproviderdocsgen"
)

func Run(args []string) int {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	commands := initCommands(ui)

	// override version flags to call version command
	if len(args) == 1 && (args[0] == "-v" || args[0] == "-version" || args[0] == "--version") {
		args = []string{"version"}
	}

	cli := cli.CLI{
		Name:       Name,
		Args:       args,
		Commands:   commands,
		HelpFunc:   cli.BasicHelpFunc(Name),
		HelpWriter: os.Stdout,
		// Version:    fmt.Sprintf("%s, build %s", Version, Build),
	}

	exitCode, err := cli.Run()
	if err != nil {
		return 1
	}
	return exitCode
}

func initCommands(ui cli.Ui) map[string]cli.CommandFactory {
	baseCommand := baseCommand{
		UI: ui,
	}

	return map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return &runCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &versionCommand{
				baseCommand: baseCommand,
			}, nil
		},
	}
}
