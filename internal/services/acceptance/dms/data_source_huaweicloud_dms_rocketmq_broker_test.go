package dms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
)

func TestAccDatasourceDmsRocketMQBroker_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcso_dms_rocketmq_broker.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDmsRocketMQBroker_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "brokers.0", "broker-0"),
				),
			},
		},
	})
}

func testAccDatasourceDmsRocketMQBroker_base(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "4.8.0"
  storage_space     = 600
  vpc_id            = hcso_vpc.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  security_group_id = hcso_networking_secgroup.test.id

  availability_zones = [
    data.hcso_availability_zones.test.names[0]
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourceDmsRocketMQBroker_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_dms_rocketmq_broker" "test" {
  instance_id = hcso_dms_rocketmq_instance.test.id
}
`, testAccDatasourceDmsRocketMQBroker_base(name))
}
