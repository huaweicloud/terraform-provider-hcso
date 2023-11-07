package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/csbs/v1/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
)

func TestAccCSBSBackupPolicyV1_basic(t *testing.T) {
	var policy policies.BackupPolicy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("tf-acc-test-update-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("hcso_csbs_backup_policy.backup_policy", &policy),
					resource.TestCheckResourceAttr(
						"hcso_csbs_backup_policy.backup_policy", "name", rName),
					resource.TestCheckResourceAttr(
						"hcso_csbs_backup_policy.backup_policy", "status", "suspended"),
				),
			},
			{
				ResourceName:      "hcso_csbs_backup_policy.backup_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCSBSBackupPolicyV1_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("hcso_csbs_backup_policy.backup_policy", &policy),
					resource.TestCheckResourceAttr(
						"hcso_csbs_backup_policy.backup_policy", "name", updateName),
				),
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_timeout(t *testing.T) {
	var policy policies.BackupPolicy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_timeout(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("hcso_csbs_backup_policy.backup_policy", &policy),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	policyClient, err := config.CsbsV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating csbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_csbs_backup_policy" {
			continue
		}

		_, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("backup policy still exists")
		}
	}

	return nil
}

func testAccCheckCSBSBackupPolicyV1Exists(n string, policy *policies.BackupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		policyClient, err := config.CsbsV1Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating csbs client: %s", err)
		}

		found, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("backup policy not found")
		}

		*policy = *found

		return nil
	}
}

func testAccCSBSBackupPolicyV1_basic(rName string) string {
	return fmt.Sprintf(`
data "hcso_networking_secgroup" "test" {
  name = "default"
}

resource "hcso_compute_instance_v2" "instance_1" {
  name               = "%s"
  image_id           = "%s"
  security_group_ids = [data.hcso_networking_secgroup.test.id]
  availability_zone  = "%s"
  flavor_id          = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "hcso_csbs_backup_policy" "backup_policy" {
  name = "%s"
  resource {
    id   = hcso_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }
  scheduled_operation {
    name            ="mybackup"
    enabled         = true
    operation_type  ="backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}
`, rName, acceptance.HCSO_IMAGE_ID, acceptance.HCSO_AVAILABILITY_ZONE, acceptance.HCSO_FLAVOR_ID, acceptance.HCSO_NETWORK_ID, rName)
}

func testAccCSBSBackupPolicyV1_update(rName, updateName string) string {
	return fmt.Sprintf(`
data "hcso_networking_secgroup" "test" {
  name = "default"
}

resource "hcso_compute_instance_v2" "instance_1" {
  name               = "%s"
  image_id           = "%s"
  security_group_ids = [data.hcso_networking_secgroup.test.id]
  availability_zone  = "%s"
  flavor_id          = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "hcso_csbs_backup_policy" "backup_policy" {
  name = "%s"
  resource {
    id   = hcso_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }
  scheduled_operation {
    name            ="mybackup"
    enabled         = true
    operation_type  ="backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}
`, rName, acceptance.HCSO_IMAGE_ID, acceptance.HCSO_AVAILABILITY_ZONE, acceptance.HCSO_FLAVOR_ID, acceptance.HCSO_NETWORK_ID, updateName)
}

func testAccCSBSBackupPolicyV1_timeout(rName string) string {
	return fmt.Sprintf(`
data "hcso_networking_secgroup" "test" {
  name = "default"
}

resource "hcso_compute_instance_v2" "instance_1" {
  name               = "%s"
  image_id           = "%s"
  security_group_ids = [data.hcso_networking_secgroup.test.id]
  availability_zone  = "%s"
  flavor_id          = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "hcso_csbs_backup_policy" "backup_policy" {
  name = "%s"
  resource {
    id   = hcso_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }
  scheduled_operation {
    name            ="mybackup"
    enabled         = true
    operation_type  ="backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, rName, acceptance.HCSO_IMAGE_ID, acceptance.HCSO_AVAILABILITY_ZONE, acceptance.HCSO_FLAVOR_ID, acceptance.HCSO_NETWORK_ID, rName)
}
