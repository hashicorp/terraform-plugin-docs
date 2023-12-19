# terraform-plugin-docs

This repository contains tools and packages for creating Terraform plugin docs (currently only provider plugins).
The primary way users will interact with this is the `tfplugindocs` CLI tool to generate and validate plugin documentation.

## `tfplugindocs`

The `tfplugindocs` CLI has two main commands, `validate` and `generate` (`generate` is the default).
This tool will let you generate documentation for your provider from live example `.tf` files and markdown templates.
It will also export schema information from the provider (using `terraform providers schema -json`),
and sync the schema with the reference documents.

If your documentation only consists of simple examples and schema information,
the tool can also generate missing template files to make website creation extremely simple for most providers.

### Installation

You can install a copy of the binary manually from the [releases](https://github.com/hashicorp/terraform-plugin-docs/releases),
or you can optionally use the [tools.go model](https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md)
for tool installation.

### Usage

```shell
$ tfplugindocs --help
Usage: tfplugindocs [--version] [--help] <command> [<args>]

Available commands are:
                the generate command is run by default
    generate    generates a plugin website from code, templates, and examples
    validate    validates a plugin website for the current directory
       
```

`generate` command:

```shell
$ tfplugindocs generate --help

Usage: tfplugindocs generate [<args>]

    --examples-dir <ARG>             examples directory based on provider-dir                                                                                           (default: "examples")
    --ignore-deprecated <ARG>        don't generate documentation for deprecated resources and data-sources                                                             (default: "false")
    --provider-dir <ARG>             relative or absolute path to the root provider code directory when running the command outside the root provider code directory  
    --provider-name <ARG>            provider name, as used in Terraform configurations                                                                               
    --providers-schema <ARG>         path to the providers schema JSON file, which contains the output of the terraform providers schema -json command. Setting this flag will skip building the provider and calling Terraform CLI                                                                               
    --rendered-provider-name <ARG>   provider name, as generated in documentation (ex. page titles, ...)                                                              
    --rendered-website-dir <ARG>     output directory based on provider-dir                                                                                             (default: "docs")
    --tf-version <ARG>               terraform binary version to download                                                                                             
    --website-source-dir <ARG>       templates directory based on provider-dir                                                                                          (default: "templates")
    --website-temp-dir <ARG>         temporary directory (used during generation)  
```

`validate` command:

```shell
$ tfplugindocs validate --help

Usage: tfplugindocs validate [<args>]
```

`migrate` command:

```shell
$ tfplugindocs migrate --help

Usage: tfplugindocs migrate [<args>]

    --examples-dir <ARG>             examples directory based on provider-dir                                                                                           (default: "examples")
    --old-website-source-dir <ARG>   old website directory based on provider-dir; files will be migrated from this directory                                            (default: "website")
    --provider-dir <ARG>             relative or absolute path to the root provider code directory when running the command outside the root provider code directory
    --templates-dir <ARG>            new website templates directory based on provider-dir; files will be migrated to this directory                                    (default: "templates")
```

### How it Works

When you run `tfplugindocs`, by default from the root directory of a provider codebase, the tool takes the following actions:

* Copy all the templates and static files to a temporary directory
* Build (`go build`) a temporary binary of the provider source code
* Collect schema information using `terraform providers schema -json`
* Generate a default provider template file, if missing (**index.md**)
* Generate resource template files, if missing
* Generate data source template files, if missing
* Copy all non-template files to the output website directory
* Process all the remaining templates to generate files for the output website directory

For inspiration, you can look at the templates and output of the
[`terraform-provider-random`](https://github.com/hashicorp/terraform-provider-random)
and [`terraform-provider-tls`](https://github.com/hashicorp/terraform-provider-tls).
You can browse their respective docs on the Terraform Registry,
[here](https://registry.terraform.io/providers/hashicorp/random/latest/docs)
and [here](https://registry.terraform.io/providers/hashicorp/tls/latest/docs).

#### About the `id` attribute

If the provider schema didn't set `id` for the given resource/data-source, the documentation generated
will place it under the "Read Only" section and provide a simple description.

Otherwise, the provider developer can set an arbitrary description like this:

```golang
    // ...
    Schema: map[string]*schema.Schema{
        // ...
        "id": {
            Type:     schema.TypeString,
            Computed: true,
            Description: "Unique identifier for this resource",
		},
        // ...
    }
    // ...
```

#### Migrate subcommand

The `migrate` subcommand can be used to migrate from the old directory structure (`website/docs/r`) to the `tfplugindocs` supported structure (`templates/`).

The `migrate` subcommand takes the following actions:
- Copies the `--old-website-source-dir` folder to the `--tempates-dir` folder (will create this folder if it doesn't exist)
- Renames `docs/d/` and `docs/r/` subdirectories to `data-sources/` and `resources/` respectively
- Change file suffixes for template files from `.html.markdown` to `.md.tmpl`
- Extracts code blocks from website docs to create individual example files in `--examples-dir` (will create this folder if it doesn't exist)
- replace extracted example code in website templates with `tfplugindocs` template code referring to example files.
- Copies non-template files to `--templates-dir` folder

### Conventional Paths

The generation of missing documentation is based on a number of assumptions / conventional paths.

> **NOTE:** In the following conventional paths, `<data source name>` and `<resource name>` include the provider prefix as well.
> For example, the data source [`caller_identity`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) in the `aws` provider would have an "example" conventional path of: `examples/data-sources/aws_caller_identity/data-source.tf`

For templates:

| Path                                                      | Description                            |
|-----------------------------------------------------------|----------------------------------------|
| `templates/`                                              | Root of templated docs                 |
| `templates/index.md[.tmpl]`                               | Docs index page (or template)          |
| `templates/data-sources.md[.tmpl]`                        | Generic data source page (or template) |
| `templates/data-sources/<data source name>.md[.tmpl]`     | Data source page (or template)         |
| `templates/resources.md[.tmpl]`                           | Generic resource page (or template)    |
| `templates/resources/<resource name>.md[.tmpl]`           | Resource page (or template)            |

Note: the `.tmpl` extension is necessary, for the file to be correctly handled as a template.

For examples:

| Path                                                      | Description                     |
|-----------------------------------------------------------|---------------------------------|
| `examples/`                                               | Root of examples                |
| `examples/provider/provider.tf`                           | Provider example config         |
| `examples/data-sources/<data source name>/data-source.tf` | Data source example config      |
| `examples/resources/<resource name>/resource.tf`          | Resource example config         |
| `examples/resources/<resource name>/import.sh`            | Resource example import command |

#### Migration

The `migrate` subcommand assumes the following conventional paths for the old website structure:

| Path                                              | Description          |
|---------------------------------------------------|----------------------|
| `website/`                                        | Root of website docs |
| `website/docs/index.html.markdown`                | Docs index page      |
| `website/docs/d/<data source name>.html.markdown` | Data source page     |
| `website/docs/r/<resource name>.html.markdown`    | Resource page        |

### Templates

The templates are implemented with Go [`text/template`](https://golang.org/pkg/text/template/)
using the following data fields and functions:

#### Data fields

##### Provider

|                   Field |  Type  | Description                                                                               |
|------------------------:|:------:|-------------------------------------------------------------------------------------------|
|          `.Description` | string | Provider description                                                                      |
|           `.HasExample` |  bool  | Is there an example file?                                                                 |
|          `.ExampleFile` | string | Path to the file with the terraform configuration example                                 |
|         `.ProviderName` | string | Canonical provider name (ex. `terraform-provider-random`)                                 |
|    `.ProviderShortName` | string | Short version of the provider name (ex. `random`)                                         |
| `.RenderedProviderName` | string | Value provided via argument `--rendered-provider-name`, otherwise same as `.ProviderName` |
|       `.SchemaMarkdown` | string | a Markdown formatted Provider Schema definition                                           |

##### Resources / Data Source

|                   Field |  Type  | Description                                                                               |
|------------------------:|:------:|-------------------------------------------------------------------------------------------|
|                 `.Name` | string | Name of the resource/data-source (ex. `tls_certificate`)                                  |
|                 `.Type` | string | Either `Resource` or `Data Source`                                                        |
|          `.Description` | string | Resource / Data Source description                                                        |
|           `.HasExample` |  bool  | Is there an example file?                                                                 |
|          `.ExampleFile` | string | Path to the file with the terraform configuration example                                 |
|            `.HasImport` |  bool  | Is there an import file?                                                                  |
|           `.ImportFile` | string | Path to the file with the command for importing the resource                              |
|         `.ProviderName` | string | Canonical provider name (ex. `terraform-provider-random`)                                 |
|    `.ProviderShortName` | string | Short version of the provider name (ex. `random`)                                         |
| `.RenderedProviderName` | string | Value provided via argument `--rendered-provider-name`, otherwise same as `.ProviderName` |
|       `.SchemaMarkdown` | string | a Markdown formatted Resource / Data Source Schema definition                             |

#### Functions

| Function        | Description                                                                                       |
|-----------------|---------------------------------------------------------------------------------------------------|
| `codefile`      | Create a Markdown code block with the content of a file. Path is relative to the repository root. |
| `lower`         | Equivalent to [`strings.ToLower`](https://pkg.go.dev/strings#ToLower).                            |
| `plainmarkdown` | Render Markdown content as plaintext.                                                             |
| `prefixlines`   | Add a prefix to all (newline-separated) lines in a string.                                        |
| `printf`        | Equivalent to [`fmt.Printf`](https://pkg.go.dev/fmt#Printf).                                      |
| `split`         | Split string into sub-strings, by a given separator (ex. `split .Name "_"`).                      |
| `title`         | Equivalent to [`cases.Title`](https://pkg.go.dev/golang.org/x/text/cases#Title).                  |
| `tffile`        | A special case of the `codefile` function, designed for Terraform files (i.e. `.tf`).             |
| `trimspace`     | Equivalent to [`strings.TrimSpace`](https://pkg.go.dev/strings#TrimSpace).                        |
| `upper`         | Equivalent to [`strings.ToUpper`](https://pkg.go.dev/strings#ToUpper).                            |
 
## Disclaimer

This is still under development: while it's being used for production-ready providers, you might still find bugs
and limitations. In those cases, please report [issues](https://github.com/hashicorp/terraform-plugin-docs/issues)
or, if you can, submit a [pull-request](https://github.com/hashicorp/terraform-plugin-docs/pulls).

Your help and patience is truly appreciated.

## Contributing

### License Headers
All source code files in this repository (excluding autogenerated files like `go.mod`, prose, and files excluded in [.copywrite.hcl](.copywrite.hcl)) must have a license header at the top.

This can be autogenerated by running `make generate` or running `go generate ./...` in the [/tools](/tools) directory.

### Acceptance Tests

This repo uses the `testscript` [package](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript) for acceptance testing.

There are two types of acceptance tests: full provider build tests in `tfplugindocs/testdata/scripts/provider-build` and provider schema json tests in `tfplugindocs/testdata/scripts/schema-json`.

Provider build tests run the default `tfplugindocs` command which builds the provider source code and runs Terraform to retrieve the schema. These tests require the full provider source code to build a valid provider binary. 

Schema json tests run the `tfplugindocs` command with the `--providers-schema=<arg>` flag to specify a provider schemas json file. This allows the test to skip the provider build and Terraform CLI call, instead using the specified file to generate docs. 

You can run `make testacc` to run the full suite of acceptance tests. By default, the provider build acceptance tests will create a temporary directory in `/tmp/tftmp` for testing, but you can change this location in `cmd/tfplugindocs/main_test.go`. The schema json tests uses the `testscript` package's [default work directory](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript#Params.WorkdirRoot).

The test scripts are defined in the `tfplugindocs/testdata/scripts` directory. Each script includes the test, golden files, and the provider source code or schema JSON file needed to run the test.

Each script is a [text archive](https://pkg.go.dev/golang.org/x/tools/txtar). You can install the `txtar` CLI locally by running `go install golang.org/x/exp/cmd/txtar@latest` to extract the files in the test script for debugging. 
You can also use `txtar` CLI archive files into the `.txtar` format to create new tests or modify existing ones.
