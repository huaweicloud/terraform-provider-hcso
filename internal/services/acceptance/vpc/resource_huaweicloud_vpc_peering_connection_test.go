package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/peerings"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getPeeringConnectionResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC Peering Connection client: %s", err)
	}
	return peerings.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcPeeringConnection_basic(t *testing.T) {
	var peering peerings.Peering

	randName := acceptance.RandomAccResourceName()
	updateName := randName + "_update"
	basicDesc := "vpc1 peers to vpc2"
	updateDesc := "vpc1 peering to vpc2"
	resourceName := "hcso_vpc_peering_connection.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&peering,
		getPeeringConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnection_config(randName, randName, basicDesc),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", basicDesc),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "hcso_vpc.test1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "peer_vpc_id", "hcso_vpc.test2", "id"),
				),
			},
			{
				Config: testAccVpcPeeringConnection_config(randName, updateName, updateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", updateDesc),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "hcso_vpc.test1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "peer_vpc_id", "hcso_vpc.test2", "id"),
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

func testAccVpcPeeringConnection_config(vpcName, peerName, desc string) string {
	return fmt.Sprintf(`
resource "hcso_vpc" "test1" {
  name = "%[1]s_1"
  cidr = "172.16.0.0/20"
}

resource "hcso_vpc" "test2" {
  name = "%[1]s_2"
  cidr = "172.16.128.0/20"
}

resource "hcso_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = hcso_vpc.test1.id
  peer_vpc_id = hcso_vpc.test2.id
  description = "%s"
}
`, vpcName, peerName, desc)
}
