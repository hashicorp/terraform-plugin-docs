package cmd

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/hashicorp/terraform-plugin-docs/internal/provider"
	"golang.org/x/tools/go/buildutil"
)

type generateCmd struct {
	commonCmd

	flagLegacySidebar bool
	buildTags         string
}

func (cmd *generateCmd) Synopsis() string {
	return "generates a plugin website from code, templates, and examples for the current directory"
}

func (cmd *generateCmd) Help() string {
	buf := bytes.Buffer{}
	flags := cmd.Flags()
	flags.SetOutput(&buf)
	// PrintDefaults implicitly prints to output
	// thus our buffer
	flags.PrintDefaults()
	return fmt.Sprintf(`Usage: tfplugindocs generate [options]
Available options:
%s
`, buf.String())
}

func (cmd *generateCmd) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.BoolVar(&cmd.flagLegacySidebar, "legacy-sidebar", false, "generate the legacy .erb sidebar file")
	fs.StringVar(&cmd.buildTags, "tags", "", buildutil.TagsFlagDoc)
	return fs
}

func (cmd *generateCmd) Run(args []string) int {
	fs := cmd.Flags()
	err := fs.Parse(args)
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("unable to parse flags: %s", err))
		return 1
	}

	return cmd.run(cmd.runInternal)
}

func (cmd *generateCmd) runInternal() error {
	err := provider.Generate(cmd.ui, cmd.flagLegacySidebar, cmd.buildTags)
	if err != nil {
		return fmt.Errorf("unable to generate website: %w", err)
	}

	return nil
}
