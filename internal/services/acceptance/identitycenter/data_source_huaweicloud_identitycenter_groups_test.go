package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourceIdentityCenterGroups_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.hcso_identitycenter_groups.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.name"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.description"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceIdentityCenterGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_identitycenter_groups" "test" {
  identity_store_id = data.hcso_identitycenter_instance.test.identity_store_id
  name              = hcso_identitycenter_group.test.name
}

data "hcso_identitycenter_groups" "name_filter" {
  identity_store_id = data.hcso_identitycenter_instance.test.identity_store_id
  name              = hcso_identitycenter_group.test.name
}

locals {
  name_filter_result = [for v in data.hcso_identitycenter_groups.name_filter.groups[*].name:
  v == data.hcso_identitycenter_groups.test.groups.0.name]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`, testIdentityCenterGroup_basic(name))
}
