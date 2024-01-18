---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcso_compute_interface_attach

Attaches a Network Interface to an Instance.

## Example Usage

### Attach a port (under the specified network) to the ECS instance and generate a random IP address

```hcl
variable "instance_id" {}
variable "network_id" {}

resource "hcso_compute_interface_attach" "test" {
  instance_id = var.instance_id
  network_id  = var.network_id
}
```

### Attach a port (under the specified network) to the ECS instance and use the custom security groups

```hcl
variable "instance_id" {}
variable "network_id" {}
variable "security_group_ids" {
  type = list(string)
}

resource "hcso_compute_interface_attach" "test" {
  instance_id        = var.instance_id
  network_id         = var.network_id
  fixed_ip           = "192.168.10.199"
  security_group_ids = var.security_group_ids
}
```

### Attach a custom port to the ECS instance

```hcl
variable "secgroup_id" {}

data "hcso_availability_zones" "myaz" {}

data "hcso_compute_flavors" "myflavor" {
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcso_vpc_subnets" "mynet" {
  name = "subnet-default"
}

data "hcso_images_image" "myimage" {
  name_regex  = "^Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "hcso_compute_keypair" "my-keypair" {
  name = "my-keypair"
}

data "hcso_networking_port" "myport" {
  network_id = data.hcso_vpc_subnets.mynet.subnets[0].id
  fixed_ip   = "192.168.0.100"
}

resource "hcso_compute_instance" "myinstance" {
  name               = "myinstance"
  image_id           = data.hcso_images_image.myimage.id
  flavor_id          = data.hcso_compute_flavors.myflavor.ids[0]
  security_group_ids = [var.secgroup_id]
  availability_zone  = data.hcso_availability_zones.myaz.names[0]
  key_pair           = hcso_compute_keypair.my-keypair.name
  system_disk_type   = "SSD"
  system_disk_size   = 40

  network {
    uuid = data.hcso_vpc_subnets.mynet.subnets[0].id
  }
}

resource "hcso_compute_interface_attach" "attached" {
  instance_id = hcso_compute_instance.myinstance.id
  port_id     = data.hcso_networking_port.myport.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the network interface attache resource. If
  omitted, the provider-level region will be used. Changing this creates a new network interface attache resource.

* `instance_id` - (Required, String, ForceNew) The ID of the Instance to attach the Port or Network to.

* `port_id` - (Optional, String, ForceNew) The ID of the Port to attach to an Instance.
  This option and `network_id` are mutually exclusive.

* `network_id` - (Optional, String, ForceNew) The ID of the Network to attach to an Instance. A port will be created
  automatically.
  This option and `port_id` are mutually exclusive.

* `fixed_ip` - (Optional, String, ForceNew) An IP address to assosciate with the port.

  ->This option cannot be used with port_id. You must specify a network_id. The IP address must lie in a range on
  the supplied network.

* `source_dest_check` - (Optional, Bool) Specifies whether the ECS processes only traffic that is destined specifically
  for it. This function is enabled by default but should be disabled if the ECS functions as a SNAT server or has a
  virtual IP address bound to it.

* `security_group_ids` - (Optional, List) Specifies the list of security group IDs bound to the specified port.  
  Defaults to the default security group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of ECS instance ID and port ID separated by a slash.
* `mac` - The MAC address of the NIC.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Interface Attachments can be imported using the Instance ID and Port ID separated by a slash, e.g.

```shell
$ terraform import hcso_compute_interface_attach.ai_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
