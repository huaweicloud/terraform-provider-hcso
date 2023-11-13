---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcso_vpc_peering_connection_accepter

Provides a resource to manage the accepter's side of a VPC Peering Connection.

-> **NOTE:** When a cross-tenant (requester's tenant differs from the accepter's tenant) VPC Peering Connection
  is created, a VPC Peering Connection resource is automatically created in the accepter's account.
  The requester can use the `hcso_vpc_peering_connection` resource to manage its side of the connection and
  the accepter can use the `hcso_vpc_peering_connection_accepter` resource to accept its side of the connection
  into management.

## Example Usage

```hcl
provider "hcso" {
  alias = "main"
}

provider "hcso" {
  alias = "peer"
}

resource "hcso_vpc" "vpc_main" {
  provider = "hcso.main"
  name     = var.vpc_name
  cidr     = var.vpc_cidr
}

resource "hcso_vpc" "vpc_peer" {
  provider = "hcso.peer"
  name     = var.peer_vpc_name
  cidr     = var.peer_vpc_cidr
}

# Requester's side of the connection.
resource "hcso_vpc_peering_connection" "peering" {
  provider       = "hcso.main"
  name           = var.peer_name
  vpc_id         = hcso_vpc.vpc_main.id
  peer_vpc_id    = hcso_vpc.vpc_peer.id
  peer_tenant_id = var.tenant_id
}

# Accepter's side of the connection.
resource "hcso_vpc_peering_connection_accepter" "peer" {
  provider = "hcso.peer"
  accept   = true

  vpc_peering_connection_id = hcso_vpc_peering_connection.peering.id
}
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the vpc peering connection accepter. If omitted,
  the provider-level region will be used. Changing this creates a new VPC peering connection accepter resource.

* `vpc_peering_connection_id` - (Required, String, ForceNew) The VPC Peering Connection ID to manage. Changing this
  creates a new VPC peering connection accepter.

* `accept` - (Optional, Bool) Whether or not to accept the peering request. Defaults to `false`.

## Removing hcso_vpc_peering_connection_accepter from your configuration

Huawei Cloud Stack Online allows a cross-tenant VPC Peering Connection to be deleted from either the requester's or
accepter's side. However, Terraform only allows the VPC Peering Connection to be deleted from the requester's side
by removing the corresponding `hcso_vpc_peering_connection` resource from your configuration.
Removing a `hcso_vpc_peering_connection_accepter` resource from your configuration will remove it from your
state file and management, but will not destroy the VPC Peering Connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VPC peering connection ID.

* `name` - The VPC peering connection name.

* `status` - The VPC peering connection status.

* `description` - The description of the VPC peering connection.

* `vpc_id` - The ID of requester VPC involved in a VPC peering connection.

* `peer_vpc_id` - The VPC ID of the accepter tenant.

* `peer_tenant_id` - The Tenant Id of the accepter tenant.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
