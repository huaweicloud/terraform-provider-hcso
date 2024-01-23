package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func TestAccDatasourceMembers_basic(t *testing.T) {
	rName := "data.hcso_elb_members.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMembers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "members.#"),
					resource.TestCheckResourceAttrSet(rName, "members.0.name"),
					resource.TestCheckResourceAttrSet(rName, "members.0.id"),
					resource.TestCheckResourceAttrSet(rName, "members.0.address"),
					resource.TestCheckResourceAttrSet(rName, "members.0.protocol_port"),
					resource.TestCheckResourceAttrSet(rName, "members.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "members.0.weight"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("member_id_filter_is_useful", "true"),
					resource.TestCheckOutput("address_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_port_filter_is_useful", "true"),
					resource.TestCheckOutput("weight_filter_is_useful", "true"),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("member_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccElbMemberConfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

resource "hcso_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = hcso_vpc.test.id
  ipv4_subnet_id = hcso_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.hcso_availability_zones.test.names[0]
  ]
  backend_subnets = [
    hcso_vpc_subnet.test.id
 ]
}

resource "hcso_elb_listener" "test" {
  name                        = "%[2]s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8083
  loadbalancer_id             = hcso_elb_loadbalancer.test.id
  advanced_forwarding_enabled = false
}

resource "hcso_elb_pool" "test" {
  name        = "%[2]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = hcso_elb_listener.test.id
}

resource "hcso_elb_member" "test" {
  name          = "%[2]s"
  address       = "192.168.0.10"
  protocol_port = 8080
  pool_id       = hcso_elb_pool.test.id
  subnet_id     = hcso_vpc_subnet.test.ipv4_subnet_id
}
`, common.TestVpc(name), name)
}

func testAccDatasourceMembers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_elb_members" "test" {
  pool_id    = hcso_elb_pool.test.id
  depends_on = [hcso_elb_member.test]
}

data "hcso_elb_members" "name_filter" {
  pool_id    = hcso_elb_pool.test.id
  name       = "%[2]s"
  depends_on = [hcso_elb_member.test]
}

output "name_filter_is_useful" {
  value = length(data.hcso_elb_members.name_filter.members) > 0 && alltrue(
  [for v in data.hcso_elb_members.name_filter.members[*].name :v == "%[2]s"]
  )  
}

data "hcso_elb_members" "member_id_filter" {
  pool_id    = hcso_elb_pool.test.id
  member_id  = hcso_elb_member.test.id
  depends_on = [hcso_elb_member.test]
}

locals {
  member_id = hcso_elb_member.test.id
}

output "member_id_filter_is_useful" {
  value = length(data.hcso_elb_members.member_id_filter.members) > 0 && alltrue(
  [for v in data.hcso_elb_members.member_id_filter.members[*].id : v == local.member_id]
  )  
}

data "hcso_elb_members" "address_filter" {
  pool_id    = hcso_elb_pool.test.id
  address    = hcso_elb_member.test.address
  depends_on = [hcso_elb_member.test]
}

locals {
  address = hcso_elb_member.test.address
}

output "address_filter_is_useful" {
  value = length(data.hcso_elb_members.address_filter.members) > 0 && alltrue(
  [for v in data.hcso_elb_members.address_filter.members[*].address : v == local.address]
  )  
}

data "hcso_elb_members" "protocol_port_filter" {
  pool_id       = hcso_elb_pool.test.id
  protocol_port = hcso_elb_member.test.protocol_port
  depends_on    = [hcso_elb_member.test]
}

locals {
  protocol_port = hcso_elb_member.test.protocol_port
}

output "protocol_port_filter_is_useful" {
  value = length(data.hcso_elb_members.protocol_port_filter.members) > 0 && alltrue(
  [for v in data.hcso_elb_members.protocol_port_filter.members[*].protocol_port : v == local.protocol_port]
  )  
}

data "hcso_elb_members" "weight_filter" {
  pool_id    = hcso_elb_pool.test.id
  weight     = hcso_elb_member.test.weight
  depends_on = [hcso_elb_member.test]
}

locals {
  weight = hcso_elb_member.test.weight
}

output "weight_filter_is_useful" {
  value = length(data.hcso_elb_members.weight_filter.members) > 0 && alltrue(
  [for v in data.hcso_elb_members.weight_filter.members[*].weight : v == local.weight]
  )  
}

data "hcso_elb_members" "subnet_id_filter" {
  pool_id    = hcso_elb_pool.test.id
  subnet_id  = hcso_elb_member.test.subnet_id
  depends_on = [hcso_elb_member.test]
}

locals {
  subnet_id = hcso_elb_member.test.subnet_id
}

output "subnet_id_filter_is_useful" {
  value = length(data.hcso_elb_members.subnet_id_filter.members) > 0 && alltrue(
  [for v in data.hcso_elb_members.subnet_id_filter.members[*].subnet_id : v == local.subnet_id]
  )  
}

data "hcso_elb_members" "member_type_filter" {
  pool_id     = hcso_elb_pool.test.id
  member_type = "instance"
  depends_on  = [hcso_elb_member.test]
}

output "member_type_filter_is_useful" {
  value = length(data.hcso_elb_members.member_type_filter.members) > 0 && alltrue(
  [for v in data.hcso_elb_members.member_type_filter.members[*].member_type : v == "instance"]
  )  
}
`, testAccElbMemberConfig_basic(name), name)
}
