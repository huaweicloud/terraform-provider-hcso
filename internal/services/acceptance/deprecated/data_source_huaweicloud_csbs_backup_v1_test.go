package deprecated

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCSBSBackupV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupV1DataSourceID("data.hcso_csbs_backup_v1.csbs"),
					resource.TestCheckResourceAttr("data.hcso_csbs_backup_v1.csbs", "backup_name", "csbs-test"),
					resource.TestCheckResourceAttr("data.hcso_csbs_backup_v1.csbs", "resource_name", "instance_1"),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find backup data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("backup data source ID not set ")
		}

		return nil
	}
}

var testAccCSBSBackupV1DataSource_basic = fmt.Sprintf(`
resource "hcso_compute_instance_v2" "instance_1" {
  name = "instance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "hcso_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test"
  description      = "test-code"
  resource_id = "${hcso_compute_instance_v2.instance_1.id}"
  resource_type = "OS::Nova::Server"
}
data "hcso_csbs_backup_v1" "csbs" {
  id = "${hcso_csbs_backup_v1.csbs.id}"
}
`, acceptance.HCSO_IMAGE_ID, acceptance.HCSO_AVAILABILITY_ZONE, acceptance.HCSO_FLAVOR_ID, acceptance.HCSO_NETWORK_ID)
