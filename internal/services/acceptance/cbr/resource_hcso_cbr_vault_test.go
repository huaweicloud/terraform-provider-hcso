package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func getVaultResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CbrV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR v3 client: %s", err)
	}
	return vaults.Get(c, state.Primary.ID).Extract()
}

func TestAccVault_backupServer(t *testing.T) {
	var (
		vault vaults.Vault

		resourceName = "hcso_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		basicConfig  = testAccVault_base(name)
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupServer_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "backup_name_prefix", "test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_az", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					// TODO excludes caused error: exclude_volumes cannot in extra_info.
					// resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
				),
			},
			{
				Config: testAccVault_backupServer_step2(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "backup_name_prefix", "test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_az", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					// TODO excludes caused error: exclude_volumes cannot in extra_info.
					// resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
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

func testAccVault_base(name string) string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    volume_type = string
    size        = number
    device_type = string
  }))
  default = [
    // TODO only "SSD" volume_type is valid.
    {volume_type = "SSD", size = 50, device_type = "VBD"},
    {volume_type = "SSD", size = 100, device_type = "VBD"},
    {volume_type = "SSD", size = 50, device_type = "VBD"},
  ]
}

// base compute resources
%[1]s

resource "hcso_compute_instance" "test" {
  availability_zone = data.hcso_availability_zones.test.names[0]
  name              = "%[2]s"
  image_id          = data.hcso_images_image.test.id
  flavor_id         = data.hcso_compute_flavors.test.ids[0]

  system_disk_type = "SSD"
  system_disk_size = 50

  security_group_ids = [
    hcso_networking_secgroup.test.id
  ]

  network {
    uuid = hcso_vpc_subnet.test.id
  }
}

resource "hcso_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.hcso_availability_zones.test.names[0]
  volume_type       = var.volume_configuration[count.index].volume_type
  name              = format("%[2]s_%%d", count.index)
  size              = var.volume_configuration[count.index].size
  device_type       = var.volume_configuration[count.index].device_type
}

resource "hcso_compute_volume_attach" "test" {
  count = length(var.volume_configuration)

  instance_id = hcso_compute_instance.test.id
  volume_id   = hcso_evs_volume.test[count.index].id
}`, common.TestBaseComputeResources(name), name)
}

func testAccVault_backupServer_step1(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "0"
  backup_name_prefix    = "test-prefix-"
  is_multi_az           = true

  resources {
    server_id = hcso_compute_instance.test.id
    // TODO excludes caused error: exclude_volumes cannot in extra_info.
    #excludes  = slice(hcso_compute_volume_attach.test[*].volume_id, 0, 2)
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, basicConfig, name)
}

func testAccVault_backupServer_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 300
  enterprise_project_id = "0"
  backup_name_prefix    = "test-prefix-"
  is_multi_az           = true

  resources {
    server_id = hcso_compute_instance.test.id
    // TODO excludes caused error: exclude_volumes cannot in extra_info.
    # excludes  = slice(hcso_compute_volume_attach.test[*].volume_id, 1, 3)
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, basicConfig, name)
}

func TestAccVault_volume(t *testing.T) {
	var (
		vault vaults.Vault

		resourceName = "hcso_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		basicConfig  = testAccVault_base(name)
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_volume_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "50"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
				),
			},
			{
				Config: testAccVault_volume_step2(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "true"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
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

func testAccVault_volume_step1(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 50
  enterprise_project_id = "0"

  resources {
    includes = slice(hcso_compute_volume_attach.test[*].volume_id, 0, 2)
  }
}
`, basicConfig, name)
}

func testAccVault_volume_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 100
  auto_expand           = true
  enterprise_project_id = "0"

  resources {
    includes = slice(hcso_compute_volume_attach.test[*].volume_id, 1, 3)
  }
}
`, basicConfig, name)
}

func TestAccVault_backupTurbo(t *testing.T) {
	var (
		vault vaults.Vault

		resourceName = "hcso_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupTurbo_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccVault_backupTurbo_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
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

// Vaults of type 'turbo'
func testAccVault_backupTurbo_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

resource "hcso_sfs_turbo" "test" {
  name              = "%[2]s_sfs_turbo"
  size              = 500
  share_proto       = "NFS"
  share_type        = "PERFORMANCE"
  vpc_id            = hcso_vpc.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  security_group_id = hcso_networking_secgroup.test.id
  availability_zone = data.hcso_availability_zones.test.names[0]
}`, common.TestBaseNetwork(name), name)
}

func testAccVault_backupTurbo_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 800
  enterprise_project_id = "0"

  resources {
    includes = [hcso_sfs_turbo.test.id]
  }
}
`, testAccVault_backupTurbo_base(rName), rName)
}

func testAccVault_backupTurbo_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 1000
  enterprise_project_id = "0"

  resources {
    includes =  [hcso_sfs_turbo.test.id]
  }
}
`, testAccVault_backupTurbo_base(rName), rName)
}

