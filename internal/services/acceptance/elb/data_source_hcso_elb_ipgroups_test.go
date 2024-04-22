package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourceIpGroups_basic(t *testing.T) {
	rName := "data.hcso_elb_ipgroups.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIpGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "ipgroups.#"),
					resource.TestCheckResourceAttrSet(rName, "ipgroups.0.name"),
					resource.TestCheckResourceAttrSet(rName, "ipgroups.0.id"),
					resource.TestCheckResourceAttrSet(rName, "ipgroups.0.description"),
					resource.TestCheckResourceAttrSet(rName, "ipgroups.0.project_id"),
					resource.TestCheckResourceAttrSet(rName, "ipgroups.0.ip_list.0.ip"),
					resource.TestCheckResourceAttrSet(rName, "ipgroups.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "ipgroups.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("ipgroup_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_address_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceIpGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_elb_ipgroups" "test" {
  depends_on = [hcso_elb_ipgroup.test]
}

data "hcso_elb_ipgroups" "name_filter" {
  depends_on = [hcso_elb_ipgroup.test]
  name       = "%[2]s"
}
output "name_filter_is_useful" {
  value = length(data.hcso_elb_ipgroups.name_filter.ipgroups) > 0 && alltrue(
  [for v in data.hcso_elb_ipgroups.name_filter.ipgroups[*].name :v == "%[2]s"]
  )  
}

locals {
  ipgroup_id = hcso_elb_ipgroup.test.id
}
data "hcso_elb_ipgroups" "ipgroup_id_filter" {
  ipgroup_id = hcso_elb_ipgroup.test.id
}
output "ipgroup_id_filter_is_useful" {
  value = length(data.hcso_elb_ipgroups.ipgroup_id_filter.ipgroups) > 0 && alltrue(
  [for v in data.hcso_elb_ipgroups.ipgroup_id_filter.ipgroups[*].id : v == local.ipgroup_id]
  )  
}

locals {
  ip_address = hcso_elb_ipgroup.test.ip_list[0].ip
}
data "hcso_elb_ipgroups" "ip_address_filter" {
  ip_address = hcso_elb_ipgroup.test.ip_list[0].ip
}
output "ip_address_filter_is_useful" {
  value = length(data.hcso_elb_ipgroups.ip_address_filter.ipgroups) > 0 && alltrue(
  [for v in data.hcso_elb_ipgroups.ip_address_filter.ipgroups[*].ip_list[0].ip : v == local.ip_address]
  )  
}

locals {
  description = hcso_elb_ipgroup.test.description
}
data "hcso_elb_ipgroups" "description_filter" {
  description = hcso_elb_ipgroup.test.description
}
output "description_filter_is_useful" {
  value = length(data.hcso_elb_ipgroups.description_filter.ipgroups) > 0 && alltrue(
  [for v in data.hcso_elb_ipgroups.description_filter.ipgroups[*].description : v == local.description]
  )  
}

`, testAccElbV3IpGroupConfig_basic(name), name)
}
