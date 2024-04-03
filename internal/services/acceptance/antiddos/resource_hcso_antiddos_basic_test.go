package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	antiddossdk "github.com/chnsz/golangsdk/openstack/antiddos/v1/antiddos"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

// TODO Failed by the OpenAPI of service Anti-DDoS not opened.
func TestAccAntiDdos_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcso_antiddos_basic.antiddos_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAntiDdosDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdos_config(rName, 200),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAntiDdosExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_threshold", "200"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "hcso_vpc_eip.eip_1", "address"),
				),
			},
			{
				Config: testAccAntiDdos_config(rName, 300),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "traffic_threshold", "300"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "hcso_vpc_eip.eip_1", "address"),
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

func testAccCheckAntiDdosDestroy(s *terraform.State) error {
	// the cloud native AntiDdos always exists
	return nil
}

func testAccCheckAntiDdosExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := cfg.AntiDDosV1Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating antiddos client: %s", err)
		}

		_, err = antiddossdk.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("error retrieving cloud native AntiDdos: %s", err)
		}

		return nil
	}
}

func testAccAntiDdos_config(rName string, threshold int) string {
	return fmt.Sprintf(`
resource "hcso_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "%s"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "hcso_antiddos_basic" "antiddos_1" {
  eip_id            = hcso_vpc_eip.eip_1.id
  traffic_threshold = %d
}
`, rName, threshold)
}
