---
subcategory: "Example"
layout: "example"
page_title: "Example: example_thing"
description: |-
  Example description.
---

# Data Source: example_thing

Byline.

## Example Usage

```ts
import { Construct } from "construct";
import { TerraformStack } from "cdktf";
import { DataExample } from "./.gen/providers/example/data_example_thing";

class MyStack extends TerraformStack {
  constructs(scope: Construct, name: string) {
    super(scope, name);

    new DataExampleThing(this, "example", {
      name: "example",
    });
  }
}
```

## Argument Reference

- `name` - (Required) Name of thing.

## Attribute Reference

- `id` - Name of thing.