func TestAccVault_AutoBind(t *testing.T) {
	var (
		vault vaults.Vault

		randName     = acceptance.RandomAccResourceName()
		resourceName = "hcso_cbr_vault.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_autoBind(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "auto_bind", "true"),
					resource.TestCheckResourceAttr(resourceName, "bind_rules.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "bind_rules.key", "value"),
				),
			},
			{
				Config: testAccVault_autoBindUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "auto_bind", "false"),
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

func testAccVault_autoBind(rName string) string {
	return fmt.Sprintf(`
resource "hcso_cbr_vault" "test" {
  name            = "%[1]s"
  type            = "server"
  protection_type = "backup"
  size            = 800
  auto_bind       = true

  bind_rules = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccVault_autoBindUpdate(rName string) string {
	return fmt.Sprintf(`
resource "hcso_cbr_vault" "test" {
  name            = "%[1]s"
  type            = "server"
  protection_type = "backup"
  size            = 800
  auto_bind       = false
}
`, rName)
}

func TestAccVault_bindPolicies(t *testing.T) {
	var (
		vault vaults.Vault

		randName     = acceptance.RandomAccResourceName()
		mainRcName   = "hcso_cbr_vault.test"
		legacyRcName = "hcso_cbr_vault.legacy"
	)

	mainRc := acceptance.InitResourceCheck(
		mainRcName,
		&vault,
		getVaultResourceFunc,
	)

	legacyRc := acceptance.InitResourceCheck(
		legacyRcName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckReplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      mainRc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_bindPolicies_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					mainRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(mainRcName, "name", randName),
					resource.TestCheckResourceAttr(mainRcName, "policy.#", "2"),
					legacyRc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(legacyRcName, "policy_id",
						"hcso_cbr_policy.backup.0", "id"),
				),
			},
			{
				Config: testAccVault_bindPolicies_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					mainRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(mainRcName, "name", randName),
					resource.TestCheckResourceAttr(mainRcName, "policy.#", "2"),
					legacyRc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(legacyRcName, "policy_id",
						"hcso_cbr_policy.backup.1", "id"),
				),
			},
			{
				ResourceName:      mainRcName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      legacyRcName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"policy_id",
				},
			},
		},
	})
}

func testAccVault_bindPolicies_base(name string) string {
	return fmt.Sprintf(`
variable "backup_configuration" {
  type = list(object({
    interval        = number
    days            = string
    execution_times = list(string)
  }))
  default = [{
    interval        = 5
    days            = null
    execution_times = ["06:00", "18:00"]
  },
  {
    interval        = null
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }]
}

resource "hcso_cbr_policy" "backup" {
  count = 2

  name            = "%[1]s"
  type            = "backup"
  backup_quantity = 5

  backup_cycle {
    interval        = var.backup_configuration[count.index].interval
    days            = var.backup_configuration[count.index].days
    execution_times = var.backup_configuration[count.index].execution_times
  }
}

resource "hcso_cbr_policy" "replication" {
  count = 2

  name                   = "%[1]s"
  type                   = "replication"
  destination_region     = "%[2]s"
  destination_project_id = "%[3]s"
  time_period            = 20

  backup_cycle {
    interval        = var.backup_configuration[count.index].interval
    days            = var.backup_configuration[count.index].days
    execution_times = var.backup_configuration[count.index].execution_times
  }
}

resource "hcso_cbr_vault" "destination" {
  region          = "%[2]s"
  name            = "%[1]s"
  type            = "server"
  protection_type = "replication"
  size            = 200
}
`, name, acceptance.HCSO_DEST_REGION, acceptance.HCSO_DEST_PROJECT_ID)
}

func testAccVault_bindPolicies_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "legacy" {
  name            = "%[2]s_legacy"
  type            = "server"
  protection_type = "backup"
  size            = 200
  policy_id       = hcso_cbr_policy.backup[0].id
}

resource "hcso_cbr_vault" "test" {
  name            = "%[2]s"
  type            = "server"
  protection_type = "backup"
  size            = 200

  policy {
    id = hcso_cbr_policy.backup[0].id
  }
  policy {
    id                   = hcso_cbr_policy.replication[0].id
    destination_vault_id = hcso_cbr_vault.destination.id
  }
}
`, testAccVault_bindPolicies_base(name), name)
}

func testAccVault_bindPolicies_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "legacy" {
  name            = "%[2]s_legacy"
  type            = "server"
  protection_type = "backup"
  size            = 200
  policy_id       = hcso_cbr_policy.backup[1].id
}

resource "hcso_cbr_vault" "test" {
  name            = "%[2]s"
  type            = "server"
  protection_type = "backup"
  size            = 200

  policy {
    id = hcso_cbr_policy.backup[1].id
  }
  policy {
    id                   = hcso_cbr_policy.replication[1].id
    destination_vault_id = hcso_cbr_vault.destination.id
  }
}
`, testAccVault_bindPolicies_base(name), name)
}

func TestAccVault_backupWorkspace(t *testing.T) {
	var (
		vault vaults.Vault

		resourceName = "hcso_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		basicConfig  = testAccVault_backupWorkspace_base(name)
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupWorkspace_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeWorkspace),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccVault_backupWorkspace_step2(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeWorkspace),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
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

func testAccVault_backupWorkspace_base(name string) string {
	wsDesktopName := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

%[1]s

resource "hcso_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = hcso_vpc.test.id
  network_ids = [
    hcso_vpc_subnet.test.id,
  ]
}

resource "hcso_workspace_desktop" "test" {
  count = 2

  flavor_id         = "workspace.x86.ultimate.large2"
  image_type        = "market"
  image_id          = "63aa8670-27ad-4747-8c44-6d8919e785a7"
  availability_zone = data.hcso_availability_zones.test.names[0]
  vpc_id            = hcso_vpc.test.id
  security_groups   = [
    hcso_workspace_service.test.desktop_security_group.0.id,
    hcso_networking_secgroup.test.id,
  ]

  nic {
    network_id = hcso_vpc_subnet.test.id
  }

  name       = format("%[2]s-%%d", count.index)
  user_name  = format("user-%[3]s-%%d", count.index)
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }

  delete_user = true
}
`, common.TestBaseNetwork(name), wsDesktopName, name)
}

