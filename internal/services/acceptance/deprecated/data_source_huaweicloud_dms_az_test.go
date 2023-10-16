package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDmsAZDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcso_dms_az.az1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsAZDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv6_enable"),
				),
			},
		},
	})
}

var testAccDmsAZDataSource_basic = fmt.Sprintf(`
data "hcso_dms_az" "az1" {}
`)
