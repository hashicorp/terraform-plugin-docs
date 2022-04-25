# 0.8.0 (April 25, 2022)

ENHANCEMENTS:

* template functions: Added the `split` command to split a string into substrings
 
# 0.7.0 (March 15, 2022)

ENHANCEMENTS:

* cmd/tfplugindocs: Use existing Terraform CLI binary if available on PATH, otherwise download latest Terraform CLI binary (https://github.com/hashicorp/terraform-plugin-docs/pull/124)
* cmd/tfplugindocs: Added `tf-version` flag for specifying Terraform CLI binary version to download, superseding the PATH lookup (https://github.com/hashicorp/terraform-plugin-docs/pull/124)

BUG FIXES:

* cmd/tfplugindocs: Swapped `.Type` and `.Name` resource and data source template fields so they correctly align (https://github.com/hashicorp/terraform-plugin-docs/pull/44)
* schemamd: Switched attribute name rendering from bold text to code blocks so the Terraform Registry treats them as anchor links (https://github.com/hashicorp/terraform-plugin-docs/pull/59)

# 0.6.0 (March 14, 2022)

NOTES:

* The `github.com/hashicorp/terraform-exec` dependency has been updated to match `terraform-plugin-sdk`, which replaced the removed `tfinstall` package with `github.com/hashicorp/hc-install`. This will resolve Go build errors for projects that import both `terraform-plugin-docs` and `terraform-plugin-sdk`.
