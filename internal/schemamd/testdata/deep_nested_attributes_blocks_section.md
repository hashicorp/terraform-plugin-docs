## Schema

### Required Attributes

- `level_one` (Attributes) (see [below for nested schema](#nestedatt--level_one))

### Read-Only

- `id` (String) Example identifier

<a id="nestedatt--level_one"></a>
### Nested Schema for `level_one`

Optional Attributes:

- `level_two` (Attributes) (see [below for nested schema](#nestedatt--level_one--level_two))

<a id="nestedatt--level_one--level_two"></a>
### Nested Schema for `level_one.level_two`

Optional Attributes:

- `level_three` (Attributes) (see [below for nested schema](#nestedatt--level_one--level_two--level_three))

<a id="nestedatt--level_one--level_two--level_three"></a>
### Nested Schema for `level_one.level_two.level_three`

Optional Attributes:

- `level_four_primary` (Attributes) (see [below for nested schema](#nestedatt--level_one--level_two--level_three--level_four_primary))
- `level_four_secondary` (String)

<a id="nestedatt--level_one--level_two--level_three--level_four_primary"></a>
### Nested Schema for `level_one.level_two.level_three.level_four_primary`

Optional Attributes:

- `level_five` (Attributes) Parent should be level_one.level_two.level_three.level_four_primary. (see [below for nested schema](#nestedatt--level_one--level_two--level_three--level_four_primary--level_five))
- `level_four_primary_string` (String) Parent should be level_one.level_two.level_three.level_four_primary.

<a id="nestedatt--level_one--level_two--level_three--level_four_primary--level_five"></a>
### Nested Schema for `level_one.level_two.level_three.level_four_primary.level_five`

Optional Attributes:

- `level_five_string` (String) Parent should be level_one.level_two.level_three.level_four_primary.level_five.
