package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/vbs/v2/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
)

func TestAccVBSBackupPolicyV2_basic(t *testing.T) {
	var policy policies.Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("tf-acc-test-update-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupPolicyV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("hcso_vbs_backup_policy.vbs", &policy),
					resource.TestCheckResourceAttr(
						"hcso_vbs_backup_policy.vbs", "name", rName),
					resource.TestCheckResourceAttr(
						"hcso_vbs_backup_policy.vbs", "status", "ON"),
				),
			},
			{
				ResourceName:      "hcso_vbs_backup_policy.vbs",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVBSBackupPolicyV2_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("hcso_vbs_backup_policy.vbs", &policy),
					resource.TestCheckResourceAttr(
						"hcso_vbs_backup_policy.vbs", "name", updateName),
					resource.TestCheckResourceAttr(
						"hcso_vbs_backup_policy.vbs", "status", "ON"),
					resource.TestCheckResourceAttr(
						"hcso_vbs_backup_policy.vbs", "rentention_num", "7"),
				),
			},
		},
	})
}

func TestAccVBSBackupPolicyV2_rentention_day(t *testing.T) {
	var policy policies.Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupPolicyV2_rentention_day(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("hcso_vbs_backup_policy.vbs", &policy),
					resource.TestCheckResourceAttr(
						"hcso_vbs_backup_policy.vbs", "name", rName),
					resource.TestCheckResourceAttr(
						"hcso_vbs_backup_policy.vbs", "status", "ON"),
					resource.TestCheckResourceAttr(
						"hcso_vbs_backup_policy.vbs", "rentention_day", "30"),
				),
			},
		},
	})
}

func testAccVBSBackupPolicyV2Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	vbsClient, err := config.VbsV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating internal sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_vbs_backup_policy" {
			continue
		}

		_, err := policies.List(vbsClient, policies.ListOpts{ID: rs.Primary.ID})
		if err != nil {
			return fmtp.Errorf("Backup Policy still exists")
		}
	}

	return nil
}

func testAccVBSBackupPolicyV2Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		vbsClient, err := config.VbsV2Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating internal vbs client: %s", err)
		}

		policyList, err := policies.List(vbsClient, policies.ListOpts{ID: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := policyList[0]
		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("backup policy not found")
		}

		*policy = found

		return nil
	}
}

func testAccVBSBackupPolicyV2_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vbs_backup_policy" "vbs" {
  name                = "%s"
  start_time          = "12:00"
  status              = "ON"
  retain_first_backup = "N"
  rentention_num      = 2
  frequency           = 1
  tags {
    key   = "k2"
    value = "v2"
  }
}
`, rName)
}

func testAccVBSBackupPolicyV2_update(updateName string) string {
	return fmt.Sprintf(`
resource "hcso_vbs_backup_policy" "vbs" {
  name                = "%s"
  start_time          = "12:00"
  status              = "ON"
  retain_first_backup = "N"
  rentention_num      = 7
  frequency           = 1
  tags {
    key   = "k2"
    value = "v2"
  }
}
`, updateName)
}

func testAccVBSBackupPolicyV2_rentention_day(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vbs_backup_policy" "vbs" {
  name                = "%s"
  start_time          = "00:00,12:00"
  retain_first_backup = "N"
  rentention_day      = 30
  week_frequency      = ["MON", "WED", "SAT"]
}`, rName)
}
