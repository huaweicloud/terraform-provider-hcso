package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func getEipResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
	}
	return eips.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcEip_basic(t *testing.T) {
	var (
		eip eips.PublicIp

		resourceName = "hcso_vpc_eip.test"
		randName     = acceptance.RandomAccResourceName()
		udpateName   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "5"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.share_type", "PER"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.charge_mode", "traffic"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				Config: testAccVpcEip_update(udpateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", udpateName),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", udpateName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "8"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
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

func TestAccVpcEip_share(t *testing.T) {
	var (
		eip eips.PublicIp

		randName     = acceptance.RandomAccResourceName()
		resourceName = "hcso_vpc_eip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_share(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
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

func TestAccVpcEip_WithEpsId(t *testing.T) {
	var (
		eip eips.PublicIp

		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcso_vpc_eip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_epsId(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccVpcEip_deprecated(t *testing.T) {
	var (
		eip eips.PublicIp

		randName        = acceptance.RandomAccResourceName()
		resourceName    = "hcso_vpc_eip.test"
		vipResourceName = "hcso_networking_vip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_deprecated(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "BOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.ip_version", "4"),
					resource.TestCheckResourceAttrPair(resourceName, "private_ip", vipResourceName, "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", vipResourceName, "id"),
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

func testAccVpcEip_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vpc_eip" "test" {
  name = "%[1]s"

  publicip {
    type       = "5_bgp"
    ip_version = 4
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = 5
    charge_mode = "traffic"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccVpcEip_update(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vpc_eip" "test" {
  name = "%[1]s"

  publicip {
    type       = "5_bgp"
    #ip_version = 6 (update field 'ip_version' is unsupported by the API(PUT /v1/{project_id}/pubilcips/{ID})
    ip_version = 4
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = 8
    charge_mode = "traffic"
  }

  tags = {
    foo  = "bar1"
    key1 = "value"
  }
}
`, rName)
}

func testAccVpcEip_epsId(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vpc_eip" "test" {
  enterprise_project_id = "%[1]s"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[2]s"
    size        = 5
    charge_mode = "traffic"
  }
}
`, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST, rName)
}

func testAccVpcEip_share(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vpc_bandwidth" "test" {
  name = "%s"
  size = 5
}

resource "hcso_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = hcso_vpc_bandwidth.test.id
  }
}
`, rName)
}

func testAccVpcEip_deprecated(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcso_vpc_subnet" "test" {
  vpc_id     = hcso_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(hcso_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(hcso_vpc.test.cidr, 4, 1), 1)
}

resource "hcso_networking_vip" "test" {
  name       = "%[1]s"
  network_id = hcso_vpc_subnet.test.id
}

resource "hcso_vpc_eip" "test" {
  publicip {
    type    = "5_bgp"
    port_id = hcso_networking_vip.test.id
  }

  bandwidth {
    name        = "%[1]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, rName)
}
