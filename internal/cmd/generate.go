package cmd

import (
	"flag"
	"fmt"

	"github.com/hashicorp/terraform-plugin-docs/internal/provider"
)

type generateCmd struct {
	commonCmd

	flagLegacySidebar bool

	flagProviderName         string
	flagRenderedProviderName string

	flagRenderedWebsiteDir string
	flagExamplesDir        string
	flagWebsiteTmpDir      string
	flagWebsiteSourceDir   string
	tfVersion              string
}

func (cmd *generateCmd) Synopsis() string {
	return "generates a plugin website from code, templates, and examples for the current directory"
}

func (cmd *generateCmd) Help() string {
	return `Usage: tfplugindocs generate`
}

func (cmd *generateCmd) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.BoolVar(&cmd.flagLegacySidebar, "legacy-sidebar", false, "generate the legacy .erb sidebar file")
	fs.StringVar(&cmd.flagProviderName, "provider-name", "", "provider name, as used in Terraform configurations")
	fs.StringVar(&cmd.flagRenderedProviderName, "rendered-provider-name", "", "provider name, as generated in documentation (ex. page titles, ...)")
	fs.StringVar(&cmd.flagRenderedWebsiteDir, "rendered-website-dir", "docs", "output directory")
	fs.StringVar(&cmd.flagExamplesDir, "examples-dir", "examples", "examples directory")
	fs.StringVar(&cmd.flagWebsiteTmpDir, "website-temp-dir", "", "temporary directory (used during generation)")
	fs.StringVar(&cmd.flagWebsiteSourceDir, "website-source-dir", "templates", "templates directory")
	fs.StringVar(&cmd.tfVersion, "tf-version", "", "terraform binary version to download")
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
	err := provider.Generate(
		cmd.ui,
		cmd.flagLegacySidebar,
		cmd.flagProviderName,
		cmd.flagRenderedProviderName,
		cmd.flagRenderedWebsiteDir,
		cmd.flagExamplesDir,
		cmd.flagWebsiteTmpDir,
		cmd.flagWebsiteSourceDir,
		cmd.tfVersion,
	)
	if err != nil {
		return fmt.Errorf("unable to generate website: %w", err)
	}

	return nil
}
