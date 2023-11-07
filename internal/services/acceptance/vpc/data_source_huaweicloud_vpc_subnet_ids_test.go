package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccVpcSubnetIdsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcso_vpc_subnet_ids.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetIdsDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "2"),
				),
			},
		},
	})
}

func testAccSubnetIdsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

resource "hcso_vpc" "test" {
  name = "%s"
  cidr = "172.16.128.0/20"
}

resource "hcso_vpc_subnet" "test1" {
  name       = "%s"
  cidr       = "172.16.140.0/22"
  gateway_ip = "172.16.140.1"
  vpc_id     = hcso_vpc.test.id
}

resource "hcso_vpc_subnet" "test2" {
  name       = "%s"
  cidr       = "172.16.136.0/22"
  gateway_ip = "172.16.136.1"
  vpc_id     = hcso_vpc.test.id
}

data "hcso_vpc_subnet_ids" "test" {
  vpc_id = hcso_vpc_subnet.test1.vpc_id
}
`, rName, rName, rName)
}
