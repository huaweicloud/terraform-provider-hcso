package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
)

func TestAccGeminiDBInstanceDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceDataSourceID("data.hcso_gaussdb_cassandra_instance.test"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find GaussDB cassandra instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("GaussDB cassandra instance data source ID not set ")
		}

		return nil
	}
}

func testAccGeminiDBInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_gaussdb_cassandra_instance" "test" {
  name        = "%s"
  password    = "Test@12345678"
  flavor      = "geminidb.cassandra.xlarge.4"
  volume_size = 100
  vpc_id      = hcso_vpc.test.id
  subnet_id   = hcso_vpc_subnet.test.id
  node_num    = 4

  security_group_id = hcso_networking_secgroup.test.id
  availability_zone = data.hcso_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "hcso_gaussdb_cassandra_instance" "test" {
  name = hcso_gaussdb_cassandra_instance.test.name
  depends_on = [
    hcso_gaussdb_cassandra_instance.test,
  ]
}
`, common.TestBaseNetwork(rName), rName)
}
