package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/bms/v1/baremetalservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func TestAccBmsInstance_basic(t *testing.T) {
	var instance baremetalservers.CloudServer

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcso_bms_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_basic(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccBmsInstance_basic(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func testAccCheckBmsInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	bmsClient, err := config.BmsV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud bms client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_bms_instance" {
			continue
		}

		server, err := baremetalservers.Get(bmsClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckBmsInstanceExists(n string, instance *baremetalservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		bmsClient, err := config.BmsV1Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud bms client: %s", err)
		}

		found, err := baremetalservers.Get(bmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccBmsInstance_base(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

data "hcso_bms_flavors" "test" {
  availability_zone = try(element(data.hcso_availability_zones.test.names, 0), "")
}

resource "hcso_kps_keypair" "test" {
  name = "%s"
}`, common.TestBaseNetwork(rName), rName)
}

func testAccBmsInstance_basic(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[2]s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcso_bms_instance" "test" {
  security_groups   = [hcso_networking_secgroup.test.id]
  availability_zone = data.hcso_availability_zones.test.names[0]
  vpc_id            = hcso_vpc.test.id
  flavor_id         = data.hcso_bms_flavors.test.flavors[0].id
  key_pair          = hcso_kps_keypair.test.name
  image_id          = "519ea918-1fea-4ebc-911a-593739b1a3bc" # CentOS 7.4 64bit for BareMetal

  name                  = "%[2]s"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = hcso_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "%[5]v"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HCSO_USER_ID, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST, isAutoRenew)
}
