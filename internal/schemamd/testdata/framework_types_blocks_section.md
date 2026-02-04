## Schema

### Required Attributes

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `required_write_only_string_attribute` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) example required write-only string attribute

### Optional Attributes

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `bool_attribute` (Boolean) example bool attribute
- `float64_attribute` (Number) example float64 attribute
- `int64_attribute` (Number) example int64 attribute
- `list_attribute` (List of String) example list attribute
- `map_attribute` (Map of String) example map attribute
- `number_attribute` (Number) example number attribute
- `object_attribute` (Object) example object attribute (see [below for nested schema](#nestedatt--object_attribute))
- `object_attribute_with_nested_object_attribute` (Object) example object attribute with nested object attribute (see [below for nested schema](#nestedatt--object_attribute_with_nested_object_attribute))
- `sensitive_bool_attribute` (Boolean, Sensitive) example sensitive bool attribute
- `sensitive_float64_attribute` (Number, Sensitive) example sensitive float64 attribute
- `sensitive_int64_attribute` (Number, Sensitive) example sensitive int64 attribute
- `sensitive_list_attribute` (List of String, Sensitive) example sensitive list attribute
- `sensitive_map_attribute` (Map of String, Sensitive) example sensitive map attribute
- `sensitive_number_attribute` (Number, Sensitive) example sensitive number attribute
- `sensitive_object_attribute` (Object, Sensitive) example sensitive object attribute (see [below for nested schema](#nestedatt--sensitive_object_attribute))
- `sensitive_set_attribute` (Set of String, Sensitive) example sensitive set attribute
- `sensitive_string_attribute` (String, Sensitive) example sensitive string attribute
- `set_attribute` (Set of String) example set attribute
- `string_attribute` (String) example string attribute
- `write_only_string_attribute` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) example write only string attribute

### Blocks

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `list_nested_block` (Block List) example list nested block (see [below for nested schema](#nestedblock--list_nested_block))
- `list_nested_block_sensitive_nested_attribute` (Block List) (see [below for nested schema](#nestedblock--list_nested_block_sensitive_nested_attribute))
- `set_nested_block` (Block Set) example set nested block (see [below for nested schema](#nestedblock--set_nested_block))
- `single_nested_block` (Block, Optional) example single nested block (see [below for nested schema](#nestedblock--single_nested_block))
- `single_nested_block_sensitive_nested_attribute` (Block, Optional) example sensitive single nested block (see [below for nested schema](#nestedblock--single_nested_block_sensitive_nested_attribute))

### Read-Only

- `id` (String) The ID of this resource.
- `set_nested_block_sensitive_nested_attribute` (Block Set) example sensitive set nested block (see [below for nested schema](#nestedblock--set_nested_block_sensitive_nested_attribute))

<a id="nestedatt--object_attribute"></a>
### Nested Schema for `object_attribute`

Optional Attributes:

- `object_attribute_attribute` (String)


<a id="nestedatt--object_attribute_with_nested_object_attribute"></a>
### Nested Schema for `object_attribute_with_nested_object_attribute`

Optional Attributes:

- `nested_object` (Object) (see [below for nested schema](#nestedobjatt--object_attribute_with_nested_object_attribute--nested_object))
- `object_attribute_attribute` (String)

<a id="nestedobjatt--object_attribute_with_nested_object_attribute--nested_object"></a>
### Nested Schema for `object_attribute_with_nested_object_attribute.nested_object`

Optional Attributes:

- `nested_object_attribute` (String)



<a id="nestedatt--sensitive_object_attribute"></a>
### Nested Schema for `sensitive_object_attribute`

Optional Attributes:

- `object_attribute_attribute` (String)


<a id="nestedblock--list_nested_block"></a>
### Nested Schema for `list_nested_block`

Optional Attributes:

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `list_nested_block_attribute` (String) example list nested block attribute
- `list_nested_block_attribute_with_default` (String) example list nested block attribute with default
- `list_nested_block_write_only_attribute` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) example list nested block write-only attribute

Blocks:

- `nested_list_block` (Block List) (see [below for nested schema](#nestedblock--list_nested_block--nested_list_block))

<a id="nestedblock--list_nested_block--nested_list_block"></a>
### Nested Schema for `list_nested_block.nested_list_block`

Optional Attributes:

- `nested_block_string_attribute` (String) example nested block string attribute



<a id="nestedblock--list_nested_block_sensitive_nested_attribute"></a>
### Nested Schema for `list_nested_block_sensitive_nested_attribute`

Optional Attributes:

- `list_nested_block_attribute` (String) example list nested block attribute
- `list_nested_block_sensitive_attribute` (String, Sensitive) example sensitive list nested block attribute


<a id="nestedblock--set_nested_block"></a>
### Nested Schema for `set_nested_block`

Optional Attributes:

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `set_nested_block_attribute` (String) example set nested block attribute
- `set_nested_block_write_only_attribute` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) example set nested block write-only attribute


<a id="nestedblock--single_nested_block"></a>
### Nested Schema for `single_nested_block`

Optional Attributes:

- `single_nested_block_attribute` (String) example single nested block attribute


<a id="nestedblock--single_nested_block_sensitive_nested_attribute"></a>
### Nested Schema for `single_nested_block_sensitive_nested_attribute`

Optional Attributes:

- `single_nested_block_attribute` (String) example single nested block attribute
- `single_nested_block_sensitive_attribute` (String, Sensitive) example sensitive single nested block attribute


<a id="nestedblock--set_nested_block_sensitive_nested_attribute"></a>
### Nested Schema for `set_nested_block_sensitive_nested_attribute`

Read-Only:

- `set_nested_block_attribute` (String) example set nested block attribute
- `set_nested_block_sensitive_attribute` (String, Sensitive) example sensitive set nested block attribute
