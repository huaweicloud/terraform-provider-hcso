---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcso_compute_volume_attach

Attaches a volume to an ECS Instance.

## Example Usage

### Basic attachment of a single volume to a single instance

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

resource "hcso_evs_volume" "myvol" {
  name              = "volume"
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  volume_type       = "SSD"
  size              = 10
}

resource "hcso_compute_volume_attach" "attached" {
  instance_id = hcso_compute_instance.myinstance.id
  volume_id   = hcso_evs_volume.myvol.id
}
```

### Attaching multiple volumes to a single instance

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

resource "hcso_evs_volume" "myvol" {
  count             = 2
  name              = "volume"
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  volume_type       = "SSD"
  size              = 10
}

resource "hcso_compute_volume_attach" "attachments" {
  count       = 2
  instance_id = hcso_compute_instance.myinstance.id
  volume_id   = element(hcso_evs_volume.myvol[*].id, count.index)
}

output "volume_devices" {
  value = hcso_compute_volume_attach.attachments[*].device
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the volume resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the Instance to attach the Volume to.

* `volume_id` - (Required, String, ForceNew) Specifies the ID of the Volume to attach to an Instance.

* `device` - (Optional, String) Specifies the device of the volume attachment (ex: `/dev/vdc`).

  -> Being able to specify a device is dependent upon the hypervisor in use. There is a chance that the device
  specified in Terraform will not be the same device the hypervisor chose. If this happens, Terraform will wish to
  update the device upon subsequent applying which will cause the volume to be detached and reattached indefinitely.
  Please use with caution.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `pci_address` - PCI address of the block device.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Volume Attachments can be imported using the Instance ID and Volume ID separated by a slash, e.g.

```shell
$ terraform import hcso_compute_volume_attach.va_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
