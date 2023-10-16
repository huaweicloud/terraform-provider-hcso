package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccIdentityGroupDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcso_identity_group.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupDataSource_by_name(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "users.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityGroupDataSource_with_user(t *testing.T) {
	dataSourceName := "data.hcso_identity_group.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupDataSource_with_user(rName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "users.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.id"),
				),
			},
		},
	})
}

func testAccIdentityGroupDataSource_by_name(rName string) string {
	return fmt.Sprintf(`
resource "hcso_identity_group" "test" {
  name        = "%s"
  description = "An ACC test group"
}

data "hcso_identity_group" "test" {
  name = hcso_identity_group.test.name
  
  depends_on = [
    hcso_identity_group.test
  ]
}
`, rName)
}

func testAccIdentityGroupDataSource_with_user(rName, password string) string {
	return fmt.Sprintf(`
resource "hcso_identity_group" "test" {
  name        = "%[1]s"
  description = "An ACC test group"
}

resource "hcso_identity_user" "test" {
  count    = 2
  name     = "%[1]s-${count.index}"
  password = "%[2]s"
  enabled  = true
}

resource "hcso_identity_group_membership" "test" {
  group = hcso_identity_group.test.id
  users = hcso_identity_user.test.*.id
}

data "hcso_identity_group" "test" {
  name = hcso_identity_group.test.name

  depends_on = [
    hcso_identity_group_membership.test
  ]
}
`, rName, password)
}
