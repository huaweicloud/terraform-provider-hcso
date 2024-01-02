---
subcategory: "Scalable File Service (SFS)"
---

# hcso_sfs_turbo

Provides a Shared File System (SFS) Turbo resource.

## Example Usage

### Create a STANDARD Shared File System (SFS) Turbo

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "test_az" {}

resource "hcso_sfs_turbo" "test" {
  name              = "sfs-turbo-standard-1"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.test_az

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a PERFORMANCE Enhanced Shared File System (SFS) Turbo

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "test_az" {}

resource "hcso_sfs_turbo" "test" {
  name              = "sfs-turbo-performance-1"
  size              = 500
  share_proto       = "NFS"
  share_type        = "PERFORMANCE"
  enhanced          = true
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.test_az

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SFS Turbo resource. If omitted, the
  provider-level region will be used. Changing this creates a new SFS Turbo resource.

* `name` - (Required, String) Specifies the name of an SFS Turbo file system. The value contains 4 to 64
  characters and must start with a letter.

* `share_proto` - (Optional, String, ForceNew) Specifies the protocol for sharing file systems. The valid value is NFS.
  Changing this will create a new resource.

* `share_type` - (Optional, String, ForceNew) Specifies the file system type. Changing this will create a new resource.
  Valid values are **STANDARD**, **PERFORMANCE**. Defaults to **STANDARD**.

* `size` - (Required, Int) Specifies the capacity of a sharing file system, in GB.
  + If `share_type` is set to **STANDARD** or **PERFORMANCE**, the value ranges from 500 to 32768, and ranges from
  10240 to 327680 for an enhanced file system.
  -> The file system capacity can only be expanded, not reduced.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone where the file system is located.
  Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of the subnet. Changing this will create a new
  resource.

* `security_group_id` - (Required, String) Specifies the security group ID.

* `enhanced` - (Optional, Bool, ForceNew) Specifies whether the file system is enhanced or not. Changing this will
  create a new resource. The default value is `false`. This parameter is valid only when `share_type` is set
  to **STANDARD** or **PERFORMANCE**.

* `crypt_key_id` - (Optional, String, ForceNew) Specifies the ID of a KMS key to encrypt the file system. Changing this
  will create a new resource.

* `dedicated_flavor` - (Optional, String, ForceNew) Specifies the VM flavor used for creating a dedicated file system.

* `dedicated_storage_id` - (Optional, String, ForceNew) Specifies the ID of the dedicated distributed storage used
  when creating a dedicated file system.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the file system. Changing this
  will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the SFS Turbo.

-> **NOTE:**
SFS Turbo will create two private IP addresses and one virtual IP address under the subnet you specified. To ensure
normal use, SFS Turbo will enable the inbound rules for ports *111*, *445*, *2049*, *2051*, *2052*, and *20048* in the
security group you specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the SFS Turbo file system.

* `region` - The region of the SFS Turbo file system.

* `status` - The status of the SFS Turbo file system.

* `version` - The version ID of the SFS Turbo file system.

* `export_location` - The mount point of the SFS Turbo file system.

* `available_capacity` - The available capacity of the SFS Turbo file system in the unit of GB.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 10 minutes.

## Import

SFS Turbo can be imported using the `id`, e.g.

```
$ terraform import hcso_sfs_turbo 1e3d5306-24c9-4316-9185-70e9787d71ab
```