func testAccVault_backupWorkspace_step1(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "workspace"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "0"

  resources {
    server_id = hcso_workspace_desktop.test[0].id
  }
}
`, basicConfig, name)
}

func testAccVault_backupWorkspace_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "workspace"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 300
  enterprise_project_id = "0"

  resources {
    server_id = hcso_workspace_desktop.test[1].id
  }
}
`, basicConfig, name)
}

func TestAccVault_backupVMware(t *testing.T) {
	var (
		vault vaults.Vault

		resourceName  = "hcso_cbr_vault.test"
		resourceName1 = "hcso_cbr_vault.test.0"
		resourceName2 = "hcso_cbr_vault.test.1"
		name          = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupVMware_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					resource.TestCheckResourceAttr(resourceName1, "name", name+"_0"),
					resource.TestCheckResourceAttr(resourceName1, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName1, "type", cbr.VaultTypeVMware),
					resource.TestCheckResourceAttr(resourceName1, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName1, "size", "100"),
					resource.TestCheckResourceAttr(resourceName1, "enterprise_project_id", "0"),

					resource.TestCheckResourceAttr(resourceName2, "name", name+"_1"),
					resource.TestCheckResourceAttr(resourceName2, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName2, "type", cbr.VaultTypeVMware),
					resource.TestCheckResourceAttr(resourceName2, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName2, "size", "200"),
					resource.TestCheckResourceAttr(resourceName2, "enterprise_project_id", "0"),
				),
			},
		},
	})
}

func testAccVault_backupVMware_step1(name string) string {
	return fmt.Sprintf(`
variable "file_backup_configuration" {
  type = list(object({
    consistent_level = string
    size             = number
  }))

  default = [
    {
      consistent_level = "app_consistent"
      size             = 100
    },
    {
      consistent_level = "crash_consistent"
      size             = 200
    }
  ]
}

resource "hcso_cbr_vault" "test" {
  count = 2

  name              = format("%[1]s_%%d", count.index)
  type              = "vmware"
  consistent_level  = var.file_backup_configuration[count.index]["consistent_level"]
  protection_type   = "backup"
  size              = var.file_backup_configuration[count.index]["size"]
}
`, name)
}

func TestAccVault_backupFile(t *testing.T) {
	var (
		vault vaults.Vault

		resourceName  = "hcso_cbr_vault.test"
		resourceName1 = "hcso_cbr_vault.test.0"
		resourceName2 = "hcso_cbr_vault.test.1"
		name          = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupFile_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					resource.TestCheckResourceAttr(resourceName1, "name", name+"_0"),
					resource.TestCheckResourceAttr(resourceName1, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName1, "type", cbr.VaultTypeFile),
					resource.TestCheckResourceAttr(resourceName1, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName1, "size", "100"),
					resource.TestCheckResourceAttr(resourceName1, "enterprise_project_id", "0"),

					resource.TestCheckResourceAttr(resourceName2, "name", name+"_1"),
					resource.TestCheckResourceAttr(resourceName2, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName2, "type", cbr.VaultTypeFile),
					resource.TestCheckResourceAttr(resourceName2, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName2, "size", "200"),
					resource.TestCheckResourceAttr(resourceName2, "enterprise_project_id", "0"),
				),
			},
		},
	})
}

func testAccVault_backupFile_step1(name string) string {
	return fmt.Sprintf(`
variable "file_backup_configuration" {
  type = list(object({
    consistent_level = string
    size             = number
  }))

  default = [
    {
      consistent_level = "app_consistent"
      size             = 100
    },
    {
      consistent_level = "crash_consistent"
      size             = 200
    }
  ]
}

resource "hcso_cbr_vault" "test" {
  count = 2

  name              = format("%[1]s_%%d", count.index)
  type              = "file"
  consistent_level  = var.file_backup_configuration[count.index]["consistent_level"]
  protection_type   = "backup"
  size              = var.file_backup_configuration[count.index]["size"]
}
`, name)
}
