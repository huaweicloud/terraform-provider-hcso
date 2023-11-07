package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func getHostGroupFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcHssV5Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS v5 client: %s", err)
	}

	return hss.QueryHostGroupById(client, acceptance.HCSO_REGION_NAME, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST,
		state.Primary.ID)
}

func TestAccHostGroup_basic(t *testing.T) {
	var (
		group *hssv5model.HostGroupItem

		name  = acceptance.RandomAccResourceName()
		rName = "hcso_hss_host_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&group,
		getHostGroupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccHostGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "1"),
				),
			},
			{
				Config: testAccHostGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccHostGroup_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_kps_keypair" "test" {
  name = "%[2]s"
}

resource "hcso_compute_instance" "test" {
  count = 2

  name               = "%[2]s"
  image_id           = data.hcso_images_image.test.id
  flavor_id          = data.hcso_compute_flavors.test.ids[0]
  security_groups    = [hcso_networking_secgroup.test.name]
  availability_zone  = data.hcso_availability_zones.test.names[0]

  key_pair   = hcso_kps_keypair.test.name
  agent_list = "hss"

  network {
    uuid = hcso_vpc_subnet.test.id
  }
}
`, common.TestBaseComputeResources(name), name)
}

func testAccHostGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_hss_host_group" "test" {
  name     = "%[2]s"
  host_ids = slice(hcso_compute_instance.test[*].id, 0, 1)
}
`, testAccHostGroup_base(name), name)
}

func testAccHostGroup_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_hss_host_group" "test" {
  name     = "%[2]s-update"
  host_ids = hcso_compute_instance.test[*].id
}
`, testAccHostGroup_base(name), name)
}
