# Copyright (c) HashiCorp, Inc.

provider "netorca" {
  url    = "https://api.netorca.example.com"
  apikey = "<netorca-api-key>"
}

data "netorca_change_instances" "a_record_changes" {
  pov = "serviceowner"
  filters {
    service_name = "a_record"
  }
}

output "change_instance_ids" {
  value = [for i in data.netorca_change_instances.a_record_changes.change_instances : i.id]
}
