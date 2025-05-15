## Schema

### Optional

- `certificate_authority_arn` (String)
- `certificate_body` (String)
- `certificate_chain` (String)
- `domain_name` (String)
- `options` (Block List, Max: 1) (see [below for nested schema](#nestedblock--options))
- `private_key` (String, Sensitive)
- `subject_alternative_names` (Set of Strings)
- `tags` (Map of Strings)
- `tags_all` (Map of Strings)
- `validation_method` (String)

### Read-Only

- `arn` (String)
- `domain_validation_options` (Set of Objects) (see [below for nested schema](#nestedatt--domain_validation_options))
- `id` (String) The ID of this resource.
- `status` (String)
- `validation_emails` (List of Strings)

<a id="nestedblock--options"></a>
### Nested Schema for `options`

Optional:

- `certificate_transparency_logging_preference` (String)


<a id="nestedatt--domain_validation_options"></a>
### Nested Schema for `domain_validation_options`

Read-Only:

- `domain_name` (String)
- `resource_record_name` (String)
- `resource_record_type` (String)
- `resource_record_value` (String)
