---
subcategory: "Image Management Service (IMS)"
---

# hcso_images_image_share

Manages an IMS image share resource within Huawei Cloud Stack Online.

## Example Usage

```hcl
variable "source_image_id" {}
variable "target_project_ids" {
  type = set(string)
}

resource "hcso_images_image_share" "test" {
  source_image_id    = var.source_image_id
  target_project_ids = var.target_project_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `source_image_id` - (Required, String, ForceNew) Specifies the ID of the source image.

  Changing this parameter will create a new resource.

* `target_project_ids` - (Required, List) Specifies the IDs of the target projects.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.
