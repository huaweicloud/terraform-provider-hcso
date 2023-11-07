package eip

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func TestAccEIPAssociate_basic(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "hcso_vpc_eip_associate.test"
	resourceName := "hcso_vpc_eip.test"
	partten := `^((25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))\.){3}(25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))$`

	// hcso_vpc_eip_associate and hcso_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(associateName, "status", "BOUND"),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
					resource.TestMatchOutput("public_ip_address", regexp.MustCompile(partten)),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEIPAssociate_port(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "hcso_vpc_eip_associate.test"
	resourceName := "hcso_vpc_eip.test"

	// hcso_vpc_eip_associate and hcso_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_port(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(associateName, "status", "BOUND"),
					resource.TestCheckResourceAttrPtr(
						associateName, "port_id", &eip.PortID),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEIPAssociate_compatible(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "hcso_networking_eip_associate.test"
	resourceName := "hcso_vpc_eip.test"

	// hcso_networking_eip_associate and hcso_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_compatible(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(
						associateName, "port_id", &eip.PortID),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEIPAssociate_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%s"
    charge_mode = "traffic"
  }
}`, common.TestVpc(rName), rName)
}

func testAccEIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

data "hcso_compute_flavors" "test" {
  availability_zone = data.hcso_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

data "hcso_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "hcso_networking_secgroup" "test" {
  name                 = "%[2]s"
  delete_default_rules = true
}

resource "hcso_kps_keypair" "test" {
  name = "%[2]s"
}

resource "hcso_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.hcso_images_image.test.id
  flavor_id          = data.hcso_compute_flavors.test.ids[0]
  availability_zone  = data.hcso_availability_zones.test.names[0]
  security_group_ids = [hcso_networking_secgroup.test.id]

  key_pair = hcso_kps_keypair.test.name

  network {
    uuid = hcso_vpc_subnet.test.id
  }
}

resource "hcso_vpc_eip_associate" "test" {
  public_ip  = hcso_vpc_eip.test.address
  network_id = hcso_compute_instance.test.network[0].uuid
  fixed_ip   = hcso_compute_instance.test.network[0].fixed_ip_v4
}

data "hcso_compute_instance" "test" {
  depends_on = [hcso_vpc_eip_associate.test]

  name = "%[2]s"
}

output "public_ip_address" {
  value = data.hcso_compute_instance.test.public_ip
}
`, testAccEIPAssociate_base(rName), rName)
}

func testAccEIPAssociate_port(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_networking_vip" "test" {
  name       = "%s"
  network_id = hcso_vpc_subnet.test.id
}

resource "hcso_vpc_eip_associate" "test" {
  public_ip = hcso_vpc_eip.test.address
  port_id   = hcso_networking_vip.test.id
}
`, testAccEIPAssociate_base(rName), rName)
}

func testAccEIPAssociate_compatible(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_networking_vip" "test" {
  name       = "%s"
  network_id = hcso_vpc_subnet.test.id
}
  
resource "hcso_networking_eip_associate" "test" {
  public_ip = hcso_vpc_eip.test.address
  port_id   = hcso_networking_vip.test.id
}
`, testAccEIPAssociate_base(rName), rName)
}
