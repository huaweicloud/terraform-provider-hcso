package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func getELBResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	eipID := state.Primary.Attributes["ipv4_eip_id"]
	eipType := state.Primary.Attributes["iptype"]
	if eipType != "" && eipID != "" {
		eipClient, err := c.NetworkingV1Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
		}

		if _, err := eips.Get(eipClient, eipID).Extract(); err != nil {
			return nil, err
		}
	}

	return loadbalancers.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttr(resourceName, "backend_subnets.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "backend_subnets.0",
						"hcso_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "true"),
					resource.TestCheckResourceAttr(resourceName, "backend_subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
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

func TestAccElbV3LoadBalancer_withEpsId(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_withEIP(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEIP(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_eip_id"),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_withEIP_Bandwidth_Id(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEIP_Bandwidth_Id(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_eip_id"),
				),
			},
		},
	})
}

func testAccElbV3LoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

resource "hcso_vpc_subnet" "test_1" {
  name       = "%[2]s_1"
  vpc_id     = hcso_vpc.test.id
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}

resource "hcso_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = hcso_vpc.test.id
  ipv4_subnet_id = hcso_vpc_subnet.test.ipv4_subnet_id
	
  availability_zone = [
    data.hcso_availability_zones.test.names[0]
  ]

  backend_subnets = [
    hcso_vpc_subnet.test.id
  ]

  protection_status = "nonProtection"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3LoadBalancerConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

resource "hcso_vpc_subnet" "test_1" {
  name       = "%[2]s_1"
  vpc_id     = hcso_vpc.test.id
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}

resource "hcso_elb_loadbalancer" "test" {
  name              = "%[3]s"
  cross_vpc_backend = true
  vpc_id            = hcso_vpc.test.id
  ipv4_subnet_id    = hcso_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.hcso_availability_zones.test.names[0]
  ]

  backend_subnets = [
    hcso_vpc_subnet.test.id,
    hcso_vpc_subnet.test_1.id,
  ]

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, common.TestVpc(rName), rName, rNameUpdate)
}

func testAccElbV3LoadBalancerConfig_withEpsId(rName string) string {
	return fmt.Sprintf(`
data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_availability_zones" "test" {}

resource "hcso_elb_loadbalancer" "test" {
  name                  = "%s"
  ipv4_subnet_id        = data.hcso_vpc_subnet.test.ipv4_subnet_id
  enterprise_project_id = "%s"

  availability_zone = [
    data.hcso_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccElbV3LoadBalancerConfig_withEIP(rName string) string {
	return fmt.Sprintf(`
data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_availability_zones" "test" {}

resource "hcso_elb_loadbalancer" "test" {
  name           = "%s"
  ipv4_subnet_id = data.hcso_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.hcso_availability_zones.test.names[0]
  ]

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 5
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_withEIP_Bandwidth_Id(rName string) string {
	return fmt.Sprintf(`
data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_availability_zones" "test" {}

resource "hcso_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "hcso_elb_loadbalancer" "test" {
  name           = "%[1]s"
  ipv4_subnet_id = data.hcso_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.hcso_availability_zones.test.names[0]
  ]

  iptype       = "5_bgp"
  sharetype    = "WHOLE"
  bandwidth_id = hcso_vpc_bandwidth.test.id
}
`, rName)
}
