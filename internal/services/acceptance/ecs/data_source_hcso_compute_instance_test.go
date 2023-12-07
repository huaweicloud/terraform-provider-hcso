package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func getEcsInstanceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ComputeV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating compute v1 client: %s", err)
	}

	resourceID := state.Primary.ID
	found, err := cloudservers.Get(client, resourceID).Extract()
	if err == nil && found.Status == "DELETED" {
		return nil, fmt.Errorf("the resource %s has been deleted", resourceID)
	}

	return found, err
}

func TestAccComputeInstanceDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.hcso_compute_instance.this"
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
				Config: testAccComputeInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume_attached.#"),
					resource.TestCheckResourceAttrSet("data.hcso_compute_instance.byID", "status"),
				),
			},
		},
	})
}

func testAccComputeInstanceDataSource_basic(rName string) string {
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
}

data "hcso_compute_instance" "this" {
  name = hcso_compute_instance.test.name
}

data "hcso_compute_instance" "byID" {
  instance_id = hcso_compute_instance.test.id
}
`, testAccCompute_data, rName)
}
