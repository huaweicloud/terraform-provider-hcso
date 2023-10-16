package sfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccSFSFileSystemV2DataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcso_sfs_file_system.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "available"),
					resource.TestCheckResourceAttr(dataSourceName, "size", "10"),
				),
			},
		},
	})
}

func testAccSFSFileSystemV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "hcso_vpc" "test" {
  name = "vpc-default"
}

data "hcso_availability_zones" "test" {}

resource "hcso_sfs_file_system" "test" {
  share_proto       = "NFS"
  size              = 10
  name              = "%s"
  description       = "sfs_c2c_test-file"
  access_to         = data.hcso_vpc.test.id
  access_level      = "rw"
  availability_zone = data.hcso_availability_zones.test.names[0]
}

data "hcso_sfs_file_system" "test" {
  id = hcso_sfs_file_system.test.id
}
`, rName)
}
