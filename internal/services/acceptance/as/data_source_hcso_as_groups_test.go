package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDataSourceASGroup_basic(t *testing.T) {
	dataSourceName := "data.hcso_as_groups.groups"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceASGroup_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "groups.0.scaling_group_name", name),
				),
			},
		},
	})
}

func testAccDataSourceASGroup_conf(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_as_groups" "groups" {
  name = hcso_as_group.acc_as_group.scaling_group_name

  depends_on = [hcso_as_group.acc_as_group]
}
`, testASGroup_basic(name))
}
