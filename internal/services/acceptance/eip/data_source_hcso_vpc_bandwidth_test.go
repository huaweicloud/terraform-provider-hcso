package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccBandWidthDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcso_vpc_bandwidth.test"
	eipResourceName := "hcso_vpc_eip.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "size", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "publicips.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "publicips.0.id",
						eipResourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "publicips.0.ip_address",
						eipResourceName, "address"),
				),
			},
		},
	})
}

func testAccBandWidthDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcso_vpc_bandwidth" "test" {
  name = "%s"
  size = 10
}

resource "hcso_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type = "WHOLE"
    id         = hcso_vpc_bandwidth.test.id
  }
}

data "hcso_vpc_bandwidth" "test" {
  depends_on = [hcso_vpc_eip.test]

  name = hcso_vpc_bandwidth.test.name
}
`, rName)
}
