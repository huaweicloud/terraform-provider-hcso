package sfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sfs_turbo/v1/shares"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func getSfsTurboResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SfsV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS client: %s", err)
	}

	resourceID := state.Primary.ID
	share, err := shares.Get(client, resourceID).Extract()
	if err != nil {
		return nil, err
	}

	if share.ID == resourceID {
		return &share, nil
	}

	return nil, fmt.Errorf("the sfs turbo %s does not exist", resourceID)
}

func TestAccSFSTurbo_basic(t *testing.T) {
	var turbo shares.Turbo
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcso_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "PERFORMANCE"),
					resource.TestCheckResourceAttr(resourceName, "enhanced", "false"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcso_networking_secgroup.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSFSTurbo_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "size", "600"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcso_networking_secgroup.test_update", "id"),
					resource.TestCheckResourceAttr(resourceName, "status", "232"),
				),
			},
		},
	})
}

// TODO failed by resource "hcso_kms_key" of service DEW not supported currently.
func TestAccSFSTurbo_crypt(t *testing.T) {
	var turbo shares.Turbo
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcso_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_crypt(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "enhanced", "false"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
					resource.TestCheckResourceAttrSet(resourceName, "crypt_key_id"),
				),
			},
		},
	})
}

func TestAccSFSTurbo_withEpsId(t *testing.T) {
	var turbo shares.Turbo
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcso_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccSFSTurbo_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_sfs_turbo" "test" {
  name              = "%s"
  size              = 500
  share_proto       = "NFS"
  share_type        = "PERFORMANCE"
  vpc_id            = hcso_vpc.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  security_group_id = hcso_networking_secgroup.test.id
  availability_zone = data.hcso_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccSFSTurbo_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_networking_secgroup" "test_update" {
  name                 = "%[2]s_update"
  delete_default_rules = true
}

data "hcso_availability_zones" "test" {}

resource "hcso_sfs_turbo" "test" {
  name              = "%[2]s_update"
  size              = 600
  share_proto       = "NFS"
  share_type        = "PERFORMANCE"
  vpc_id            = hcso_vpc.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  security_group_id = hcso_networking_secgroup.test_update.id
  availability_zone = data.hcso_availability_zones.test.names[0]

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccSFSTurbo_crypt(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "hcso_sfs_turbo" "test" {
  name              = "%s"
  size              = 500
  share_proto       = "NFS"
  share_type        = "PERFORMANCE"
  vpc_id            = hcso_vpc.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  security_group_id = hcso_networking_secgroup.test.id
  availability_zone = data.hcso_availability_zones.test.names[0]
  crypt_key_id      = hcso_kms_key.test.id
}
`, common.TestBaseNetwork(rName), rName, rName)
}

func testAccSFSTurbo_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_sfs_turbo" "test" {
  name                   = "%s"
  size                   = 500
  share_proto            = "NFS"
  share_type             = "PERFORMANCE"
  vpc_id                 = hcso_vpc.test.id
  subnet_id              = hcso_vpc_subnet.test.id
  security_group_id      = hcso_networking_secgroup.test.id
  availability_zone      = data.hcso_availability_zones.test.names[0]
  enterprise_project_id  = "%s"
}
`, common.TestBaseNetwork(rName), rName, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST)
}
