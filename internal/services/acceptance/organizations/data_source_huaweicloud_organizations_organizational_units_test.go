package organizations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourceOrganizationalUnits_basic(t *testing.T) {
	rName := "data.hcso_organizations_organizational_units.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceOrganizationalUnits_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "children.#"),
					resource.TestCheckResourceAttrSet(rName, "children.0.id"),
					resource.TestCheckResourceAttrSet(rName, "children.0.name"),
					resource.TestCheckResourceAttrSet(rName, "children.0.urn"),
					resource.TestCheckResourceAttrSet(rName, "children.0.created_at"),
				),
			},
		},
	})
}

func testAccDatasourceOrganizationalUnits_basic() string {
	return `
data "hcso_organizations_organization" "test" {}

data "hcso_organizations_organizational_units" "test" {
  parent_id = data.hcso_organizations_organization.test.root_id
}
`
}
