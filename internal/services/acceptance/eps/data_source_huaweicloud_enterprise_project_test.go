package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccEnterpriseProjectDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcso_enterprise_project.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProjectDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "default"),
					resource.TestCheckResourceAttr(dataSourceName, "id", "0"),
				),
			},
		},
	})
}

const testAccEnterpriseProjectDataSource_basic = `
data "hcso_enterprise_project" "test" {
  name = "default"
}
`
