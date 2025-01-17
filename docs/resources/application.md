---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fortitokencloud_application Resource - fortitokencloud"
subcategory: ""
description: |-
  
---

# fortitokencloud_application (Resource)

## Example Usage
```terraform
resource "fortitokencloud_application" "test" {
  name = "fgt_sslvpn"
  # local
  realm_id = data.fortitokencloud_realm.test.id

  sp_entity_id = "https://<fgt_vpn_host>/remote/saml/metadata"
  sp_acs_url   = "https://<fgt_vpn_host>/remote/saml/login"
  sp_slo_url   = "https://<fgt_vpn_host>/remote/saml/logout"
  user_source_ids = [
    fortitokencloud_usersource.test.id
  ]
  ttl = 900
  attr_mapping = jsonencode({
    username = "Username"
  })
}

```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `realm_id` (String)

### Optional

- `attr_mapping` (String)
- `branding_id` (String)
- `signing_cert_id` (String)
- `sp_acs_url` (String)
- `sp_entity_id` (String)
- `sp_slo_url` (String)
- `ttl` (Number)
- `user_source_ids` (Set of String)

### Read-Only

- `entity_id` (String)
- `id` (String) The ID of this resource.
- `prefix` (String)
- `slo_url` (String)
- `sso_url` (String)
