package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/servicestage/v1/repositories"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccRepoPwdAuth_basic(t *testing.T) {
	var (
		auth         repositories.Authorization
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcso_servicestage_repo_password_authorization.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&auth,
		getAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRepoPwdAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRepoPwdAuth_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "devcloud"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"user_name",
					"password",
				},
			},
		},
	})
}

func testAccRepoPwdAuth_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcso_servicestage_repo_password_authorization" "test" {
  type      = "devcloud"
  name      = "%s"
  user_name = "%s/%s"
  password  = "%s"
}
`, rName, acceptance.HCSO_DOMAIN_NAME, acceptance.HCSO_USER_NAME, acceptance.HCSO_GITHUB_REPO_PWD)
}
