package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/routes"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getRouteV2ResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud Network client: %s", err)
	}
	return routes.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcRouteV2_basic(t *testing.T) {
	var route routes.Route

	randName := acceptance.RandomAccResourceName()
	resourceName := "hcso_vpc_route_v2.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&route,
		getRouteV2ResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteV2_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "peering"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "nexthop",
						"${hcso_vpc_peering_connection.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "destination",
						"${hcso_vpc.test2.cidr}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "vpc_id",
						"${hcso_vpc.test1.id}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRouteV2_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vpc" "test1" {
  name = "%s_1"
  cidr = "172.16.0.0/20"
}

resource "hcso_vpc" "test2" {
  name = "%s_2"
  cidr = "172.16.128.0/20"
}

resource "hcso_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = hcso_vpc.test1.id
  peer_vpc_id = hcso_vpc.test2.id
}

resource "hcso_vpc_route_v2" "test" {
  type        = "peering"
  nexthop     = hcso_vpc_peering_connection.test.id
  destination = hcso_vpc.test2.cidr
  vpc_id      = hcso_vpc.test1.id
}
`, rName, rName, rName)
}
