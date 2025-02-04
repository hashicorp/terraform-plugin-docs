## Schema

### Required

- `level_one` (Attributes) (See [below for nested schema](#nestedatt--level_one))

### Read-Only

- `id` (String) Example identifier

<a id="nestedatt--level_one"></a>
### Nested Schema for `level_one`

Optional:

- `level_two` (Attributes) (See [below for nested schema](#nestedatt--level_one--level_two))

<a id="nestedatt--level_one--level_two"></a>
### Nested Schema for `level_one.level_two`

Optional:

- `level_three` (Attributes) (See [below for nested schema](#nestedatt--level_one--level_two--level_three))

<a id="nestedatt--level_one--level_two--level_three"></a>
### Nested Schema for `level_one.level_two.level_three`

Optional:

- `level_four_primary` (Attributes) (See [below for nested schema](#nestedatt--level_one--level_two--level_three--level_four_primary))
- `level_four_secondary` (String)

<a id="nestedatt--level_one--level_two--level_three--level_four_primary"></a>
### Nested Schema for `level_one.level_two.level_three.level_four_primary`

Optional:

- `level_five` (Attributes) Parent should be level_one.level_two.level_three.level_four_primary. (See [below for nested schema](#nestedatt--level_one--level_two--level_three--level_four_primary--level_five))
- `level_four_primary_string` (String) Parent should be level_one.level_two.level_three.level_four_primary.

<a id="nestedatt--level_one--level_two--level_three--level_four_primary--level_five"></a>
### Nested Schema for `level_one.level_two.level_three.level_four_primary.level_five`

Optional:

- `level_five_string` (String) Parent should be level_one.level_two.level_three.level_four_primary.level_five.
