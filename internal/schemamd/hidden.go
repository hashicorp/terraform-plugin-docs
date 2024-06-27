package schemamd

// hiddenFields is the list of fields we don't want to appear in the documentation
// because users should not have to interact with them.
var hiddenFields = []string{
	"id",
	"kind",
	"metadata.namespace",
	"metadata.revision",
}
