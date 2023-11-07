package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccVpcRouteDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcso_vpc_route.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckDeprecated(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "type", "peering"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${hcso_vpc.test1.id}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "destination",
						"${hcso_vpc.test2.cidr}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "nexthop",
						"${hcso_vpc_peering_connection.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcRouteDataSource_byVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcso_vpc_route.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteDataSource_byVpcId(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "type", "peering"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${hcso_vpc.test1.id}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "destination",
						"${hcso_vpc.test2.cidr}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "nexthop",
						"${hcso_vpc_peering_connection.test.id}"),
				),
			},
		},
	})
}

func testAccRouteDataSource_base(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vpc" "test1" {
  name = "%s_1"
  cidr = "192.168.128.0/20"
}

resource "hcso_vpc" "test2" {
  name = "%s_2"
  cidr = "192.168.192.0/20"
}

resource "hcso_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = hcso_vpc.test1.id
  peer_vpc_id = hcso_vpc.test2.id
}

resource "hcso_vpc_route" "test" {
  type        = "peering"
  nexthop     = hcso_vpc_peering_connection.test.id
  destination = hcso_vpc.test2.cidr
  vpc_id      = hcso_vpc.test1.id
}
`, rName, rName, rName)
}

func testAccRouteDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_vpc_route" "test" {
  id = hcso_vpc_route.test.id
}
`, testAccRouteDataSource_base(rName))
}

func testAccRouteDataSource_byVpcId(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_vpc_route" "test" {
  vpc_id = hcso_vpc_route.test.vpc_id
}
`, testAccRouteDataSource_base(rName))
}
