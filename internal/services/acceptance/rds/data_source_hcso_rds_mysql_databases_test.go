package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccMysqlDatabases_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.hcso_rds_mysql_databases.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDatabases_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "databases.#"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.name"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.character_set"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.description"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("character_set_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccMysqlDatabases_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_rds_mysql_databases" "test" {
  depends_on  = [hcso_rds_mysql_database.test]
  instance_id = hcso_rds_instance.test.id
}

data "hcso_rds_mysql_databases" "name_filter" {
  depends_on  = [hcso_rds_mysql_database.test]
  instance_id = hcso_rds_instance.test.id
  name        = hcso_rds_mysql_database.test.name
}

locals {
  name = hcso_rds_mysql_database.test.name
}
	
output "name_filter_is_useful" {
  value = length(data.hcso_rds_mysql_databases.name_filter.databases) > 0 && alltrue(
  [for v in data.hcso_rds_mysql_databases.name_filter.databases[*].name : v == local.name]
  )
}

data "hcso_rds_mysql_databases" "character_set_filter" {
  depends_on    = [hcso_rds_mysql_database.test]
  instance_id   = hcso_rds_instance.test.id
  character_set = hcso_rds_mysql_database.test.character_set
}

locals {
  character_set = hcso_rds_mysql_database.test.character_set
}
	
output "character_set_filter_is_useful" {
  value = length(data.hcso_rds_mysql_databases.character_set_filter.databases) > 0 && alltrue(
  [for v in data.hcso_rds_mysql_databases.character_set_filter.databases[*].character_set : v == local.character_set]
  )
}
`, testMysqlDatabase_basic(name, "", "test_database"))
}
