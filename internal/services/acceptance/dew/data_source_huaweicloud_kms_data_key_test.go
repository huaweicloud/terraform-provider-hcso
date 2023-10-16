package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccKmsDataKeyV1DataSource_basic(t *testing.T) {
	keyAlias := acceptance.RandomAccResourceName()
	datasourceName := "data.hcso_kms_data_key.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsDataKeyV1DataSource_basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(
						datasourceName, "plain_text"),
					resource.TestCheckResourceAttrSet(
						datasourceName, "cipher_text"),
				),
			},
		},
	})
}

func testAccKmsDataKeyV1DataSource_basic(keyAlias string) string {
	return fmt.Sprintf(`
resource "hcso_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

data "hcso_kms_data_key" "test" {
  key_id         = hcso_kms_key.test.id
  datakey_length = "512"
}
`, keyAlias)
}
