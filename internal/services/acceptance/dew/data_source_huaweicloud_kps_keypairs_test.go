package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDataKpsKeypairs_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "data.hcso_kps_keypairs.test"
	publicKey, _, _ := acctest.RandSSHKeyPair("Generated-by-AccTest")

	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKpsKeypairs_basic(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "keypairs.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "keypairs.0.public_key", publicKey),
					resource.TestCheckResourceAttrPair(resourceName, "keypairs.0.scope",
						"hcso_kps_keypair.test", "scope"),
					resource.TestCheckResourceAttrPair(resourceName, "keypairs.0.fingerprint",
						"hcso_kps_keypair.test", "fingerprint"),
					resource.TestCheckResourceAttrPair(resourceName, "keypairs.0.is_managed",
						"hcso_kps_keypair.test", "is_managed"),
				),
			},
		},
	})
}

func testAccDataKpsKeypairs_basic(rName, key string) string {
	return fmt.Sprintf(`
%s

data "hcso_kps_keypairs" "test" {
  name = hcso_kps_keypair.test.name

  depends_on = [hcso_kps_keypair.test]
}
`, testKeypair_publicKey(rName, key))
}
