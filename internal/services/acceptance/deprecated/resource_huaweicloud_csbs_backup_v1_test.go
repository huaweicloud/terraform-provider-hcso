package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/csbs/v1/backup"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
)

func TestAccCSBSBackupV1_basic(t *testing.T) {
	var backups backup.Backup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("hcso_csbs_backup.csbs", &backups),
					resource.TestCheckResourceAttr(
						"hcso_csbs_backup.csbs", "backup_name", rName),
					resource.TestCheckResourceAttr(
						"hcso_csbs_backup.csbs", "resource_type", "OS::Nova::Server"),
				),
			},
			{
				ResourceName:      "hcso_csbs_backup.csbs",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccCSBSBackupV1_timeout(t *testing.T) {
	var backups backup.Backup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_timeout(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("hcso_csbs_backup.csbs", &backups),
				),
			},
		},
	})
}

func testAccCSBSBackupV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	backupClient, err := config.CsbsV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating csbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_csbs_backup" {
			continue
		}

		_, err := backup.Get(backupClient, rs.Primary.ID).ExtractBackup()
		if err == nil {
			return fmtp.Errorf("Backup still exists")
		}
	}

	return nil
}

func testAccCSBSBackupV1Exists(n string, backups *backup.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Backup not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		backupClient, err := config.CsbsV1Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating csbs client: %s", err)
		}

		found, err := backup.Get(backupClient, rs.Primary.ID).ExtractBackup()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmtp.Errorf("backup not found")
		}

		*backups = *found

		return nil
	}
}

func testAccCSBSBackupV1_basic(rName string) string {
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
resource "hcso_csbs_backup" "csbs" {
  backup_name      = "%s"
  description      = "test-code"
  resource_id      = hcso_compute_instance_v2.instance_1.id
  resource_type    = "OS::Nova::Server"
}
`, rName, acceptance.HCSO_IMAGE_ID, acceptance.HCSO_AVAILABILITY_ZONE, acceptance.HCSO_FLAVOR_ID, acceptance.HCSO_NETWORK_ID, rName)
}

func testAccCSBSBackupV1_timeout(rName string) string {
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
resource "hcso_csbs_backup" "csbs" {
  backup_name      = "%s"
  description      = "test-code"
  resource_id      = hcso_compute_instance_v2.instance_1.id
  resource_type    = "OS::Nova::Server"
}
`, rName, acceptance.HCSO_IMAGE_ID, acceptance.HCSO_AVAILABILITY_ZONE, acceptance.HCSO_FLAVOR_ID, acceptance.HCSO_NETWORK_ID, rName)
}
