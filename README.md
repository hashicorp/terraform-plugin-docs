# terraform-plugin-docs

This repository contains tools and packages for creating Terraform plugin docs (currently only provider plugins). The primary way users will interact with this is the **tfpluginwebsite** CLI tool to generate and validate plugin documentation.

## tfpluginwebsite

The **tfpluginwebsite** CLI has two main commands, `validate` and `generate` (`generate` is the default). This tool will let you generate documentation for your provider from live example .tf files and markdown templates. It will also export schema information from the provider (using `terraform providers schema -json`), and sync the schema with the reference documents. If your documentation only consists of simple examples and schema information, the tool can also generate missing template files to make website creation extremely simple for most providers.

### How it Works

When you run `tfpluginwebsite` from root directory of the provider the tool takes the following actions:

* Copy all the templates and static files to a temporary directory
* Build (`go build`) a temporary binary of the provider source code
* Collect schema information using `terraform providers schema -json`
* Generate a default provider template file, if missing (**index.md**)
* Generate resource template files, if missing
* Generate data source template files, if missing
* Copy all non-template files to the output website directory
* Process all the remaining templates to generate files for the output website directory

### Conventional Paths

The generation of missing documentation is based on a number of assumptions / conventional paths:

| Path                                                      | Description                     |
|-----------------------------------------------------------|---------------------------------|
| `templates/`                                              | Root of templated docs          |
| `templates/index.md[.tmpl]`                               | Docs index page (or template)   |
| `examples/provider/provider.tf`                           | Provider example config*        |
| `templates/data-sources/<data source name>.md[.tmpl]`     | Data source page (or template)  |
| `examples/data-sources/<data source name>/data-source.tf` | Data source example config*     |
| `templates/resources/<resource name>.md[.tmpl]`           | Resource page (or template)     |
| `examples/resources/<resource name>/resource.tf`          | Resource example config*        |
| `examples/resources/<resource name>/import.sh`            | Resource example import command |

### Templates

The templates are implemented with Go [`text/template`](https://golang.org/pkg/text/template/) using the following objects and functions:

#### Template Objects

TBD

#### Template Functions

| Function        | Description                                                                                                        |
|-----------------|--------------------------------------------------------------------------------------------------------------------|
| `codefile`      | Create a Markdown code block and populate it with the contents of a file. Path is relative to the repository root. |
| `tffile`        | A special case of the `codefile` function. In addition this will elide lines with an `OMIT` comment.               |
| `trimspace`     | `strings.TrimSpace`                                                                                                |
| `plainmarkdown` | Render Markdown content as plaintext                                                                               |

### Installation

You can install a copy of the binary manually from the releases, or you can optionally use the [tools.go model](https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md) for tool installation.

## Disclaimer

This experimental repository contains software which is still being developed and in the alpha testing stage. It is not ready for production use.
