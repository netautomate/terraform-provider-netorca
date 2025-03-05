# Copyright (c) HashiCorp, Inc.

provider "netorca" {
  url    = "https://api.netorca.example.com"
  apikey = "<netorca-api-key>"
}

data "netorca_service_items" "completed_changes" {
  pov = "serviceowner"
  filters {
    change_state = "ALL_CHANGES_COMPLETED"
  }
}

output "service_item_ids" {
  value = [for i in data.netorca_service_items.service_items.completed_changes : i.id]
}

output "service_item_declarations" {
  value = [for i in data.netorca_service_items.service_items.completed_changes : i.declaration]
}
