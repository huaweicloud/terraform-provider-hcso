package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourcePgPlugins_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.hcso_rds_pg_plugins.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgPlugins_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.name"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.version"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.created"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.description"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("version_filter_is_useful", "true"),

					resource.TestCheckOutput("created_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePgPlugins_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_rds_pg_plugins" "test" {
  instance_id   = hcso_rds_instance.test.id
  database_name = hcso_rds_pg_database.test.name
}

data "hcso_rds_pg_plugins" "name_filter" {
  instance_id   = hcso_rds_instance.test.id
  database_name = hcso_rds_pg_database.test.name
  name          = data.hcso_rds_pg_plugins.test.plugins[0].name
}
locals {
  name = data.hcso_rds_pg_plugins.test.plugins[0].name
}
output "name_filter_is_useful" {
  value = length(data.hcso_rds_pg_plugins.name_filter.plugins) > 0 && alltrue(
    [for v in data.hcso_rds_pg_plugins.name_filter.plugins[*].name : v == local.name]
  )  
}

data "hcso_rds_pg_plugins" "version_filter" {
  instance_id   = hcso_rds_instance.test.id
  database_name = hcso_rds_pg_database.test.name
  version       = data.hcso_rds_pg_plugins.test.plugins[0].version

}
locals {
  version = data.hcso_rds_pg_plugins.test.plugins[0].version
}
output "version_filter_is_useful" {
  value = length(data.hcso_rds_pg_plugins.version_filter.plugins) > 0 && alltrue(
    [for v in data.hcso_rds_pg_plugins.version_filter.plugins[*].version : v == local.version]
  )  
}

data "hcso_rds_pg_plugins" "created_filter" {
  instance_id   = hcso_rds_instance.test.id
  database_name = hcso_rds_pg_database.test.name
  created       = data.hcso_rds_pg_plugins.test.plugins[0].created

}
locals {
  created = data.hcso_rds_pg_plugins.test.plugins[0].created
}
output "created_filter_is_useful" {
  value = length(data.hcso_rds_pg_plugins.created_filter.plugins) > 0 && alltrue(
    [for v in data.hcso_rds_pg_plugins.created_filter.plugins[*].created : v == local.created]
  )  
}
`, testPgDatabase_basic(name, ""))
}
