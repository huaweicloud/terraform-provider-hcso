package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/nat/v3/snats"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getPrivateSnatRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 client: %s", err)
	}

	return snats.Get(client, state.Primary.ID)
}

func TestAccPrivateSnatRule_basic(t *testing.T) {
	var (
		obj snats.Rule

		rName = "hcso_nat_private_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateSnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gateway_id",
						"hcso_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id",
						"hcso_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"hcso_vpc_subnet.test", "id"),
				),
			},
			{
				Config: testAccPrivateSnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id",
						"hcso_nat_private_transit_ip.standby", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPrivateSnatRule_transitIpConfig(name string) string {
	return fmt.Sprintf(`
resource "hcso_vpc" "transit_ip_used" {
  name = "%[1]s-transit-ip"
  cidr = "172.16.0.0/16"
}

resource "hcso_vpc_subnet" "transit_ip_used" {
  vpc_id     = hcso_vpc.transit_ip_used.id
  name       = "%[1]s-transit-ip"
  cidr       = cidrsubnet(hcso_vpc.transit_ip_used.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(hcso_vpc.transit_ip_used.cidr, 4, 1), 1)
}

resource "hcso_nat_private_transit_ip" "test" {
  subnet_id             = hcso_vpc_subnet.transit_ip_used.id
  enterprise_project_id = "0"
}

resource "hcso_nat_private_transit_ip" "standby" {
  subnet_id             = hcso_vpc_subnet.transit_ip_used.id
  enterprise_project_id = "0"
}
`, name)
}

func testAccPrivateSnatRule_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "hcso_nat_private_gateway" "test" {
  subnet_id             = hcso_vpc_subnet.test.id
  name                  = "%[3]s"
  enterprise_project_id = "0"
}
`, common.TestBaseNetwork(name), testAccPrivateSnatRule_transitIpConfig(name), name)
}

func testAccPrivateSnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_nat_private_snat_rule" "test" {
  gateway_id    = hcso_nat_private_gateway.test.id
  description   = "Created by acc test"
  transit_ip_id = hcso_nat_private_transit_ip.test.id
  subnet_id     = hcso_vpc_subnet.test.id
}
`, testAccPrivateSnatRule_base(name))
}

func testAccPrivateSnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_nat_private_snat_rule" "test" {
  gateway_id    = hcso_nat_private_gateway.test.id
  transit_ip_id = hcso_nat_private_transit_ip.standby.id
  subnet_id     = hcso_vpc_subnet.test.id
}
`, testAccPrivateSnatRule_base(name))
}

func TestAccPrivateSnatRule_cidr(t *testing.T) {
	var (
		obj snats.Rule

		rName = "hcso_nat_private_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateSnatRule_cidr_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gateway_id",
						"hcso_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id",
						"hcso_nat_private_transit_ip.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "cidr",
						"hcso_vpc_subnet.test", "cidr"),
				),
			},
			{
				Config: testAccPrivateSnatRule_cidr_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id",
						"hcso_nat_private_transit_ip.standby", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPrivateSnatRule_cidr_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_nat_private_snat_rule" "test" {
  gateway_id    = hcso_nat_private_gateway.test.id
  description   = "Created by acc test"
  transit_ip_id = hcso_nat_private_transit_ip.test.id
  cidr          = hcso_vpc_subnet.test.cidr
}
`, testAccPrivateSnatRule_base(name))
}

func testAccPrivateSnatRule_cidr_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_nat_private_snat_rule" "test" {
  gateway_id    = hcso_nat_private_gateway.test.id
  transit_ip_id = hcso_nat_private_transit_ip.standby.id
  cidr          = hcso_vpc_subnet.test.cidr
}
`, testAccPrivateSnatRule_base(name))
}
