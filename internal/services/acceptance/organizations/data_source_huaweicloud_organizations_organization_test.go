package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourceOrganization_basic(t *testing.T) {
	rName := "data.hcso_organizations_organization.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceOrganization_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "root_tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "root_tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccDatasourceOrganization_basic() string {
	return fmt.Sprintf(`
%s

data "hcso_organizations_organization" "test" {
  depends_on = [hcso_organizations_organization.test]
}
`, testOrganization_basic())
}
