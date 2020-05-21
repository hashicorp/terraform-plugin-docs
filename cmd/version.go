package cmd

import "fmt"

var (
	GitCommit string

	Version string
)

type versionCommand struct {
	baseCommand
}

func (c *versionCommand) Run(_ []string) int {
	c.UI.Output(fmt.Sprintf("%s %s (%s)", Name, Version, GitCommit))
	return 0
}

func (c *versionCommand) Help() string {
	return ""
}

func (c *versionCommand) Synopsis() string {
	return "Prints the version"
}
