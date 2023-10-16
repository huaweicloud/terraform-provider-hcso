package deprecated

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
)

func TestAccVBSBackupPolicyV2DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupPolicyV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupPolicyV2DataSource("data.hcso_vbs_backup_policy_v2.policies"),
					resource.TestCheckResourceAttr("data.hcso_vbs_backup_policy_v2.policies", "name", "policy_001"),
					resource.TestCheckResourceAttr("data.hcso_vbs_backup_policy_v2.policies", "status", "ON"),
				),
			},
		},
	})
}

func testAccCheckVBSBackupPolicyV2DataSource(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find backup policy data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("backup policy ID not set ")
		}

		return nil
	}
}

var testAccVBSBackupPolicyV2DataSource_basic = `
resource "hcso_vbs_backup_policy_v2" "vbs_1" {
  name = "policy_001"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1
      tags =[
        {
          key = "k2"
          value = "v2"
          }]
}
data "hcso_vbs_backup_policy_v2" "policies" {
  id = "${hcso_vbs_backup_policy_v2.vbs_1.id}"
}
`
