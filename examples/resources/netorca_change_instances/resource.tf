# Copyright (c) HashiCorp, Inc.

provider "netorca" {
  url    = "https://api.netorca.example.com"
  apikey = "<netorca-api-key>"
}

data "netorca_change_instances" "change_instances" {
  pov = "serviceowner"
  filters {
    service_id = 12
  }
}

resource "local_file" "test_file" {
  for_each = { for i in data.netorca_change_instances.change_instances.change_instances : jsondecode(i.service_item.declaration).name => jsondecode(i.service_item.declaration) }
  content  = "Zone: ${each.value.zone}\nName: ${each.value.name}\nAddresses: ${join(",", each.value.addresses)}\n"
  filename = each.value.name
}


resource "netorca_change_instances" "example" {
  depends_on = [
    local_file.test_file
  ]
  for_each = { for i in data.netorca_change_instances.change_instances.change_instances : i.id => i }
  id       = each.value.id
  state    = "APPROVED"
  pov      = "serviceowner"
  deployed_item = jsonencode(
    {
      "deployed" : true,
    }
  )
}

output "change_instances" {
  value = [for i in resource.netorca_change_instances.example : i]
}

