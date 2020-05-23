package main

import (
	"os"

	"github.com/hashicorp/terraform-plugin-docs/internal/cmd"
)

func main() {
	name := "tfpluginwebsite"
	version := name + " Version " + version
	if builtBy != "" {
		version += ", built by " + builtBy
	}
	if commit != "" {
		version += " from commit " + commit
	}
	if date != "" {
		version += " on " + date
	}

	os.Exit(cmd.Run(
		name,
		version,
		os.Args[1:],
		os.Stdin, os.Stdout, os.Stderr,
	))
}
