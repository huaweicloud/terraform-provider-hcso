package internal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIECNetworkACLDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIECNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hcso_iec_network_acl.by_name", "name", rName),
					resource.TestCheckResourceAttr(
						"data.hcso_iec_network_acl.by_id", "name", rName),
				),
			},
		},
	})
}

func testAccDataSourceIECNetworkACL_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcso_iec_network_acl" "test" {
  name        = "%s"
  description = "IEC network acl for acc test"
}

data "hcso_iec_network_acl" "by_name" {
  name = hcso_iec_network_acl.test.name
}

data "hcso_iec_network_acl" "by_id" {
  id = hcso_iec_network_acl.test.id
}
`, rName)
}
