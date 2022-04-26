## Schema

### Optional

- `kms_key_id` (String) The Amazon Resource Name (ARN) of the CMK to use when encrypting log data.
- `log_group_name` (String) The name of the log group. If you don't specify a name, AWS CloudFormation generates a unique ID for the log group.
- `retention_in_days` (Number) The number of days to retain the log events in the specified log group. Possible values are: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, and 3653.

### Read-Only

- `arn` (String) The CloudWatch log group ARN.
- `id` (String) Uniquely identifies the resource.