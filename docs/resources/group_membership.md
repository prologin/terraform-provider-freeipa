---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "freeipa_group_membership Resource - terraform-provider-freeipa"
subcategory: ""
description: |-
  Manage FreeIPA group membership
---

# freeipa_group_membership (Resource)

Manage FreeIPA group membership

## Example Usage

```terraform
resource "freeipa_group_membership" "managers" {
  group  = "ipausers"
  member = "root"
  type   = "user"

  manager = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group` (String) Group name (CN)
- `member` (String) Member identifier

### Optional

- `manager` (Boolean) The member is a manager of the group (connot be used with type service)
- `type` (String) Member type (must be one of "user", "group" or "service")

### Read-Only

- `id` (String) The ID of this resource.


