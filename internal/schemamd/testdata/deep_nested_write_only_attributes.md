## Schema

### Required

- `level_one` (Attributes) (see [below for nested schema](#nestedatt--level_one))

### Read-Only

- `id` (String) Example identifier

<a id="nestedatt--level_one"></a>
### Nested Schema for `level_one`

Optional:

- `level_two` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) (see [below for nested schema](#nestedatt--level_one--level_two))

<a id="nestedatt--level_one--level_two"></a>
### Nested Schema for `level_one.level_two`

Optional:

- `level_three` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) (see [below for nested schema](#nestedatt--level_one--level_two--level_three))

<a id="nestedatt--level_one--level_two--level_three"></a>
### Nested Schema for `level_one.level_two.level_three`

Optional:

- `level_four_primary` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) (see [below for nested schema](#nestedatt--level_one--level_two--level_three--level_four_primary))
- `level_four_secondary` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments))

<a id="nestedatt--level_one--level_two--level_three--level_four_primary"></a>
### Nested Schema for `level_one.level_two.level_three.level_four_primary`

Optional:

- `level_five` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Parent should be level_one.level_two.level_three.level_four_primary. (see [below for nested schema](#nestedatt--level_one--level_two--level_three--level_four_primary--level_five))
- `level_four_primary_string` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Parent should be level_one.level_two.level_three.level_four_primary.

<a id="nestedatt--level_one--level_two--level_three--level_four_primary--level_five"></a>
### Nested Schema for `level_one.level_two.level_three.level_four_primary.level_five`

Optional:

- `level_five_string` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Parent should be level_one.level_two.level_three.level_four_primary.level_five.
