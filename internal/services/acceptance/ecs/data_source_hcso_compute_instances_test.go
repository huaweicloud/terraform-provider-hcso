package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccComputeInstancesDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.hcso_compute_instances.test"
	var instance cloudservers.CloudServer

	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rc := acceptance.InitResourceCheck(
		"hcso_compute_instance.test",
		&instance,
		getEcsInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.network.#"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr("data.hcso_compute_instances.byID", "instances.#", "1"),
				),
			},
		},
	})
}

func testAccComputeInstancesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_compute_instance" "test" {
  name               = "%s"
  image_id           = data.hcso_images_image.test.id
  flavor_id          = data.hcso_compute_flavors.test.ids[0]
  security_group_ids = [data.hcso_networking_secgroups.test.security_groups[0].id]
  availability_zone  = data.hcso_availability_zones.test.names[0]
  system_disk_type   = "SSD"
  network {
    uuid = data.hcso_vpc_subnets.test.subnets[0].id
  }

  tags = {
    foo = "bar"
  }
}

data "hcso_compute_instances" "test" {
  name = hcso_compute_instance.test.name
}

data "hcso_compute_instances" "byID" {
  instance_id = hcso_compute_instance.test.id
}
`, testAccCompute_data, rName)
}
