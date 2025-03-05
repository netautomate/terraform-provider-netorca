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
