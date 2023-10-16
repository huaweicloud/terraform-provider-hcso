package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccGaussdbMysqlConfigurationDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussdbMysqlConfigurationDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hcso_gaussdb_mysql_configuration.test", "name", "Default-GaussDB-for-MySQL 8.0"),
				),
			},
		},
	})
}

const testAccGaussdbMysqlConfigurationDataSource_basic = `
data "hcso_gaussdb_mysql_configuration" "test" {
  name = "Default-GaussDB-for-MySQL 8.0"
}
`
