package main

import (
	"os"

	"github.com/hashicorp/terraform-plugin-docs/internal/cmd"
	"github.com/mattn/go-colorable"
)

func main() {
	name := "tfplugindocs"
	version := name + " Version " + version
	if commit != "" {
		version += " from commit " + commit
	}

	os.Exit(cmd.Run(
		name,
		version,
		os.Args[1:],
		os.Stdin,
		colorable.NewColorableStdout(),
		colorable.NewColorableStderr(),
	))
}
