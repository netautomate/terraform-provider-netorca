---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netorca_change_instances Data Source - netorca"
subcategory: ""
description: |-
  Use this data provider to return a list of change instances.
---

# netorca_change_instances (Data Source)

Use this data provider to return a list of change instances.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `pov` (String) The POV from which to make the request (serviceowner|consumer).

### Optional

- `filters` (Block, Optional) (see [below for nested schema](#nestedblock--filters))

### Read-Only

- `change_instance_count` (Number) The number of change instances that the request has matched.
- `change_instances` (Block List) (see [below for nested schema](#nestedblock--change_instances))

<a id="nestedblock--filters"></a>
### Nested Schema for `filters`

Optional:

- `application_id` (Number)
- `change_type` (String)
- `commit_id` (String)
- `consumer_team_id` (Number)
- `limit` (Number)
- `offset` (Number)
- `ordering` (String)
- `service_id` (Number)
- `service_item_id` (Number)
- `service_name` (String)
- `service_owner_team_id` (Number)
- `state` (String)
- `submission_id` (Number)


<a id="nestedblock--change_instances"></a>
### Nested Schema for `change_instances`

Read-Only:

- `consumer_team` (Object) (see [below for nested schema](#nestedatt--change_instances--consumer_team))
- `created` (String)
- `id` (Number)
- `modified` (String)
- `owner` (Object) (see [below for nested schema](#nestedatt--change_instances--owner))
- `service_item` (Object) (see [below for nested schema](#nestedatt--change_instances--service_item))
- `state` (String)
- `submission` (Object) (see [below for nested schema](#nestedatt--change_instances--submission))
- `url` (String)

<a id="nestedatt--change_instances--consumer_team"></a>
### Nested Schema for `change_instances.consumer_team`

Read-Only:

- `id` (Number)
- `metadata` (String)
- `name` (String)


<a id="nestedatt--change_instances--owner"></a>
### Nested Schema for `change_instances.owner`

Read-Only:

- `id` (Number)
- `name` (String)


<a id="nestedatt--change_instances--service_item"></a>
### Nested Schema for `change_instances.service_item`

Read-Only:

- `application` (Object) (see [below for nested schema](#nestedobjatt--change_instances--service_item--application))
- `change_state` (String)
- `created` (String)
- `declaration` (String)
- `deployed_item` (String)
- `id` (Number)
- `modified` (String)
- `name` (String)
- `runtime_state` (String)
- `service` (Object) (see [below for nested schema](#nestedobjatt--change_instances--service_item--service))
- `url` (String)

<a id="nestedobjatt--change_instances--service_item--application"></a>
### Nested Schema for `change_instances.service_item.application`

Read-Only:

- `id` (Number)
- `metadata` (String)
- `name` (String)
- `owner` (Number)


<a id="nestedobjatt--change_instances--service_item--service"></a>
### Nested Schema for `change_instances.service_item.service`

Read-Only:

- `healthcheck` (Boolean)
- `id` (Number)
- `name` (String)
- `owner` (Object) (see [below for nested schema](#nestedobjatt--change_instances--service_item--service--owner))

<a id="nestedobjatt--change_instances--service_item--service--owner"></a>
### Nested Schema for `change_instances.service_item.service.owner`

Read-Only:

- `id` (Number)
- `name` (String)




<a id="nestedatt--change_instances--submission"></a>
### Nested Schema for `change_instances.submission`

Read-Only:

- `commit_id` (String)
- `id` (Number)
