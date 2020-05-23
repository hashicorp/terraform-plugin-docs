package cmd

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-docs/internal/provider"
)

type generateCmd struct {
	commonCmd
}

func (cmd *generateCmd) Run(args []string) int {
	// TODO: flags?
	return cmd.run(cmd.runInternal)
}

func (cmd *generateCmd) runInternal() error {
	err := provider.Generate(func(format string, a ...interface{}) {
		cmd.ui.Info(fmt.Sprintf(format, a...))
	})
	if err != nil {
		return fmt.Errorf("unable to generate website: %w", err)
	}

	return nil
}
