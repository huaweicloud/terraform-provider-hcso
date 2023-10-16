package cce

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCCENodePoolV3DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.hcso_cce_node_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckCCENodePoolV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find node pools data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Node pool data source ID not set ")
		}

		return nil
	}
}

func testAccCCENodePoolV3DataSource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cce_node_pool" "test" {
  cluster_id               = hcso_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = data.hcso_compute_flavors.test.ids[0]
  initial_node_count       = 1
  availability_zone        = data.hcso_availability_zones.test.names[0]
  key_pair                 = hcso_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}

data "hcso_cce_node_pool" "test" {
  cluster_id = hcso_cce_cluster.test.id
  name       = hcso_cce_node_pool.test.name
}
`, testAccNodePool_base(name), name)
}
