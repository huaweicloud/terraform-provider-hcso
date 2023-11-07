package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func TestAccDataCciNamespaces_basic(t *testing.T) {
	dataSourceName := "data.hcso_cci_namespaces.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCciNamespaces_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.auto_expend_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.container_network_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.rbac_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.status", "Active"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.type", "general-computing"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.recycling_interval"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.warmup_pool_size"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.name",
						"hcso_cci_network.test", "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.security_group_id",
						"hcso_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.vpc.0.id",
						"hcso_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.vpc.0.subnet_id",
						"hcso_vpc_subnet.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.vpc.0.subnet_cidr",
						"hcso_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.vpc.0.network_id",
						"hcso_vpc_subnet.test", "id"),
				),
			},
		},
	})
}

func TestAccDataCciNamespaces_noNetwork(t *testing.T) {
	dataSourceName := "data.hcso_cci_namespaces.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCciNamespaces_noNetwork(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.type", "general-computing"),
				),
			},
		},
	})
}

func testAccDataCciNamespaces_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

resource "hcso_cci_namespace" "test" {
  name = "%[2]s"
  type = "general-computing"
}

resource "hcso_cci_network" "test" {
  availability_zone = data.hcso_availability_zones.test.names[0]
  namespace         = hcso_cci_namespace.test.name
  name              = "%[2]s"
  network_id        = hcso_vpc_subnet.test.id
  security_group_id = hcso_networking_secgroup.test.id
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccDataCciNamespaces_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_cci_namespaces" "test" {
  depends_on = [hcso_cci_network.test]

  name = "%s"
  type = "general-computing"
}
`, testAccDataCciNamespaces_base(rName), rName)
}

func testAccDataCciNamespaces_noNetwork(rName string) string {
	return fmt.Sprintf(`
resource "hcso_cci_namespace" "test" {
  name = "%[1]s"
  type = "general-computing"
}

data "hcso_cci_namespaces" "test" {
  name = hcso_cci_namespace.test.name
}
`, rName)
}
