module github.com/hashicorp/terraform-plugin-docs

go 1.14

replace github.com/hashicorp/terraform-exec => ../terraform-exec/

require (
	github.com/hashicorp/terraform-exec v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform-json v0.5.0
	github.com/mattn/go-colorable v0.1.7
	github.com/mitchellh/cli v1.1.1
	github.com/russross/blackfriday v1.5.2
	github.com/zclconf/go-cty v1.4.1
)
