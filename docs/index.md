---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netorca Provider"
subcategory: ""
description: |-
  Interact with NetOrca.
---

# netorca Provider

Interact with NetOrca.

## Example Usage

```terraform
# Copyright (c) HashiCorp, Inc.

provider "netorca" {
  url    = "https://api.netorca.example.com"
  apikey = "<netorca-api-key>"
}

data "netorca_change_instances" "a_records" {
  pov = "serviceowner"
  filters {
    service_name = "a_record"
  }
}

resource "local_file" "test_file" {
  for_each = { for i in data.netorca_change_instances.a_records.change_instances : jsondecode(i.service_item.declaration).name => jsondecode(i.service_item.declaration) }
  content  = "Zone: ${each.value.zone}\nName: ${each.value.name}\nAddresses: ${join(",", each.value.addresses)}\n"
  filename = each.value.name
}

resource "netorca_change_instances" "example" {
  depends_on = [
    local_file.test_file
  ]
  for_each = { for i in data.netorca_change_instances.a_records.change_instances : i.id => i }
  id       = each.value.id
  state    = "COMPLETED"
  pov      = "serviceowner"
  deployed_item = jsonencode(
    {
      "deployed" : true,
    }
  )
}

output "change_instance_ids" {
  value = [for i in data.netorca_change_instances.a_records.change_instances : i.id]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `apikey` (String) Api-Key for NetOrca API authentication.
- `url` (String) URL for NetOrca API.
