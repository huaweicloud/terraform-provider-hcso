package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/members"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func getBackupShareResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CbrV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR v3 client: %s", err)
	}
	var (
		backupId = state.Primary.ID
		opts     = members.ListOpts{
			BackupId: backupId,
		}
	)
	memberList, err := members.List(client, opts)
	if len(memberList) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return memberList, err
}

// TODO Failed: POST API crate checkpoint not support resource_details but resources
func TestAccBackupShare_basic(t *testing.T) {
	var (
		memberList   []members.Member
		name         = acceptance.RandomAccResourceName()
		resourceName = "hcso_cbr_backup_share.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&memberList,
		getBackupShareResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDestProjectIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupShare_basic(name, acceptance.HCSO_DEST_PROJECT_ID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "backup_id"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.0.dest_project_id", acceptance.HCSO_DEST_PROJECT_ID),
				),
			},
			{
				Config: testAccBackupShare_basic(name, acceptance.HCSO_DEST_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "backup_id"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.0.dest_project_id", acceptance.HCSO_DEST_PROJECT_ID_TEST),
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

func testAccBackupShare_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.hcso_images_image.test.id
  flavor_id          = data.hcso_compute_flavors.test.ids[0]
  security_group_ids = [hcso_networking_secgroup.test.id]
  availability_zone  = data.hcso_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = hcso_vpc_subnet.test.id
  }
}

resource "hcso_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 10

  resources {
    server_id = hcso_compute_instance.test.id
  }
}

// TODO Failed: POST API crate checkpoint not support resource_details but resources
resource "hcso_cbr_checkpoint" "test" {
  vault_id = hcso_cbr_vault.test.id
  name     = "%[2]s"

  backups {
    type        = "OS::Nova::Server"
    resource_id = hcso_compute_instance.test.id
  }
}
`, common.TestBaseComputeResources(name), name)
}

func testAccBackupShare_basic(name, destProjectId string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_backup_share" "test" {
  // TODO Failed: POST API crate checkpoint not support resource_details but resources,Using existing backup_id
  backup_id = try(tolist(hcso_cbr_checkpoint.test.backups)[0].id, "")

  members {
    # Different user (ID) in the same region.
    dest_project_id = "%[2]s"
  }
}
`, testAccBackupShare_base(name), destProjectId)
}
