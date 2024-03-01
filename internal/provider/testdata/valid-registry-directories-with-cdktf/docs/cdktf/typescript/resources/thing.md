---
subcategory: "Example"
page_title: "Example: example_thing"
description: |-
  Example description.
---

# Resource: example_thing

Byline.

## Example Usage

```ts
import { Construct } from "construct";
import { TerraformStack } from "cdktf";
import { Thing } from "./.gen/providers/example/thing";

class MyStack extends TerraformStack {
  constructs(scope: Construct, name: string) {
    super(scope, name);

    new Thing(this, "example", {
      name: "example",
    });
  }
}
```

## Argument Reference

- `name` - (Required) Name of thing.

## Attribute Reference

- `id` - Name of thing.
