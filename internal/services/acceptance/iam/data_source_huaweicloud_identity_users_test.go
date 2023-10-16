package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccIdentityUsersDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcso_identity_users.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.name"),
				),
			},
		},
	})
}

const testAccUsersDataSourceBasic string = `
data "hcso_identity_users" "test" {
}
`
