package cnad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourceAvailableObjects_basic(t *testing.T) {
	rName := "data.hcso_cnad_advanced_available_objects.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCNADInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAvailableObjects_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.id"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.type"),
					resource.TestCheckOutput("protected_object_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_address_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testAvailableObjects_basic = `
data "hcso_cnad_advanced_instances" "test" {}

data "hcso_cnad_advanced_available_objects" "test" {
  instance_id = data.hcso_cnad_advanced_instances.test.instances.0.instance_id
}

data "hcso_cnad_advanced_available_objects" "protected_object_id_filter" {
  instance_id         = data.hcso_cnad_advanced_instances.test.instances.0.instance_id
  protected_object_id = data.hcso_cnad_advanced_available_objects.test.protected_objects.0.id
}
output "protected_object_id_filter_is_useful" {
  value = alltrue([for v in data.hcso_cnad_advanced_available_objects.protected_object_id_filter.
  protected_objects[*].id : v == data.hcso_cnad_advanced_available_objects.
  test.protected_objects.0.id])
}

data "hcso_cnad_advanced_available_objects" "ip_address_filter" {
  instance_id = data.hcso_cnad_advanced_instances.test.instances.0.instance_id
  ip_address  = data.hcso_cnad_advanced_available_objects.test.protected_objects.0.ip_address
}
output "ip_address_filter_is_useful" {
  value = alltrue([for v in data.hcso_cnad_advanced_available_objects.ip_address_filter.protected_objects[*].
  ip_address : v == data.hcso_cnad_advanced_available_objects.test.protected_objects.0.ip_address])
}

data "hcso_cnad_advanced_available_objects" "type_filter" {
  instance_id = data.hcso_cnad_advanced_instances.test.instances.0.instance_id
  type       = data.hcso_cnad_advanced_available_objects.test.protected_objects.0.type
}
output "type_filter_is_useful" {
  value = alltrue([for v in data.hcso_cnad_advanced_available_objects.type_filter.protected_objects[*].
  type : v == data.hcso_cnad_advanced_available_objects.test.protected_objects.0.type])
}
`
