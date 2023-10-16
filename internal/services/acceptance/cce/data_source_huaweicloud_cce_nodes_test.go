package cce

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNodesDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcso_cce_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.0.name", rName),
				),
			},
		},
	})
}

func testAccNodesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_cce_nodes" "test" {
  cluster_id = hcso_cce_cluster.test.id
  name       = hcso_cce_node.test.name

  depends_on = [hcso_cce_node.test]
}
`, testAccCceCluster_config(rName))
}
