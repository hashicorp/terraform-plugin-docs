# 0.12.0 (June 29, 2022)

BUG FIXES:

* template data: A regression was introduced in [#155](https://github.com/hashicorp/terraform-plugin-docs/pull/155) making template data field `HasExample` and `HasImport` always true ([#162](https://github.com/hashicorp/terraform-plugin-docs/pull/162)).

NEW FEATURES:

* template functions: Added `lower`, `upper` and `title` ([#162](https://github.com/hashicorp/terraform-plugin-docs/pull/162)).

ENHANCEMENTS:

* Added documentation for all the template functions and template data fields ([#162](https://github.com/hashicorp/terraform-plugin-docs/pull/162)).

# 0.11.0 (June 28, 2022)

NEW FEATURES:

* cmd/tfplugindocs: Additional CLI argument `ignore-deprecated` allows to skip deprecated resources and data-sources when generating docs ([#154](https://github.com/hashicorp/terraform-plugin-docs/pull/154)).

BUG FIXES:

* cmd/tfplugindocs: Pass through filepaths for `examples` and `import` to allow use of `HasExample` and `HasImport` template helpers in custom templates ([#155](https://github.com/hashicorp/terraform-plugin-docs/pull/155)).
* cmd/tfplugindocs: Fixed issue with the generation of title and reference links, when nested attributes go too deep ([#56](https://github.com/hashicorp/terraform-plugin-docs/pull/56)).

# 0.10.1 (June 14, 2022)

BUG FIXES:

* cmd/tfplugindocs: Do not error when schema not found, issue log warning ([#151](https://github.com/hashicorp/terraform-plugin-docs/pull/151)).

# 0.10.0 (June 13, 2022)

BUG FIXES:

* cmd/tfplugindocs: Allow single word resources to use templates ([#147](https://github.com/hashicorp/terraform-plugin-docs/pull/147)).
* cmd/tfplugindocs: Pass in correct provider name for data-source and resource schema lookup when overidden with `rendered-provider-name` flag ([#148](https://github.com/hashicorp/terraform-plugin-docs/pull/148)).

ENHANCEMENTS:

* cmd/tfplugindocs: Expose `RenderedProviderName` to templates ([#149](https://github.com/hashicorp/terraform-plugin-docs/pull/149)).

# 0.9.0 (June 1, 2022)

NEW FEATURES:

* cmd/tfplugindocs: Additional CLI arguments `provider-name`, `rendered-provider-name`, `rendered-website-dir`, `examples-dir`, `website-temp-dir`, and `website-source-dir`. These allow to further customise generated doc ([#95](https://github.com/hashicorp/terraform-plugin-docs/pull/95)).

ENHANCEMENTS:

* cmd/tfplugindocs: Implemented usage output (i.e. `--help`) for `generate` and `validate` commands ([#95](https://github.com/hashicorp/terraform-plugin-docs/pull/95)).

# 0.8.1 (May 10, 2022)

BUG FIXES:

* cmd/tfplugindocs: Updated version of [hc-install](github.com/hashicorp/hc-install) in response to change in HashiCorp Release API [sending back a different `Content-Type` header](https://github.com/hashicorp/hc-install/issues/56). This was causing failures when the tool attempted to install Terraform. ([#135](https://github.com/hashicorp/terraform-plugin-docs/issues/135))

# 0.8.0 (May 3, 2022)

ENHANCEMENTS:

* template functions: Added `split` to help separating a string into substrings ([#70](https://github.com/hashicorp/terraform-plugin-docs/pull/70)).

BUG FIXES:

* cmd/tflugindocs: Support for schemas containing empty nested attributes or empty nested blocks ([#99](https://github.com/hashicorp/terraform-plugin-docs/pull/99), [#134](https://github.com/hashicorp/terraform-plugin-docs/pull/134)).
* schemamd: Attribute `ID` is considered "Read Only", unless there's a description defined, in which case it's handled like any other attribute in the schema ([#46](https://github.com/hashicorp/terraform-plugin-docs/pull/46), [#134](https://github.com/hashicorp/terraform-plugin-docs/pull/134)).

# 0.7.0 (March 15, 2022)

ENHANCEMENTS:

* cmd/tfplugindocs: Use existing Terraform CLI binary if available on PATH, otherwise download latest Terraform CLI binary (https://github.com/hashicorp/terraform-plugin-docs/pull/124).
* cmd/tfplugindocs: Added `tf-version` flag for specifying Terraform CLI binary version to download, superseding the PATH lookup (https://github.com/hashicorp/terraform-plugin-docs/pull/124).

BUG FIXES:

* cmd/tfplugindocs: Swapped `.Type` and `.Name` resource and data source template fields so they correctly align (https://github.com/hashicorp/terraform-plugin-docs/pull/44).
* schemamd: Switched attribute name rendering from bold text to code blocks so the Terraform Registry treats them as anchor links (https://github.com/hashicorp/terraform-plugin-docs/pull/59).

# 0.6.0 (March 14, 2022)

NOTES:

* dependencies: `github.com/hashicorp/terraform-exec` dependency has been updated to match `terraform-plugin-sdk`, which replaced the removed `tfinstall` package with `github.com/hashicorp/hc-install`. This will resolve Go build errors for projects that import both `terraform-plugin-docs` and `terraform-plugin-sdk`.
